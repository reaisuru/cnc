package packages

import (
	"encoding/json"
	"reflect"
)

// JsonEncode is ported into the JSONFunctions map
func JsonEncode(encode any) string {
	value := reflect.ValueOf(encode)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	/* base data types are handled here. */
	if value.Kind() != reflect.Map && value.Kind() != reflect.Struct && value.Kind() != reflect.Func {
		return jsonEncodeValue(encode)
	}

	reference := make(map[string]any)

	for i := 0; i < value.NumField(); i++ {
		index := value.Field(i)
		if index.Kind() == reflect.Pointer {
			index = index.Elem()
		}

		if index.Kind() == reflect.Func {
			continue
		}

		reference[reflect.TypeOf(encode).Field(i).Name] = index.Interface()
		if tag := reflect.TypeOf(encode).Field(i).Tag.Get("swash"); len(tag) > 0 {
			delete(reference, reflect.TypeOf(encode).Field(i).Name)
			reference[tag] = index.Interface()
		}
	}

	return jsonEncodeValue(reference)
}

func jsonEncodeValue(value any) string {
	content, err := json.Marshal(value)
	if err != nil {
		return err.Error()
	}

	return string(content)
}
