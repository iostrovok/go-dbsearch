package dbsearch

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func convert(List map[string]*OneRow, name string, data interface{}, Result interface{}) error {

	val := reflect.ValueOf(Result).Elem()

	if data == nil {
		val.Set(reflect.Zero(val.Type()))
		return nil
	}

	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if !dataVal.IsValid() {
		val.Set(reflect.Zero(val.Type()))
		return nil
	}

	dataValKind := dataVal.Kind()
	if dataValKind != reflect.Map {
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataValKind)
	}

	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return fmt.Errorf(
			"'%s' needs a map with string keys, has '%s' keys",
			name, dataValType.Key().Kind())
	}

	fields := make(map[*reflect.StructField]reflect.Value)
	structType := val.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)

		if fieldType.Anonymous {
			return fmt.Errorf("%s: unsupported type: %s", fieldType.Name, fieldType.Type.Kind())
		}

		fields[&fieldType] = val.Field(i)
		if !fields[&fieldType].IsValid() {
			panic("field is not valid")
		}
	}

	for fieldType, field := range fields {
		fieldName := fieldType.Name

		el := List[fieldName]
		if el == nil {
			log.Fatalf("No setup fieldName [%s]\n", fieldName)
		}

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)

		if field.CanSet() && rawMapVal != reflect.Zero(reflect.TypeOf(rawMapVal)).Interface() {
			if err := el.SetFunc(rawMapVal.Interface(), field); err != nil {
				return err
			}
		}
	}

	return nil
}

func _AnyToString(data interface{}) string {

	dataVal := reflect.ValueOf(data)

	out := ""
	switch dataVal.Kind() {
	case reflect.String:
		out = dataVal.String()
	case reflect.Bool:
		if dataVal.Bool() {
			out = "1"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = strconv.FormatInt(dataVal.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = strconv.FormatUint(dataVal.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		out = strconv.FormatFloat(dataVal.Float(), 'f', -1, 64)
	case reflect.Slice:
		// TODO : Bad solution
		if dataVal.Type().Elem().Kind() == reflect.Uint8 {
			out = string(dataVal.Interface().([]uint8))
		}
	}

	return out
}
