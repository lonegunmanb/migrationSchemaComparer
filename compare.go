package migrationSchemaComparer

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
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
			if resourceSchema.Type == schema.TypeList || resourceSchema.Type == schema.TypeSet || resourceSchema.Type == schema.TypeMap {
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
	s.ValidateFunc = nil
	s.DiffSuppressFunc = nil
	s.StateFunc = nil
	s.ValidateDiagFunc = nil
	s.DefaultFunc = nil
	if s.Type == schema.TypeList || s.Type == schema.TypeSet || s.Type == schema.TypeMap {
		switch t := s.Elem.(type) {
		case *schema.Schema:
			CleanSchemaValidation(t)
		case *schema.Resource:
			CleanObjectValidations(t.Schema)
		}
	}
}