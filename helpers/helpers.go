package helpers

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func MergeStructs(a, b interface{}) interface{} {
	// Ensure that both `a` and `b` are of the same type and are structs.
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	if aValue.Kind() != reflect.Struct || bValue.Kind() != reflect.Struct || aValue.Type() != bValue.Type() {
		return a
	}

	merged := reflect.New(aValue.Type()).Elem()

	// Iterate through the fields of `a` and copy them to `merged`.
	for i := 0; i < aValue.NumField(); i++ {
		fieldA := aValue.Field(i)
		fieldName := aValue.Type().Field(i).Name
		fieldB := bValue.FieldByName(fieldName)

		if fieldB.IsValid() && !isEmptyValue(fieldB) {
			merged.Field(i).Set(fieldB)
		} else {
			merged.Field(i).Set(fieldA)
		}
	}

	return merged.Interface()
}

// isEmptyValue checks if a value is empty.
func isEmptyValue(value reflect.Value) bool {
	zero := reflect.Zero(value.Type())
	return value.Interface() == zero.Interface()
}

type Detail struct {
	Rule    string
	Message string
}

func ParseErrors(errors validator.ValidationErrors) map[string]Detail {

	result := make(map[string]Detail)

	for _, verror := range errors {
		result[verror.Field()] = Detail{Rule: verror.Tag(), Message: verror.Error()}
	}

	return result
}
