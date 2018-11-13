package gqlgorm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/azer/snakecase"
	"github.com/jinzhu/gorm"
)

const (
	// LessThen suffix string
	LessThen = "Lt"
	// LessThenEqual suffix string
	LessThenEqual = "Lte"
	// GreaterThen suffix string
	GreaterThen = "Gt"
	// GreaterThenEqual suffix string
	GreaterThenEqual = "Gte"
)

// FilterByObject filters all fields from object
// is used to query filter objects from gqlgen
func FilterByObject(tx *gorm.DB, obj interface{}) *gorm.DB {
	s := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < s.NumField(); i++ {
		tag := t.Field(i).Tag
		if IsType(tag, SkipTag) {
			continue
		}

		field := s.Field(i)
		if field.Interface() == reflect.Zero(field.Type()).Interface() {
			continue
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		tx = queryByFieldType(tx, field.Kind(), field.Interface(), t.Field(i).Name)
	}

	return tx
}

func queryByFieldType(tx *gorm.DB, typeArg reflect.Kind, value interface{}, name string) *gorm.DB {
	switch typeArg {
	case reflect.String:
		query := fmt.Sprintf("\"%s\" LIKE ?", normalizeFieldName(name))
		tx = tx.Where(query, "%"+value.(string)+"%")
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Struct:
		query := fmt.Sprintf("\"%s\" %s ?", normalizeFieldName(name), getComparisonType(name))
		tx = tx.Where(query, value)
	default:
		panic("Invalid type field to query")
	}

	return tx
}

func getComparisonType(name string) string {
	if strings.HasSuffix(name, LessThen) {
		return "<"
	}
	if strings.HasSuffix(name, LessThenEqual) {
		return "<="
	}
	if strings.HasSuffix(name, GreaterThen) {
		return ">"
	}
	if strings.HasSuffix(name, GreaterThenEqual) {
		return ">="
	}
	return "="
}

func normalizeFieldName(name string) string {
	name = strings.TrimSuffix(name, LessThen)
	name = strings.TrimSuffix(name, LessThenEqual)
	name = strings.TrimSuffix(name, GreaterThen)
	name = strings.TrimSuffix(name, GreaterThenEqual)
	return snakecase.SnakeCase(name)
}
