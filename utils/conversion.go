package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// struct to map
func Decode(structPointer interface{}) map[string]interface{} {
	t := reflect.TypeOf(structPointer).Elem()
	v := reflect.ValueOf(structPointer).Elem()

	mapParameter := make(map[string]interface{}, 0)

	for i := 0; i < t.NumField(); i++ {
		mapParameter[t.Field(i).Name] = v.Field(i).Interface()
	}
	return mapParameter
}

func Encode_2(mapParameter interface{}, structPointer interface{}) {
	mapV := reflect.ValueOf(mapParameter)
	fmt.Println(mapV.Kind())
}

// map to struct
func Encode(mapParameter interface{}, structPointer interface{}) {
	mapV := reflect.Indirect(reflect.ValueOf(mapParameter))
	if mapV.Kind() != reflect.Map || mapV.Type().Key().Kind() != reflect.String {
		fmt.Println("unconvertible: mapParameter must be a map with string keys")
		return
	}

	strP := reflect.ValueOf(structPointer)
	if strP.Kind() != reflect.Ptr {
		fmt.Println("unconvertible: structPointer must be a pointer")
		return
	}
	strV := strP.Elem()
	if !strV.CanAddr() {
		fmt.Println("unconvertible: structPointer must be addressable")
		return
	}
	strT := reflect.TypeOf(structPointer).Elem()

	for i := 0; i < strT.NumField(); i++ {
		switch kind(strV.Field(i)) {
		case reflect.Int:
			setInt(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.Uint:
			setUint(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.Bool:
			setBool(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.String:
			setString(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.Float64:
			setFloat(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.Map:
			setMap(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.Slice:
			setSlice(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.Struct:
			setStruct(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		case reflect.Interface:
			setInterface(mapValue(mapV, strT.Field(i).Name), strV.FieldByName(strT.Field(i).Name))
		}
	}
}

func mapValue(mapParameter reflect.Value, structField string) interface{} {
	for _, v := range mapParameter.MapKeys() {
		if strings.ToLower(v.String()) == strings.ToLower(structField) {
			return mapParameter.MapIndex(v).Interface()
		}
	}
	return nil
}

func kind(value reflect.Value) reflect.Kind {
	kind := value.Kind()
	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float64
	default:
		return kind
	}
}

func setInt(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)

	switch parameter.(type) {
	case int:
		value.SetInt(int64(parV.Int()))
	case uint:
		value.SetInt(parV.Int())
	case float64:
		value.SetInt(int64(parV.Int()))
	default:
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}
	return nil
}

func setUint(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)

	switch parameter.(type) {
	case int:
		value.SetUint(uint64(parV.Int()))
	case uint:
		value.SetUint(parV.Uint())
	case float64:
		value.SetUint(uint64(parV.Float()))
	default:
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}
	return nil
}

func setBool(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)

	switch parameter.(type) {
	case bool:
		value.SetBool(parV.Bool())
	default:
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}
	return nil
}

func setString(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)
	switch parameter.(type) {
	case int, uint:
		value.SetString(fmt.Sprintf("%d", parV.Int()))
	case float64:
		value.SetString(fmt.Sprintf("%G", parV.Float()))
	case string:
		value.SetString(parV.String())
	default:
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}
	return nil
}

func setFloat(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)

	switch parameter.(type) {
	case int:
		value.SetFloat(float64(parV.Int()))
	case uint:
		value.SetFloat(float64(parV.Int()))
	case float64:
		value.SetFloat(float64(parV.Float()))
	default:
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}
	return nil
}

func setMap(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)
	switch parV.Type() {
	case value.Type():
		value.Set(parV)
	}
	return nil
}

func setSlice(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)

	switch value.Type() {
	case parV.Type():
		value.Set(parV)
	default:
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}

	return nil
}

func setStruct(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)

	switch parV.Kind() {
	case reflect.Struct:
		if value.Type() == parV.Type() {
			value.Set(parV)
		}
	case reflect.Map:
		// Encode(parameter, reflect.Indirect(reflect.New(value.Type())))
	default:
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}

	return nil
}

func setInterface(parameter interface{}, value reflect.Value) error {
	if parameter == nil {
		return nil
	}
	parV := reflect.ValueOf(parameter)
	if !parV.Type().AssignableTo(value.Type()) {
		return errors.New(fmt.Sprintf("unconvertible: %v - %v", value.Type(), parV.Type()))
	}
	value.Set(parV)
	return nil
}
