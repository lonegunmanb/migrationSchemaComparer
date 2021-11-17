package migrationSchemaComparer_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lonegunmanb/migrationSchemaComparer"
	"github.com/stretchr/testify/assert"
)

func TestCompareEmptySchema(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{}
	migrationSchema := map[string]*schema.Schema{}
	assert.Nil(t, migrationSchemaComparer.Equal(resourceSchema, migrationSchema, ""))
}

func TestCompareSinglePrimaryTypeSchema(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Required: true,
			Default:  1,
		},
	}
	migrationSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Required: true,
			Default:  1,
		},
	}
	assert.Nil(t, migrationSchemaComparer.Equal(resourceSchema, migrationSchema, ""))
}

func TestCompareMultiplePrimaryTypeSchemas(t *testing.T) {
	cases := []struct {
		name            string
		resourceSchema  map[string]*schema.Schema
		migrationSchema map[string]*schema.Schema
		expected        bool
	}{
		{
			name: "different key",
			resourceSchema: map[string]*schema.Schema{
				"field1": {
					Type:     schema.TypeInt,
					ForceNew: true,
					Required: true,
					Default:  1,
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field2": {
					Type:     schema.TypeInt,
					ForceNew: true,
					Required: true,
					Default:  1,
				},
			},
			expected: false,
		},
		{
			name:           "different count",
			resourceSchema: map[string]*schema.Schema{},
			migrationSchema: map[string]*schema.Schema{
				"field1": {},
			},
			expected: false,
		},
		{
			name: "same items different order",
			resourceSchema: map[string]*schema.Schema{
				"field1": {
					Type:     schema.TypeInt,
					Default:  1,
					Required: true,
					Computed: false,
				},
				"field2": {
					Type:     schema.TypeString,
					Default:  "hello",
					Required: false,
					Computed: true,
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field2": {
					Type:     schema.TypeString,
					Default:  "hello",
					Required: false,
					Computed: true,
				},
				"field1": {
					Type:     schema.TypeInt,
					Default:  1,
					Required: true,
					Computed: false,
				},
			},
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := migrationSchemaComparer.Equal(c.resourceSchema, c.migrationSchema, "")
			if c.expected {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestComplexTypeField(t *testing.T) {
	cases := []struct {
		name            string
		resourceSchema  map[string]*schema.Schema
		migrationSchema map[string]*schema.Schema
		expected        bool
	}{
		{
			name: "equal list",
			resourceSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
					ValidateFunc: func(i interface{}, s string) ([]string, []error) {
						return nil, nil
					},
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
			expected: true,
		},
		{
			name: "not equal list 1",
			resourceSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 2,
					ValidateFunc: func(i interface{}, s string) ([]string, []error) {
						return nil, nil
					},
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
		},
		{
			name: "not equal list 2",
			resourceSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeString,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
					ValidateFunc: func(i interface{}, s string) ([]string, []error) {
						return nil, nil
					},
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
		},
		{
			name: "object type equal",
			resourceSchema: map[string]*schema.Schema{
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
			},
			migrationSchema: map[string]*schema.Schema{
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
			},
			expected: true,
		},
		{
			name: "object type not equal 1",
			resourceSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"embedded_field1": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
							},
						},
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"embedded_field2": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
							},
						},
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
			expected: false,
		},
		{
			name: "object type not equal 2",
			resourceSchema: map[string]*schema.Schema{
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
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeList,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"embedded_field": {
								Type:     schema.TypeInt,
								Required: true,
								ForceNew: true,
							},
						},
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
			expected: false,
		},
		{
			name: "equal set",
			resourceSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeSet,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
					ValidateFunc: func(i interface{}, s string) ([]string, []error) {
						return nil, nil
					},
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeSet,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
			expected: true,
		},
		{
			name: "not equal set 1",
			resourceSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeSet,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 2,
					ValidateFunc: func(i interface{}, s string) ([]string, []error) {
						return nil, nil
					},
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeSet,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
		},
		{
			name: "not equal set 2",
			resourceSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeSet,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeString,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
					ValidateFunc: func(i interface{}, s string) ([]string, []error) {
						return nil, nil
					},
				},
			},
			migrationSchema: map[string]*schema.Schema{
				"field": {
					Type:     schema.TypeSet,
					ForceNew: true,
					Elem: schema.Schema{
						Type:     schema.TypeInt,
						ForceNew: true,
						Required: true,
					},
					MaxItems: 100,
					MinItems: 1,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := migrationSchemaComparer.Equal(c.resourceSchema, c.migrationSchema, "")
			if c.expected {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestValidationInsideListElementPrimaryType(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeList,
			ForceNew: true,
			Elem: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					return nil, nil
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
			Elem: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			MaxItems: 100,
			MinItems: 1,
		},
	}
	assert.Nil(t, migrationSchemaComparer.Equal(resourceSchema, migrationSchema, ""))
}

func TestValidationInsideListElementObjectType(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeList,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"embedded_field": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: func(i interface{}, s string) ([]string, []error) {
							return nil, nil
						},
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
	assert.Nil(t, migrationSchemaComparer.Equal(resourceSchema, migrationSchema, ""))
}

func TestValidationInsideListElementObjectType_error_should_contains_correct_field_name(t *testing.T) {
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
	err := migrationSchemaComparer.Equal(resourceSchema, migrationSchema, "")
	assert.NotNil(t, err)
	errMsg := err.Error()
	assert.True(t, strings.Contains(errMsg, "field.embedded_field"))
}

func TestResourceSchemaWithValidationAndMigrationSchemaWithoutValidation(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: func(i interface{}, s string) ([]string, []error) {
				return nil, nil
			},
		},
	}

	migrationSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	assert.Nil(t, migrationSchemaComparer.Equal(resourceSchema, migrationSchema, ""))
}

func TestCleanObjectValidation(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: func(i interface{}, s string) ([]string, []error) {
				return nil, nil
			},
		},
	}
	migrationSchemaComparer.CleanObjectValidations(resourceSchema)

	s := resourceSchema["field"]
	assert.Nil(t, s.ValidateFunc)
}

func TestCleanSchemaValidation(t *testing.T) {
	s := &schema.Schema{
		Type: schema.TypeInt,
		ValidateFunc: func(i interface{}, s string) ([]string, []error) {
			return nil, nil
		},
	}
	migrationSchemaComparer.CleanSchemaValidation(s)
	assert.Nil(t, s.ValidateFunc)
}

func TestCleanListPrimaryElementValidation(t *testing.T) {
	s := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
			ValidateFunc: func(i interface{}, s string) ([]string, []error) {
				return nil, nil
			},
		},
	}
	migrationSchemaComparer.CleanSchemaValidation(s)
	elementSchema := s.Elem.(*schema.Schema)
	assert.Nil(t, elementSchema.ValidateFunc)
}

func TestSchemaFuncShouldBeWipeOut(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
			ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
				panic("")
			},
		},
	}
	migrationSchema := map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				panic("")
			},
		},
	}
	assert.Nil(t, migrationSchemaComparer.Equal(resourceSchema, migrationSchema, ""))
}

func TestSchemaFuncShouldBeWipeOut_same_func(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"field": {
			Type:             schema.TypeInt,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppressFunc,
		},
	}
	migrationSchema := map[string]*schema.Schema{
		"field": {
			Type:             schema.TypeInt,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppressFunc,
		},
	}
	assert.Nil(t, migrationSchemaComparer.Equal(resourceSchema, migrationSchema, ""))
}

var suppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
	panic("")
}
