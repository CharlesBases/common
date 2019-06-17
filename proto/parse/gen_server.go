package parse

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/cihub/seelog"
)

const ServiceServerTemplate = `// this file is generated from {{.PkgPath}} {{$Package := .PkgPath|funcSort}}
package {{.Package}}

import (
	"fmt"
	"runtime"

	proto "github.com/CharlesBases/common/proto/parse"
	"{{.PkgPath}}"
	{{range $index, $importA := .ImportA}}{{funcImport $index $importA}}
	{{end}}
	log "github.com/cihub/seelog"
	_struct "github.com/golang/protobuf/ptypes/struct"
)
{{range $interfaceIndex, $interface := .Interfaces}}
func New{{.Name}}Server({{.Name}} {{$Package}}.{{.Name}}) {{.Name}}Server {
	return &{{.Name}}ServerImpl{
		{{.Name}}:{{.Name}},
	}
}

type {{.Name}}ServerImpl struct {
	{{.Name}} {{$Package}}.{{.Name}}
}
{{range $funcsIndex, $func := .Funcs}}{{$ParamsLen := .Params|len|funcReduce}}{{$ResultsLen := .Results|len|funcReduce}}
func ({{$interface.Name}} *{{$interface.Name}}ServerImpl) {{.Name}} (ctx context.Context, request_ *{{.Name}}Req_, respone_ *{{.Name}}Resp_) (err_ error) {
	defer func() {
		if e := recover(); e != nil {
			funcName := ""
			if pc, _, _, ok := runtime.Caller(1); ok {
				funcName = runtime.FuncForPC(pc).Name()
			}
			log.Error(fmt.Sprintf("rpc-server error: %s(%v) \n%s", funcName, e, debug.Stack()))
		}
	}()
	{{range $paramsIndex, $param := .Params}}
		{{.Name}} := {{convertServerRequest $param "request_"}}
	{{end}}

	{{if ne $ResultsLen -1}}
		{{range $resultsIndex, $result := .Results}}{{.Name}}{{if ne $resultsIndex $ResultsLen }},{{end}} {{end}} := {{end}}{{$interface.Name}}.{{$interface.Name}}.{{.Name}}(ctx, {{range $paramsIndex, $param := .Params}}{{.Name}}{{if ne $paramsIndex $ParamsLen }},{{end}}{{end}})
	{{range $resultsIndex, $result := .Results}}
	respone_.{{.Variable}} = {{convertServerRespone . ""}}
	{{end}}
	return 
} 
{{end}}
{{end}}
{{range $structsIndex, $struct := .Structs}}
	{{if $struct.IsRecursion }}
func serverModel{{.Name}}(value  model.{{.Name}}) *{{.Name}} {
	result := &{{.Name}}{} 
	{{range $fieldsIndex, $field := .Fields}}
	result.{{.Name}} = {{convertServerRespone . "value"}}
	{{end}}
	return result
}
{{end}}
{{end}}
`

