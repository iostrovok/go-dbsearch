package dbsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

// construct a regexp to extract values:
var (
	findArrayReg  = regexp.MustCompile(`^\[\].+`)
	findNoTimeReg = regexp.MustCompile(`^(date|time|timestamp)$`)

	findBoolReg    = regexp.MustCompile(`^(boolean|bool)$`)
	findByteaReg   = regexp.MustCompile(`^(bytea)$`)
	findIntReg     = regexp.MustCompile(`^(bigint|smallint|int|integer|serial|bigserial|smallserial)$`)
	findJsonReg    = regexp.MustCompile(`^(json|jsonb)$`)
	findNumericReg = regexp.MustCompile(`^(numeric|decimal|money)$`)
	findRealReg    = regexp.MustCompile(`^(real|double)$`)
	findTextReg    = regexp.MustCompile(`^(varchar|char|text)$`)
	findTimeReg    = regexp.MustCompile(`^(date|time|timestamp)$`)
)

// An action transitions stochastically to a resulting score.
type ConvertData func(data interface{}, val reflect.Value) error

func (oRow OneRow) DebugV(fieldName string, fTType reflect.Type, data interface{}) {
	if oRow.Log {
		log.Printf("ACTION on %s. oRow.Type: %s => fTType: %s for %#v\n",
			fieldName, oRow.Type, fTType, data)
	}
}

func (aRows *AllRows) ErrorConvertRunTimeMessage(fieldName string, Type string,
	fTType reflect.Type, err error, data interface{}) error {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Errorf("Error in %s line %d for %s.%s convert data from '%s' to '%s'.\n"+
		"Error System: '%s'\n"+
		"Date: %#v\n",
		file, line, aRows.SType, fieldName, Type, fTType, err, data)
}

func (aRows *AllRows) PanicInitMessage(what, fieldName, Type string) string {
	return fmt.Sprintf("Error for %s.%s. Not found '%s' for '%s'\n", aRows.SType, fieldName, what, Type)
}

func (aRows *AllRows) PanicConvert(fieldName string, Type string, fTType reflect.Type) {
	log.Panicf("Error for %s.%s convert from '%s' to '%s'\n", aRows.SType, fieldName, Type, fTType)
}

