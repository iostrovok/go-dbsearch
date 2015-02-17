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
	findArrayReg   = regexp.MustCompile(`^\[\].+`)
	findNoTimeReg  = regexp.MustCompile(`^(date|time|timestamp)$`)
	findBoolReg    = regexp.MustCompile(`^(boolean|bool)$`)
	findByteaReg   = regexp.MustCompile(`^(bytea)$`)
	findIntReg     = regexp.MustCompile(`^(bigint|smallint|int|integer|serial|bigserial|smallserial)$`)
	findJSONReg    = regexp.MustCompile(`^(json|jsonb)$`)
	findNumericReg = regexp.MustCompile(`^(numeric|decimal|money)$`)
	findRealReg    = regexp.MustCompile(`^(real|double)$`)
	findTextReg    = regexp.MustCompile(`^(varchar|char|text)$`)
	findTimeReg    = regexp.MustCompile(`^(date|time|timestamp)$`)
)

// An action transitions stochastically to a resulting score.
type convertData func(data interface{}, val reflect.Value) error

func (oRow OneRow) debugV(fieldName string, fTType reflect.Type, data interface{}) {
	if oRow.Log {
		log.Printf("ACTION on %s. oRow.Type: %s => fTType: %s for %#v\n",
			fieldName, oRow.Type, fTType, data)
	}
}

func (aRows *AllRows) errorConvertRunTimeMessage(fieldName string, Type string,
	fTType reflect.Type, err error, data interface{}) error {

	msg := "\n"
	i := 5
	for i > 0 {
		i--
		_, file, line, _ := runtime.Caller(i)
		msg += fmt.Sprintf("Error in %s line %d\n", file, line)
	}

	return fmt.Errorf(msg+"\nfor %s.%s convert data from '%s' to '%s'.\n"+
		"Error System: '%s'\n"+
		"Date: %#v\n",
		aRows.SType, fieldName, Type, fTType, err, data)
}

func (aRows *AllRows) panicInitMessage(what, fieldName, Type string) string {
	return fmt.Sprintf("Error for %s.%s. Not found '%s' for '%s'\n", aRows.SType, fieldName, what, Type)
}

func (aRows *AllRows) errorConvertMessage(fieldName string, Type string, fTType reflect.Type) error {
	_, file, line, _ := runtime.Caller(3)
	_, file0, line0, _ := runtime.Caller(2)
	return fmt.Errorf("Error in %s line %d => in %s line %d for %s.%s convert from '%s' to '%s'\n",
		file0, line0, file, line, aRows.SType, fieldName, Type, fTType)
}