func (file *File) GenServer(wr io.Writer) {
	log.Info("generating server file ...")
	t := template.New("pb.server.go")
	t.Funcs(template.FuncMap{
		"funcReduce": func(i int) int {
			return i - 1
		},
		"funcSort": func(Package string) string {
			return filepath.Base(Package)
		},
		"funcInterface": func(n string) string {
			if strings.HasSuffix(n, "Service") {
				return n
			}
			return n + "Service"
		},
		"funcImport":           genImport,
		"convertServerRequest": file.convertServerRequest,
		"convertServerRespone": file.convertServerRespone,
	})

	parsed, err := t.Parse(ServiceServerTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	parsed.Execute(wr, file)
}

func (file *File) convertServerRequest(field Field, expr string) string {
	if field.Variable == "" {
		field.Variable = field.Name
	}
	isRepeated := strings.Contains(field.ProtoType, "repeated")
	if isRepeated {
		field.ProtoType = strings.Replace(field.ProtoType, "repeated", "", -1)
		field.ProtoType = strings.TrimSpace(field.ProtoType)
	}
	fieldName := field.Name
	if field.GoExpr != "" {
		fieldName = fmt.Sprintf("%s.%s", field.GoExpr, field.FieldName)
	}
	if expr != "" {
		field.VariableCall = fmt.Sprintf("%s.%s", expr, field.Variable)
	}
	if _, ok := protoBaseType[field.ProtoType]; ok {
		if isRepeated {
			builder := strings.Builder{}

			builder.WriteString(fmt.Sprintf(`func() %s {
					list := make(%s, len(%s))
					for key, val := range %s {
						list[key] = val
					}
					return list
				}()`,
				field.GoType,
				field.GoType,
				field.VariableCall,
				field.VariableCall,
			))

			return builder.String()
		} else {
			return field.GoType + "(" + field.VariableCall + ")"
		}
		return fmt.Sprintf("%s(%s)", field.GoType, field.VariableCall)
	}
	switch field.ProtoType {
	case "google.protobuf.Value":
		if isRepeated {
			builder := strings.Builder{}

			builder.WriteString(fmt.Sprintf(`func() []interface{} {
					list := make([]interface{}, len(%s))
					for key, val := range %s {
						list[key] = proto.DecodeProtoStruct2Interface(val)
					}
					return list
				}()`,
				field.VariableCall,
				field.VariableCall,
			))

			return builder.String()
		} else {
			return "proto.DecodeProtoStruct2Interface(" + field.VariableCall + ")"
		}
	case "google.protobuf.Struct":
		if /*field.GoType == "error" {
			return "nil\n error__:=weberror.BaseWebError{} \n  proto.ConvertMapToStruct(proto.DecodeProtoStruct2Map(" + field.VariableCall + "),&error__)\n" +
				field.Name + " =error__"
		} else if field.GoType == "[]error" {
			builder := strings.Builder{}

			builder.WriteString("make([]weberror.BaseWebError")
			builder.WriteString(",len(")
			builder.WriteString(field.VariableCall)
			builder.WriteString("))\n")
			builder.WriteString("for i, fileStruct := range ")
			builder.WriteString(field.VariableCall)
			builder.WriteString("{\n")
			builder.WriteString(fieldName)
			builder.WriteString("error__:=weberror.BaseWebError{} \n  proto.ConvertMapToStruct(proto.DecodeProtoStruct2Map(fileStruct),&error__)\n")
			builder.WriteString("[i]=error__\n")
			builder.WriteString("}\n")

			return builder.String()
		} else if*/field.GoType == "map[string]interface{}" {
			return "proto.DecodeProtoStruct2Map(" + field.VariableCall + ")"
		} else if field.GoType == "[]map[string]interface{}" {
			builder := strings.Builder{}

			builder.WriteString(fmt.Sprintf(`func() interface{} {
					list := make(interface{}, len(%s))
					for key, val := range %s {
						list[key] = proto.DecodeProtoStruct2Map(val)
					}
					return list
				}()`,
				field.VariableCall,
				field.VariableCall,
			))

			return builder.String()
		}
	default:
		if isRepeated {
			builder := strings.Builder{}

			variable := strings.Replace(fieldName, ".", "", -1)
			variable = strings.Replace(variable, "[", "", -1)
			variable = strings.Replace(variable, "]", "", -1)
			variable = "slice" + strings.Title(variable)

			builder.WriteString(fmt.Sprintf("make(%s, 0)\n", field.GoType))
			builder.WriteString(fmt.Sprintf("%s := make(%s,len(%s))\n", field.GoType, field.VariableCall))
			builder.WriteString(fmt.Sprintf("for key, val := range %s {\n", field.VariableCall))
			builder.WriteString(fmt.Sprintf("%s[key] = %s\n}\n", variable,
				func() string {
					str := strings.Replace(fmt.Sprintf("%s{}\n", strings.TrimPrefix(field.GoType, "[]")), "*", "&", 1)
					for _, fileStruct := range file.Structs {
						if fileStruct.Name == field.ProtoType {
							for _, structField := range fileStruct.Fields {
								if structField.IsRecursion {
									GoType := strings.Replace(field.GoType, "[", "", -1)
									GoType = strings.Replace(GoType, "]", "", -1)
									GoType = strings.Replace(GoType, "*", "", -1)
									index := strings.Index(GoType, ".")
									if index != -1 {
										GoType = GoType[index+1:]
									}
									structField.Variable = variable + "[i]." + structField.FieldName + "="
									str += "\n" + structField.Variable + "to" + GoType + "GoModelServer(fileStruct)"
								} else {
									structField.GoExpr = fieldName
									structField.Variable = variable + "[i]." + structField.FieldName + "="
									str += "\n" + structField.Variable + file.convertServerRequest(structField, "fileStruct")
								}
							}
						}
					}
					return str
				}()))
			builder.WriteString(fmt.Sprintf("%s = %s", field.Variable, variable))

			return builder.String()
		} else {
			variableValue := strings.Builder{}
			variableValue.WriteString(strings.Replace(fmt.Sprintf("%s{\n", field.GoType), "*", "&", 1))
			for _, fileStruct := range file.Structs {
				if fileStruct.Name == field.ProtoType {
					for _, structField := range fileStruct.Fields {
						if structField.IsRecursion {
							GoType := strings.Replace(field.GoType, "[", "", -1)
							GoType = strings.Replace(GoType, "]", "", -1)
							GoType = strings.Replace(GoType, "*", "", -1)
							index := strings.Index(GoType, ".")
							if index != -1 {
								GoType = GoType[index+1:]
							}
							variableValue.WriteString(fmt.Sprintf("%s: to%sGoModelServer(val),\n", structField.FieldName, GoType))
						} else {
							variableValue.WriteString(fmt.Sprintf("%s: %s,\n", structField.FieldName, file.convertServerRequest(structField, field.VariableCall)))
						}
					}
				}
			}
			variableValue.WriteString("}\n")
			return variableValue.String()
		}
	}
	return field.Name
}

func (file *File) convertServerRespone(field Field, expr string) string {
	if field.VariableCall == "" {
		field.VariableCall = "respone_." + field.Variable
	}
	isRepeated := strings.Contains(field.ProtoType, "repeated")
	if isRepeated {
		field.ProtoType = strings.Replace(field.ProtoType, "repeated", "", -1)
		field.ProtoType = strings.TrimSpace(field.ProtoType)
	}
	fieldName := field.Name
	if expr != "" {
		fieldName = expr + "." + field.FieldName
	}

	protoType, ok := protoBaseType[field.ProtoType]
	if ok {
		if isRepeated {
			build := strings.Builder{}

			build.WriteString(fmt.Sprintf(`func() []%s {
					list := make([]%s, len(%s))
					for key, val := range %s {
						list[key] = val
					}
					return list
				}()`,
				protoType,
				protoType,
				fieldName,
				fieldName,
			))

			return build.String()
		} else {
			return protoType + "(" + fieldName + ")"
		}
	}
	switch field.ProtoType {
	case "google.protobuf.Value":
		if isRepeated {
			build := strings.Builder{}

			build.WriteString(fmt.Sprintf(`func() []*_struct.Value {
					list := make([]*_struct.Value, len(%s))
					for key, val := range %s {
						list[key] = proto.EncodeInterface2ProtoValue(val)
					}
					return list
				}()`,
				fieldName,
				fieldName,
			))

			return build.String()
		} else {
			return "proto.EncodeInterface2ProtoValue(" + fieldName + ")"
		}
	case "google.protobuf.Struct":
		if field.GoType == "[]map[string]interface{}" {
			build := strings.Builder{}

			build.WriteString(fmt.Sprintf(`func() []*_struct.Struct {
					list := make([]*_struct.Struct, len(%s))
					for key, val := range %s {
						list[key] = proto.EncodeMap2ProtoStruct(val)
					}
					return list
				}()`,
				fieldName,
				fieldName,
			))

			return build.String()
		} else if field.GoType == "[]error" {
			build := strings.Builder{}

			build.WriteString(fmt.Sprintf(`func() []*_struct.Struct {
					if %s != nil {
						list := make([]*_struct.Struct, len(%s))
						for key, val := range %s {
							list[key] = proto.EncodeMap2ProtoStruct(map[string]interface{}{"err": val})
						}
						return list
					}
					return nil
				}()`,
				fieldName,
				fieldName,
				fieldName,
			))

			return build.String()
		} else if field.GoType == "error" {
			build := strings.Builder{}

			build.WriteString(fmt.Sprintf(`func() *_struct.Struct {
					if %s != nil {
						return proto.EncodeMap2ProtoStruct(map[string]interface{}{"err": %s})
					}
					return nil
				}()`,
				fieldName,
				fieldName))

			return build.String()
		} else if field.GoType == "map[string]interface{}" {
			return "proto.EncodeMap2ProtoStruct(" + fieldName + ")"
		}
	default:
		if isRepeated {
			build := strings.Builder{}

			build.WriteString(fmt.Sprintf(`func() []*%s {
					list := make([]*%s, len(%s))
					for key, val := range %s {
						list[key] = &%s{
							%s
						}
					}
					return list
				}()`,
				field.ProtoType,
				field.ProtoType,
				fieldName,
				fieldName,
				field.ProtoType,
				func() string {
					if field.IsRecursion {
						GoType := strings.Replace(field.GoType, "[", "", -1)
						GoType = strings.Replace(GoType, "]", "", -1)
						GoType = strings.Replace(GoType, "*", "", -1)
						index := strings.Index(GoType, ".")
						if index != -1 {
							GoType = GoType[index+1:]
						}
						return fmt.Sprintf("to%sMicroModelServer(val)", GoType)
					} else {
						str := strings.Builder{}
						for _, fileStruct := range file.Structs {
							if fileStruct.Name == field.ProtoType {
								for _, structField := range fileStruct.Fields {
									str.WriteString(fmt.Sprintf("%s: %s,\n", structField.FieldName, file.convertServerRespone(structField, "val")))
								}
							}
						}
						return str.String()
					}
				}(),
			))

			return build.String()
		} else {
			build := strings.Builder{}
			if strings.Contains(field.GoType, "*") {
				build.WriteString(fmt.Sprintf(`func() *%s {
						if %s != nil {
							return %s
						}
						return nil
					}()`,
					field.ProtoType,
					fieldName,
					func() string {
						if field.IsRecursion {
							GoType := strings.Replace(field.GoType, "[", "", -1)
							GoType = strings.Replace(GoType, "]", "", -1)
							GoType = strings.Replace(GoType, "*", "", -1)
							index := strings.Index(GoType, ".")
							if index != -1 {
								GoType = GoType[index+1:]
							}
							return fmt.Sprintf("to%sMicroModelServer(%s)", GoType, fieldName)
						} else {
							str := strings.Builder{}
							str.WriteString(fmt.Sprintf("&%s{\n", field.ProtoType))
							for _, fileStruct := range file.Structs {
								if fileStruct.Name == field.ProtoType {
									for _, structField := range fileStruct.Fields {
										str.WriteString(fmt.Sprintf("%s: %s,\n", structField.FieldName, file.convertServerRespone(structField, fieldName)))
									}
								}
							}
							str.WriteString("}\n")
							return str.String()
						}
					}(),
				))

				return build.String()
			} /* else {
				if field.IsRecursion {
					GoType := strings.Replace(field.GoType, "[", "", -1)
					GoType = strings.Replace(GoType, "]", "", -1)
					GoType = strings.Replace(GoType, "*", "", -1)
					index := strings.Index(GoType, ".")
					if index != -1 {
						GoType = GoType[index+1:]
					}
					return fmt.Sprintf("to%sMicroModelServer(val)", GoType)
				} else {
					str := strings.Builder{}
					if strings.Contains(field.ProtoType, "map<") {
						str.WriteString(fieldName + "\n")
					} else {
						str.WriteString(fmt.Sprintf("new(%s)\n", field.ProtoType))
					}

					for _, fileStruct := range file.Structs {
						if fileStruct.Name == field.ProtoType {
							for _, structField := range fileStruct.Fields {
								structField.MicroExpr = field.MicroExpr + field.Variable + "."
								structField.VariableCall = field.MicroExpr + field.Variable + "." + structField.Variable + "="
								r += "\n" + structField.VariableCall + file.convertServerRespone(structField, fieldName)
							}
						}
					}
					return r
				}
			}*/
		}
	}
	return field.Name
}
