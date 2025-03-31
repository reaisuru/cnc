package refhelper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Set(structure interface{}, variableNameIgnoreCase, value string) error {
	v := reflect.ValueOf(structure).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.EqualFold(field.Name, variableNameIgnoreCase) {
			fieldValue := v.FieldByName(field.Name)
			if !fieldValue.CanSet() {
				return fmt.Errorf("cannot set field %s", field.Name)
			}

			return SetField(fieldValue, field.Type, value)
		}
	}

	return fmt.Errorf("no such field: %s", variableNameIgnoreCase)
}

func SetField(field reflect.Value, fieldType reflect.Type, value string) error {
	switch fieldType.Kind() {
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intValue))
	case reflect.String:
		field.SetString(value)
	case reflect.Slice:
		if fieldType.Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(strings.Split(value, ",")))
		}
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			timeValue, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(timeValue))
		}
	default:
		return fmt.Errorf("unsupported field type: %s", fieldType)
	}
	return nil
}
