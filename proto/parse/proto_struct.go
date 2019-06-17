package parse

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	log "github.com/cihub/seelog"
	proto "github.com/golang/protobuf/ptypes/struct"
)

func DecodeProtoStruct2Map(protoStruct *proto.Struct) map[string]interface{} {
	if protoStruct == nil {
		return nil
	}
	Map := map[string]interface{}{}
	for key, val := range protoStruct.Fields {
		Map[key] = DecodeProtoStruct2Interface(val)
	}
	return Map
}

func DecodeProtoStruct2Interface(protoStruct *proto.Value) interface{} {
	if protoStruct == nil {
		return nil
	}
	switch kind := protoStruct.Kind.(type) {
	case *proto.Value_NullValue:
		return nil
	case *proto.Value_NumberValue:
		return kind.NumberValue
	case *proto.Value_StringValue:
		return kind.StringValue
	case *proto.Value_BoolValue:
		return kind.BoolValue
	case *proto.Value_StructValue:
		return DecodeProtoStruct2Map(kind.StructValue)
	case *proto.Value_ListValue:
		Interface := make([]interface{}, len(kind.ListValue.Values))
		for key, val := range kind.ListValue.Values {
			Interface[key] = DecodeProtoStruct2Interface(val)
		}
		return Interface
	default:
		panic("protos_truct: unknown kind")
	}
}

// EncodeMap2ProtoStruct converts a map[string]interface{} to a ptypes.Struct
func EncodeMap2ProtoStruct(v map[string]interface{}) *proto.Struct {
	size := len(v)
	if size == 0 {
		return nil
	}
	fields := make(map[string]*proto.Value, size)
	for k, v := range v {
		fields[k] = EncodeInterface2ProtoStruct(v)
	}
	return &proto.Struct{
		Fields: fields,
	}
}

// EncodeInterface2ProtoStruct converts an interface{} to a ptypes.Value
func EncodeInterface2ProtoStruct(v interface{}) *proto.Value {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return &proto.Value{
			Kind: &proto.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &proto.Value{
			Kind: &proto.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		fields := make(map[string]*proto.Value, 2)
		fields["Code"] = &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: -1,
			},
		}
		fields["Err"] = &proto.Value{
			Kind: &proto.Value_StringValue{
				StringValue: v.Error(),
			},
		}
		return &proto.Value{
			Kind: &proto.Value_StructValue{
				StructValue: &proto.Struct{
					Fields: fields,
				},
			},
		}
	default:
		// Fallback to reflection for other types
		return toValue(reflect.ValueOf(v))
	}
}

// ConvertMapToStruct converts a map[string]interface{} to a struct
func ConvertMap2Struct(v map[string]interface{}, p interface{}) (ok bool) {
	if v == nil {
		return false
	}
	log.Debug(v)
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Error(err)
		return false
	}
	err = json.Unmarshal(bytes, p)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug(p)
	return true
}

// ConvertStructToMap converts a struct to a map[string]interface{}
func ConvertStruct2Map(obj interface{}) map[string]interface{} {
	if obj == nil {
		return nil
	}
	var data = make(map[string]interface{})
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Error(err)
		return nil
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Error(err)
		return nil
	}
	return data
}

func toValue(v reflect.Value) *proto.Value {
	switch v.Kind() {
	case reflect.Bool:
		return &proto.Value{
			Kind: &proto.Value_BoolValue{
				BoolValue: v.Bool(),
			},
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v.Int()),
			},
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: float64(v.Uint()),
			},
		}
	case reflect.Float32, reflect.Float64:
		return &proto.Value{
			Kind: &proto.Value_NumberValue{
				NumberValue: v.Float(),
			},
		}
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return toValue(reflect.Indirect(v))
	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return toValue(v.Elem())
	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			return nil
		}
		size := v.Len()
		values := make([]*proto.Value, size)
		for i := 0; i < size; i++ {
			values[i] = toValue(v.Index(i))
		}
		return &proto.Value{
			Kind: &proto.Value_ListValue{
				ListValue: &proto.ListValue{
					Values: values,
				},
			},
		}
	case reflect.Struct:
		t := v.Type()
		size := v.NumField()
		if size == 0 {
			return nil
		}
		fields := make(map[string]*proto.Value, size)
		for i := 0; i < size; i++ {
			field := t.Field(i)
			// 支持内嵌结构体展开
			if field.Anonymous && (field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Ptr) {
				if _, ok := field.Tag.Lookup("protoOpen"); ok {
					// 展开结构体
					sizeF := field.Type.NumField()
					typeF := field.Type

					for j := 0; j < sizeF; j++ {
						tagName := findTagNameJsonOrDefault(typeF.Field(j))
						fields[tagName] = toValue(v.FieldByIndex([]int{i, j}))
					}
					continue
				}
			}

			name := field.Name
			tagName := findTagNameJsonOrDefault(field)
			// Better way?
			if len(name) > 0 && 'A' <= name[0] && name[0] <= 'Z' {
				fields[tagName] = toValue(v.Field(i))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &proto.Value{
			Kind: &proto.Value_StructValue{
				StructValue: &proto.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.Map: // 只支持键为string的map
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		fields := make(map[string]*proto.Value, len(keys))
		for _, k := range keys {
			if k.Kind() == reflect.String {
				fields[k.String()] = toValue(v.MapIndex(k))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &proto.Value{
			Kind: &proto.Value_StructValue{
				StructValue: &proto.Struct{
					Fields: fields,
				},
			},
		}
	default:
		// Last resort
		return &proto.Value{
			Kind: &proto.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}

func findTagNameJsonOrDefault(f reflect.StructField) string {
	defaultName := f.Name

	tagName := f.Tag.Get("json")
	if len(tagName) == 0 {
		return defaultName
	}

	tagName = strings.Split(tagName, ",")[0]
	if tagName == "-" {
		tagName = defaultName
	}

	return tagName
}