/* Select function */
func (aRows *AllRows) convertSelect(oRow OneRow, fieldName string,
	fTType reflect.Type) (convertData, error) {

	fStrType := fmt.Sprintf("%s", fTType)

	if aRows.Log > 1 {
		oRow.Log = true
	}

	if f, b, err := aRows.convertFuncSlice(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceInt(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceText(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceNum(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceReal(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceBool(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceJSON(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceBytea(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}
	if f, b, err := aRows.convertFuncNoSliceDateTime(oRow, fStrType, fieldName, fTType); b {
		return f, nil
	} else if err != nil {
		return nil, err
	}

	return nil, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
}

/* Check nil */
func isNotNil(data interface{}, field reflect.Value, fieldTypeType reflect.Type) bool {
	if data == nil {
		field.Set(reflect.Zero(fieldTypeType))
		return false
	}
	return true
}

/*
	Check INT value.
	Full error message
*/
func (aRows *AllRows) isInt(v reflect.Value, oRow OneRow, fStrType, fieldName string) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return nil
	}

	tab := ""
	if aRows.TableInfo != nil {
		tab = fmt.Sprintf("Table: %s.%s. ", aRows.Schema, aRows.Table)
	}

	info := fmt.Sprintf(tab+" Structure: %s. Field: %s [%s]. ",
		aRows.SType, oRow.Name, oRow.FType)

	err := fmt.Sprintf("Data from DB: %s, expected: %s.", v.Kind(), oRow.FType)

	return fmt.Errorf("Error convert data. " + info + " " + err)
}

/* IS ARRAYS or SLICE */
func (aRows *AllRows) convertFuncSlice(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData

	if findArrayReg.FindString(oRow.Type) == "" || findNoTimeReg.FindString(oRow.Type) != "" {
		return fn, false, nil
	}

	if fStrType == "[]int" || fStrType == "[]int64" {
		switch oRow.Type {
		case "[]boolean", "[]bool":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
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
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
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
		return fn, true, nil
	}

	if fStrType == "[]uint" || fStrType == "[]byte" || fStrType == "[]uint8" || fStrType == "[]uint64" {
		switch oRow.Type {
		case "[]boolean", "[]bool":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
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
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
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
		return fn, true, nil
	}

	if fStrType == "[]string" {
		switch oRow.Type {
		case "[]boolean", "[]bool":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
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
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
					a := parseArray(_AnyToString(data))
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		}
		return fn, true, nil
	}

	if fStrType == "[]float32" || fStrType == "[]float64" {
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
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
		return fn, true, nil
	}

	if fStrType == "[]bool" {
		switch oRow.Type {
		case "[]boolean", "[]bool":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
					a := parseBoolArrayForBool(data)
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		case "[]text", "[]varchar", "[]char":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
					a := parseBoolArrayForString(data)
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		case "[]bigint", "[]smallint", "[]int", "[]integer":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
					a := parseBoolArrayForNumber(data)
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		case "[]money", "[]numeric", "[]decimal", "[]real", "[]double":
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
					a := parseBoolArrayForReal(data)
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		default:
			fn = func(data interface{}, field reflect.Value) error {
				oRow.debugV(fieldName, fTType, data)
				if isNotNil(data, field, fTType) {
					a := parseBoolArrayForBool(data)
					field.Set(reflect.ValueOf(a))
				}
				return nil
			}
		}
		return fn, true, nil
	}

	return fn, false, nil
}

/* IS --NOT--- ARRAYS or SLICE */
func (aRows *AllRows) convertFuncNoSliceInt(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {
	var fn convertData

	if findIntReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				v := reflect.ValueOf(data)
				if err := aRows.isInt(v, oRow, fStrType, fieldName); err != nil {
					return err
				}
				field.SetInt(v.Int())
			}
			return nil
		}
	case "[]int":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]int{int(data.(int64))}))
			}
			return nil
		}
	case "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]int64{data.(int64)}))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.SetUint(uint64(reflect.ValueOf(data).Int()))
			}
			return nil
		}
	case "[]uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]uint64{data.(uint64)}))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.SetFloat(float64(reflect.ValueOf(data).Int()))
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]float64{float64(reflect.ValueOf(data).Int())}))
			}
			return nil
		}
	default:
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}
	return fn, true, nil
}

func (aRows *AllRows) convertFuncNoSliceText(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData
	if findTextReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				f, err := strconv.ParseFloat(_AnyToString(data), field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.SetFloat(f)
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				f, err := strconv.ParseFloat(_AnyToString(data), 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.Set(reflect.ValueOf([]float64{f}))
			}
			return nil
		}
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.SetInt(i)
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				} else if i < 0 {
					// TODO Make "WEEK" mode when we set i = 0 when i < 0
					err := fmt.Errorf("Value from DB is less then zero")
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.SetUint(uint64(i))
			}
			return nil
		}
	case "[]int":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.Set(reflect.ValueOf([]int{int(i)}))
			}
			return nil
		}
	case "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				i, err := strconv.ParseInt(_AnyToString(data), 0, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.Set(reflect.ValueOf([]int64{i}))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			//sssssssss
			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	default:
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}
	return fn, true, nil
}

func (aRows *AllRows) convertFuncNoSliceNum(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData
	if findNumericReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}

				field.SetInt(int64(i))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}

				field.SetUint(uint64(i))
			}
			return nil
		}

	case "[]int", "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				f, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.SetFloat(float64(f))
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}
				field.Set(reflect.ValueOf([]float64{f}))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	default:
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}

	return fn, true, nil
}

