package migrationSchemaComparer

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Equal(resourceSchemas, migrationSchemas map[string]*schema.Schema, currentPath string) error {
	if len(resourceSchemas) != len(migrationSchemas) {
		return fmt.Errorf("different count, path:%s, resource:%+v, migration:%+v", currentPath, resourceSchemas, migrationSchemas)
	}
	CleanObjectValidations(resourceSchemas)
	CleanObjectValidations(migrationSchemas)
	for key, resourceSchema := range resourceSchemas {
		migrationSchema, ok := migrationSchemas[key]
		if !ok {
			return fmt.Errorf("expected %s.%s not existd in migration", currentPath, key)
		}
		if !reflect.DeepEqual(resourceSchema, migrationSchema) {
			if resourceSchema.Type == schema.TypeList || resourceSchema.Type == schema.TypeSet {
				switch t := resourceSchema.Elem.(type) {
				case *schema.Resource:
					return Equal(t.Schema, migrationSchema.Elem.(*schema.Resource).Schema, fmt.Sprintf("%s.%s", currentPath, key))
				}
			}
			return fmt.Errorf("%s.%s different, resource:%+v, migration:%+v", currentPath, key, *resourceSchema, *migrationSchema)
		}
	}
	return nil
}

func CleanObjectValidations(object map[string]*schema.Schema) {
	for _, s := range object {
		CleanSchemaValidation(s)
	}
}

func CleanSchemaValidation(s *schema.Schema) {
	setToDefault(s)
	if s.Type == schema.TypeList || s.Type == schema.TypeSet || s.Type == schema.TypeMap {
		switch t := s.Elem.(type) {
		case *schema.Schema:
			CleanSchemaValidation(t)
		case *schema.Resource:
			CleanObjectValidations(t.Schema)
		}
	}
}

func setToDefault(s *schema.Schema) {
	s.ValidateFunc = nil
	s.DiffSuppressFunc = nil
	s.StateFunc = nil
	s.ValidateDiagFunc = nil
	s.DefaultFunc = nil
	s.Default = nil
	s.Description = ""
	s.ConfigMode = 0
	s.InputDefault = ""
	s.ForceNew = false
	s.MaxItems = 0
	s.MinItems = 0
	s.ComputedWhen = nil
	s.ConflictsWith = nil
	s.ExactlyOneOf = nil
	s.AtLeastOneOf = nil
	s.RequiredWith = nil
	s.Deprecated = ""
	s.Sensitive = false
}
