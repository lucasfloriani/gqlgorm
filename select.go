package gqlgorm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/azer/snakecase"
)

// GetQueryFields returns a list of field by a GraphQL Type name
func GetQueryFields(ctx context.Context, name string) (fields []string) {
	for _, field := range graphql.CollectFieldsCtx(ctx, []string{name}) {
		fields = append(fields, snakecase.SnakeCase(field.Name))
	}
	return
}

// ConvertToSelectFields converts fields array to select columns string
func ConvertToSelectFields(fields []string, prefix string, obj interface{}) string {
	return strings.Join(deepFields(obj, fields, 0), fmt.Sprintf(", %s.", prefix))
}

func deepFields(obj interface{}, fields []string, level int) (allowedFields []string) {
	s := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < s.NumField(); i++ {
		typeField := t.Field(i)
		valueField := s.Field(i)
		tag := typeField.Tag

		if IsType(tag, SkipTag) {
			continue
		}

		if valueField.Kind() == reflect.Ptr && typeField.Type.Elem().Kind() == reflect.Struct && valueField.IsNil() {
			ptrToZero := reflect.New(typeField.Type.Elem())
			valueField = ptrToZero.Elem()
		}

		fieldName := snakecase.SnakeCase(typeField.Name)
		if IsType(tag, EmbeddedFilter) && contain(fields, fieldName) {
			allowedFields = append(allowedFields, deepFields(valueField.Interface(), fields, level+1)...)
			continue
		}

		names := GetAlias(tag)
		names = append(names, cleanString(fieldName))
		if level != 0 || contains(fields, names) {
			allowedFields = append(allowedFields, fieldName)
		}
	}

	return
}

func cleanString(name string) string {
	return strings.TrimSuffix(name, "_id")
}

func contains(s []string, e []string) bool {
	for _, a := range e {
		if contain(s, a) {
			return true
		}
	}
	return false
}

func contain(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