/* Select function */
func (aRows *AllRows) convert_select(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) ConvertData {

	if aRows.Log > 1 {
		oRow.Log = true
	}

	if f, b := aRows.convert_func_slice(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_int(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_text(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_num(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_real(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_bool(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_json(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_bytea(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	if f, b := aRows.convert_func_no_slice_datetime(oRow, fStrType, fieldName, fTType); b {
		return f
	}
	aRows.PanicConvert(fieldName, oRow.Type, fTType)
	return nil
}

func IsNotNilValue(data interface{}, field reflect.Value, fieldTypeType reflect.Type) bool {
	if data == nil {
		field.Set(reflect.Zero(fieldTypeType))
		return false
	}
	return true
}

/* IS ARRAYS or SLICE */
func (aRows *AllRows) convert_func_slice(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData

	if findArrayReg.FindString(oRow.Type) == "" || findNoTimeReg.FindString(oRow.Type) != "" {
		return fn, false
	}

	if fStrType == "[]int" || fStrType == "[]int64" {
		switch oRow.Type {
		case "[]boolean", "[]bool":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.DebugV(fieldName, fTType, data)
				if IsNotNilValue(data, field, fTType) {
					r := parseArray(_AnyToString(data))
					if fStrType == "[]int64" {
						a := make([]int64, len(r))
						for i, v := range r {
							if v == "T" || v == "t" {
								a[i] = 1
							} else {
								a[i] = 0
							}
						}
						field.Set(reflect.ValueOf(a))
					} else {
						a := make([]int, len(r))
						for i, v := range r {
							if v == "T" || v == "t" {
								a[i] = 1
							} else {
								a[i] = 0
							}
						}
						field.Set(reflect.ValueOf(a))
					}
				}
				return nil
			}
		default:
			fn = func(data interface{}, field reflect.Value) error {
				oRow.DebugV(fieldName, fTType, data)
				if IsNotNilValue(data, field, fTType) {
					if fStrType == "[]int64" {
						a := parseInt64Array(_AnyToString(data))
						field.Set(reflect.ValueOf(a))
					} else {
						a := parseIntArray(_AnyToString(data))
						field.Set(reflect.ValueOf(a))
					}
				}
				return nil
			}
		}
		return fn, true
	}

	if fStrType == "[]uint" || fStrType == "[]byte" || fStrType == "[]uint8" || fStrType == "[]uint64" {
		switch oRow.Type {
		case "[]boolean", "[]bool":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.DebugV(fieldName, fTType, data)
				if IsNotNilValue(data, field, fTType) {
					r := parseArray(_AnyToString(data))
					if fStrType == "[]uint64" {
						a := make([]uint64, len(r))
						for i, v := range r {
							if v == "T" || v == "t" {
								a[i] = 1
							} else {
								a[i] = 0
							}
						}
						field.Set(reflect.ValueOf(a))
					} else if fStrType == "[]uint" {
						a := make([]uint, len(r))
						for i, v := range r {
							if v == "T" || v == "t" {
								a[i] = 1
							} else {
								a[i] = 0
							}
						}
						field.Set(reflect.ValueOf(a))
					} else {
						a := make([]uint8, len(r))
						for i, v := range r {
							if v == "T" || v == "t" {
								a[i] = 1
							} else {
								a[i] = 0
							}
						}
						field.Set(reflect.ValueOf(a))
					}
				}
				return nil
			}
		default:
			fn = func(data interface{}, field reflect.Value) error {
				oRow.DebugV(fieldName, fTType, data)
				if IsNotNilValue(data, field, fTType) {
					if fStrType == "[]uint64" {
						a := parseUint64Array(_AnyToString(data))
						field.Set(reflect.ValueOf(a))
					} else if fStrType == "[]uint" {
						a := parseUintArray(_AnyToString(data))
						field.Set(reflect.ValueOf(a))
					} else {
						a := parseUint8Array(_AnyToString(data))
						field.Set(reflect.ValueOf(a))
					}
				}
				return nil
			}
		}
		return fn, true
	}

	if fStrType == "[]string" {
		switch oRow.Type {
		case "[]boolean", "[]bool":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.DebugV(fieldName, fTType, data)
				if IsNotNilValue(data, field, fTType) {
					r := parseArray(_AnyToString(data))
					a := make([]string, len(r))
					for i, v := range r {
						if v == "T" || v == "t" {
							a[i] = "1"
						} else {
							a[i] = ""
						}
					}
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		default:
			fn = func(data interface{}, field reflect.Value) error {
				oRow.DebugV(fieldName, fTType, data)
				if IsNotNilValue(data, field, fTType) {
					a := parseArray(_AnyToString(data))
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		}
		return fn, true
	}

	if fStrType == "[]float32" || fStrType == "[]float64" {
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				if fStrType == "[]float64" {
					a := parseFloat64Array(_AnyToString(data))
					field.Set(reflect.ValueOf(a))
				} else {
					a := parseFloat32Array(_AnyToString(data))
					field.Set(reflect.ValueOf(a))
				}
			}
			return nil
		}
		return fn, true
	}

	if fStrType == "[]bool" {
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				a := parseBoolArray(data)
				field.Set(reflect.ValueOf(a))
			}
			return nil
		}
		return fn, true
	}

	return fn, false
}

/* IS --NOT--- ARRAYS or SLICE */
func (aRows *AllRows) convert_func_no_slice_int(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {
	var fn ConvertData

	if findIntReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	v := []interface{}{"convert_func_no_slice_int"}
	oRow.DebugV(fieldName, fTType, v)
	log.Printf("fStrType: %s\n", fStrType)
	log.Printf("fTType: %s\n", fTType)
	log.Printf("fieldName: %s\n", fieldName)

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.SetInt(reflect.ValueOf(data).Int())
			}
			return nil
		}
	case "[]int":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]int{int(data.(int64))}))
			}
			return nil
		}
	case "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]int64{data.(int64)}))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.SetUint(uint64(reflect.ValueOf(data).Int()))
			}
			return nil
		}
	case "[]uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]uint64{data.(uint64)}))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.SetFloat(float64(reflect.ValueOf(data).Int()))
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]float64{float64(reflect.ValueOf(data).Int())}))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}
	return fn, true
}

