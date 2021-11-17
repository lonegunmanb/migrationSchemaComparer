# migration schema comparer for Terraform

When develop Terraform provider sometimes we need do some state migration(not schema migration) via `StateUpgraders`, in this case we need provide a Point-In-Time schema snapshot which exactly as same as the one defined in resource code. For some large resource(like `azurerm_kubernetes_cluster`) this PIT schema which expanded manually could contain hundreds even thousands lines. It's hard for human to check whether this snapshot is exactly same as one defined in resource code.

This helper function set all function field inside to nil, then use `reflect.DeepEqual` to compare two schemas. If there are some differences, the function will return an error indicating the path.

For example:

```go
resourceSchema := map[string]*schema.Schema{
	"field": {
		Type:     schema.TypeList,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
			"embedded_field": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: false,
				},
			},
		},
		MaxItems: 100,
		MinItems: 1,
	},
}
migrationSchema := map[string]*schema.Schema{
	"field": {
		Type:     schema.TypeList,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"embedded_field": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		},
		MaxItems: 100,
		MinItems: 1,
	},
}
err := migrationSchemaComparer.Equal(resourceSchema, migrationSchema,"")
```

The error message would start with `.field.embedded_field`

This helper function cannot be import into unit test that committed into repository since the resource schema would change in the future, so any unit test do this compare would fail eventually, we can only use it as a temporary unit test to make sure that the PIT schema makes no change. After that, we should delete this unit test.