func (aRows *AllRows) convertFuncNoSliceReal(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData
	if findRealReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.SetInt(int64(reflect.ValueOf(data).Float()))
			}
			return nil
		}
	case "[]int", "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.SetFloat(reflect.ValueOf(data).Float())
			}
			return nil
		}
	case "[]float32", "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]float64{reflect.ValueOf(data).Float()}))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	case "[]string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]string{_AnyToString(data)}))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				s := noNumberDots.ReplaceAllString(_AnyToString(data), "")
				i, err := strconv.ParseFloat(s, field.Type().Bits())
				if err != nil {
					// TODO Make "WEEK" mode when we set "" in error case
					return aRows.errorConvertRunTimeMessage(fieldName, oRow.Type, fTType, err, data)
				}

				field.SetUint(uint64(i))
			}
			return nil
		}
	default:
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}

	return fn, true, nil
}

func (aRows *AllRows) convertFuncNoSliceBool(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData
	if findBoolReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "int", "int8", "int16", "int32", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				i := 0
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]int{i}))
			}
			return nil
		}
	case "[]int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				var i int64
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]int64{i}))
			}
			return nil
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				var i uint
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]uint{i}))
			}
			return nil
		}
	case "[]uint64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				var i uint64
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]uint64{i}))
			}
			return nil
		}
	case "[]byte", "[]uint8":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				var i uint8
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]uint8{i}))
			}
			return nil
		}
	case "[]float32":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				var i float32
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]float32{i}))
			}
			return nil
		}
	case "[]float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				var i float64
				if reflect.ValueOf(data).Bool() {
					i = 1
				}
				field.Set(reflect.ValueOf([]float64{i}))
			}
			return nil
		}
	case "float32", "float64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				var i string
				if reflect.ValueOf(data).Bool() {
					i = "1"
				}
				field.Set(reflect.ValueOf([]string{i}))
			}
			return nil
		}
	case "bool", "boolean":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.SetBool(reflect.ValueOf(data).Bool())
			}
			return nil
		}
	case "[]bool", "[]boolean":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				field.Set(reflect.ValueOf([]bool{reflect.ValueOf(data).Bool()}))
			}
			return nil
		}
	default:
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}
	return fn, true, nil
}

func (aRows *AllRows) convertFuncNoSliceJSON(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData
	if findJSONReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "map[string]interface {}":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
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
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	default:
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}
	return fn, true, nil
}

func (aRows *AllRows) convertFuncNoSliceBytea(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData
	if findByteaReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "[]byte", "[]uint8":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				line := data.([]byte)
				field.Set(reflect.ValueOf(line))
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)

			if isNotNil(data, field, fTType) {
				field.SetString(_AnyToString(data))
			}
			return nil
		}
	default:
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}

	return fn, true, nil
}

func (oRow OneRow) prepareNoSliceDateTime(data interface{}, fieldName string) (time.Time, error) {
	var t time.Time
	var e error
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

func (aRows *AllRows) convertFuncNoSliceDateTime(oRow OneRow, fStrType, fieldName string,
	fTType reflect.Type) (convertData, bool, error) {

	var fn convertData
	if findTimeReg.FindString(oRow.Type) == "" {
		return fn, false, nil
	}

	switch fStrType {
	case "time.Time":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				t, e := oRow.prepareNoSliceDateTime(data, fieldName)
				if e != nil {
					return e
				}
				field.Set(reflect.ValueOf(t))
			}
			return nil
		}
	case "int", "int64":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				t, e := oRow.prepareNoSliceDateTime(data, fieldName)
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				t, e := oRow.prepareNoSliceDateTime(data, fieldName)
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				t, e := oRow.prepareNoSliceDateTime(data, fieldName)
				if e != nil {
					return e
				}
				field.SetFloat(float64(t.UnixNano()) / 1000)
			}
			return nil
		}
	case "string":
		fn = func(data interface{}, field reflect.Value) error {
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				t, e := oRow.prepareNoSliceDateTime(data, fieldName)
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				t, e := oRow.prepareNoSliceDateTime(data, fieldName)
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
			oRow.debugV(fieldName, fTType, data)
			if isNotNil(data, field, fTType) {
				t, e := oRow.prepareNoSliceDateTime(data, fieldName)
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
		return nil, false, aRows.errorConvertMessage(fieldName, oRow.Type, fTType)
	}

	return fn, true, nil
}
