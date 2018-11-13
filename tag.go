package gqlgorm

import (
	"reflect"
	"strings"
)

type filterType string

const (
	// EmbeddedFilter filter embedded fields in struct
	EmbeddedFilter filterType = "embedded"
	// SkipTag is the value from FilterTag that skips a field in struct
	SkipTag filterType = "skip"
)

func (ft filterType) String() string {
	return string(ft)
}

func splitTagValues(tag string) map[string]interface{} {
	mapValues := make(map[string]interface{})
	if tag == "" {
		return mapValues
	}
	typeValues := strings.Split(tag, ";")

	for _, typeValue := range typeValues {
		typeMap := strings.Split(typeValue, ":")
		if typeMap[0] == "alias" {
			mapValues[typeMap[0]] = strings.Split(typeMap[1], ",")
			continue
		}
		mapValues[typeMap[0]] = typeMap[1]
	}

	return mapValues
}

// IsType check if tag string contains the type from argument
func IsType(tag reflect.StructTag, typeArg filterType) bool {
	text, ok := tag.Lookup("filter")
	if !ok || text == "" {
		return false
	}

	values := splitTagValues(text)

	if val, ok := values["type"]; ok {
		if val.(string) == typeArg.String() {
			return true
		}
	}

	return false
}

// GetAlias from tag object
func GetAlias(tag reflect.StructTag) (alias []string) {
	text, ok := tag.Lookup("filter")
	if !ok || text == "" {
		return
	}

	values := splitTagValues(text)
	if val, ok := values["alias"]; ok {
		return val.([]string)
	}

	return
}
