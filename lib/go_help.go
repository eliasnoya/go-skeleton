package lib

import "reflect"

func GetStructName(s interface{}) string {
	t := reflect.TypeOf(s)

	// If it's a pointer, get the type of the underlying element
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	// Check if the type is a struct
	if t.Kind() == reflect.Struct {
		return t.Name()
	}

	return "Not a struct"
}
