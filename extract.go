package pocket

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func extractProp(r *http.Request, ft reflect.StructField, fv reflect.Value) error {
	if param := ft.Tag.Get("param"); param != "" {
		return coerceIntoValue(r.URL.Query().Get(param), ft, fv)
	}

	return fmt.Errorf("no target for field %s", ft.Name)
}

func coerceIntoValue(s string, ft reflect.StructField, fv reflect.Value) error {
	switch ft.Type.Kind() {
	case reflect.Bool:
		value, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		fv.SetBool(value)

	case reflect.Slice:
		switch fv.Type().Elem().Kind() {
		case reflect.Uint8:
			fv.SetBytes([]byte(s))
		}

	case reflect.Float32:
		value, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		fv.SetFloat(value)

	case reflect.Float64:
		value, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(value)

	case reflect.Int:
		value, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		fv.SetInt(int64(value))

	case reflect.Uint:
		value, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(value)

	case reflect.String:
		fv.SetString(s)

	case reflect.Struct:
		return fmt.Errorf("cannot handle struct type %s", ft.Type.Name())

	default:
		return fmt.Errorf("cannot handle type %v", ft.Type.Kind())
	}
	return nil
}