func (aRows *AllRows) convert_func_no_slice_text(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData
	if findTextReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	v := []interface{}{"convert_func_no_slice_text"}
	oRow.DebugV(fieldName, fTType, v)
	log.Printf("fStrType: %s\n", fStrType)
	log.Printf("fTType: %s\n", fTType)
	log.Printf("fieldName: %s\n", fieldName)

	switch fStrType {
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				f, err := strconv.ParseFloat(_AnyToString(data), field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.SetFloat(f)
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				f, err := strconv.ParseFloat(_AnyToString(data), 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.Set(reflect.ValueOf([]float64{f}))
			}
			return nil
		}
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				} else {
					field.SetInt(i)
				}
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				} else if i < 0 {
					// TODO Make "WEEK" mode when we set i = 0 when i < 0
					err := fmt.Errorf("Value from DB is less then zero")
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				} else {
					field.SetUint(uint64(i))
				}
			}
			return nil
		}
	case "[]int":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				} else {
					field.Set(reflect.ValueOf([]int{int(i)}))
				}
			}
			return nil
		}
	case "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				} else {
					field.Set(reflect.ValueOf([]int64{i}))
				}
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			//sssssssss
			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}
	return fn, true
}

func (aRows *AllRows) convert_func_no_slice_num(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData
	if findNumericReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}

				field.SetInt(int64(i))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}

				field.SetUint(uint64(i))
			}
			return nil
		}

	case "[]int", "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				if fStrType == "[]int" {
					field.Set(reflect.ValueOf([]int{int(i)}))
				} else {
					field.Set(reflect.ValueOf([]int64{int64(i)}))
				}
			}
			return nil
		}
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				f, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.SetFloat(float64(f))
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.Set(reflect.ValueOf([]float64{f}))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}

	return fn, true
}

func (aRows *AllRows) convert_func_no_slice_real(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData
	if findRealReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.SetInt(int64(reflect.ValueOf(data).Float()))
			}
			return nil
		}
	case "[]int", "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				if fStrType == "[]int" {
					field.Set(reflect.ValueOf([]int{int(data.(float64))}))
				} else {
					field.Set(reflect.ValueOf([]int64{int64(data.(float64))}))
				}
			}
			return nil
		}
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.SetFloat(reflect.ValueOf(data).Float())
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]float64{reflect.ValueOf(data).Float()}))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.ErrorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}

				field.SetUint(uint64(i))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}

	return fn, true
}

func (aRows *AllRows) convert_func_no_slice_bool(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData
	if findBoolReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				if reflect.ValueOf(data).Bool() {
					field.SetInt(1)
				} else {
					field.SetInt(0)
				}
			}
			return nil
		}
	case "[]int":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i int = 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]int{i}))
			}
			return nil
		}
	case "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i int64 = 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]int64{i}))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				if reflect.ValueOf(data).Bool() {
					field.SetUint(1)
				} else {
					field.SetUint(0)
				}
			}
			return nil
		}
	case "[]uint":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i uint = 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]uint{i}))
			}
			return nil
		}
	case "[]uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i uint64 = 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]uint64{i}))
			}
			return nil
		}
	case "[]byte", "[]uint8":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i uint8 = 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]uint8{i}))
			}
			return nil
		}
	case "[]float32":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i float32 = 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]float32{i}))
			}
			return nil
		}
	case "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i float64 = 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]float64{i}))
			}
			return nil
		}
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				if reflect.ValueOf(data).Bool() {
					field.SetFloat(1)
				} else {
					field.SetFloat(0)
				}
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				if reflect.ValueOf(data).Bool() {
					field.SetString("1")
				} else {
					field.SetString("")
				}
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				var i string = ""
				if reflect.ValueOf(data).Bool() {
					i = "1"
				}
				field.Set(reflect.ValueOf([]string{i}))
			}
			return nil
		}
	case "bool", "boolean":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.SetBool(reflect.ValueOf(data).Bool())
			}
			return nil
		}
	case "[]bool", "[]boolean":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				field.Set(reflect.ValueOf([]bool{reflect.ValueOf(data).Bool()}))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}
	return fn, true
}

func (aRows *AllRows) convert_func_no_slice_json(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData
	if findJsonReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	switch fStrType {
	case "map[string]interface {}":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				line := _AnyToString(data)
				if line == "" {
					line = "{}"
				}
				raw := []byte(line)
				var res map[string]interface{}
				err := json.Unmarshal(raw, &res)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(res))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}
	return fn, true
}

func (aRows *AllRows) convert_func_no_slice_bytea(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData
	if findByteaReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	switch fStrType {
	case "[]byte", "[]uint8":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				line := data.([]byte)
				field.Set(reflect.ValueOf(line))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)

			if IsNotNilValue(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}

	return fn, true
}

func (oRow OneRow) prepare_no_slice_datetime(data interface{}, fieldName string) (time.Time, error) {
	var t time.Time
	var e error = nil
	switch ch := data.(type) {
	case time.Time:
		t = data.(time.Time)
		if oRow.Type == "date" {
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
	default:
		e = fmt.Errorf("%s: unsupported type from DB: %s\n", fieldName, ch)
	}
	return t, e
}

func (aRows *AllRows) convert_func_no_slice_datetime(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (ConvertData, bool) {

	var fn ConvertData
	if findTimeReg.FindString(oRow.Type) == "" {
		return fn, false
	}

	switch fStrType {
	case "time.Time":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				t, e := oRow.prepare_no_slice_datetime(data, fieldName)
				if e != nil {
					return e
				}
				field.Set(reflect.ValueOf(t))
			}
			return nil
		}
	case "int", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				t, e := oRow.prepare_no_slice_datetime(data, fieldName)
				if e != nil {
					return e
				}
				i := t.Unix()

				if fStrType == "int" {
					field.Set(reflect.ValueOf(int(i)))
				} else {
					field.Set(reflect.ValueOf(int64(i)))
				}
			}
			return nil
		}
	case "uint", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				t, e := oRow.prepare_no_slice_datetime(data, fieldName)
				if e != nil {
					return e
				}
				i := t.Unix()
				if i < 0 {
					i = 0
				}
				if fStrType == "uint" {
					field.Set(reflect.ValueOf(uint(i)))
				} else {
					field.Set(reflect.ValueOf(uint64(i)))
				}
			}
			return nil
		}
	case "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				t, e := oRow.prepare_no_slice_datetime(data, fieldName)
				if e != nil {
					return e
				}
				field.SetFloat(float64(t.UnixNano()) / 1000)
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				t, e := oRow.prepare_no_slice_datetime(data, fieldName)
				if e != nil {
					return e
				}
				layout := "2006-01-02 15:04:05 -0700"
				field.SetString(string(t.Format(layout)))
			}
			return nil
		}
	case "map[string]int", "map[string]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				t, e := oRow.prepare_no_slice_datetime(data, fieldName)
				if e != nil {
					return e
				}
				_, i := t.Zone()
				a := map[string]int{
					"year": t.Year(), "month": int(t.Month()), "day": t.Day(),
					"hour": t.Hour(), "minute": t.Minute(), "second": t.Second(),
					"nanosecond": t.Nanosecond(), "zone": i,
				}
				field.Set(reflect.ValueOf(a))
			}
			return nil
		}
	case "[]int", "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.DebugV(fieldName, fTType, data)
			if IsNotNilValue(data, field, fTType) {
				t, e := oRow.prepare_no_slice_datetime(data, fieldName)
				if e != nil {
					return e
				}
				_, i := t.Zone()
				a := []int{
					t.Year(), int(t.Month()), t.Day(),
					t.Hour(), t.Minute(), t.Second(),
					t.Nanosecond(), i,
				}
				field.Set(reflect.ValueOf(a))
			}
			return nil
		}
	default:
		aRows.PanicConvert(fieldName, oRow.Type, fTType)
	}

	return fn, true
}
