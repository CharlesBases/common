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
	"runtime/debug"

	"{{.PkgPath}}"
	{{range $index, $importA := .ImportA}}{{generateImport $index $importA}}
	{{end}}

	proto "github.com/CharlesBases/common/proto/parse"

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
{{range $funcsIndex, $func := .Funcs}} {{$ParamsLen := .Params|len|funcReduce}} {{$ResultsLen := .Results|len|funcReduce}}
func ({{$interface.Name}} *{{$interface.Name}}ServerImpl) {{.Name}} (ctx context.Context, serviceRequest *{{.Name}}Req_, serviceResponse *{{.Name}}Resp_) (err_ error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("rpc-server error: %v \n%s", err, debug.Stack()))
		}
	}()
	{{range $paramsIndex, $param := .Params}}
		{{.Name}} := {{convertServerRequest $param "serviceRequest"}}
	{{end}}

	{{if ne $ResultsLen -1}}
		{{range $resultsIndex, $result := .Results}}{{.Name}}{{if ne $resultsIndex $ResultsLen }},{{end}} {{end}} := {{end}}{{$interface.Name}}.{{$interface.Name}}.{{.Name}}(ctx, {{range $paramsIndex, $param := .Params}}{{.Name}}{{if ne $paramsIndex $ParamsLen }},{{end}}{{end}})
	{{range $resultsIndex, $result := .Results}}
	serviceResponse.{{.Variable}} = {{convertServerResponse . ""}}
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
	result.{{.Name}} = {{convertServerResponse . "value"}}
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
		"service": func(n string) string {
			if strings.HasSuffix(n, "Service") {
				return n
			}
			return n + "Service"
		},
		"generateImport":        generateImport,
		"convertServerRequest":  file.convertServerRequest,
		"convertServerResponse": file.convertServerResponse,
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
			return fmt.Sprintf(`func() %s {
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
			)
		} else {
			return fmt.Sprintf("%s(%s)", field.GoType, field.VariableCall)
		}
	}
	switch field.ProtoType {
	case "google.protobuf.Value":
		if field.GoType == "[]error" {
			return fmt.Sprintf(`func() []error {
						errors := make([]error, len(%s))
						for key, val := range %s {
							errors[key] = fmt.Errorf("%s", val)
						}
						return errors
					}()`,
				field.VariableCall,
				field.VariableCall,
				"%v",
			)
		}
		if field.GoType == "error" {
			return fmt.Sprintf(`func() error {
						if %s != nil {
							return fmt.Errorf("%s", %s)
						} else {
							return nil
						}
					}()`,
				field.VariableCall,
				"%v",
				field.VariableCall,
			)
		}
		if field.GoType == "[]interface{}" {
			return fmt.Sprintf(`func() []interface{} {
					list := make([]interface{}, len(%s))
					for key, val := range %s {
						list[key] = proto.DecodeProtoValue2Interface(val)
					}
					return list
				}()`,
				field.VariableCall,
				field.VariableCall,
			)
		}
		if field.GoType == "interface{}" {
			return "proto.DecodeProtoValue2Interface(" + field.VariableCall + ")"
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

func (file *File) convertServerResponse(field Field, expr string) string {
	if field.VariableCall == "" {
		field.VariableCall = "serviceResponse." + field.Variable
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
			return fmt.Sprintf(`func() []%s {
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
			)
		} else {
			return fmt.Sprintf("%s(%s)", protoType, fieldName)
		}
	}
	switch field.ProtoType {
	case "google.protobuf.Value":
		if isRepeated {
			return fmt.Sprintf(`func() []*_struct.Value {
						errors := make([]*_struct.Value, len(%s))
						for key, val := range %s {
							errors[key] = proto.EncodeInterface2ProtoValue(val)
						}
						return errors
					}()`,
				fieldName,
				fieldName,
			)
		} else {
			return fmt.Sprintf("proto.DecodeProtoValue2Interface(%s)", fieldName)
		}
	default:
		if isRepeated {
			return fmt.Sprintf(`func() []*%s {
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
						return fmt.Sprintf("serverModel$s(val)", GoType)
					} else {
						str := strings.Builder{}
						for _, fileStruct := range file.Structs {
							if fileStruct.Name == field.ProtoType {
								for _, structField := range fileStruct.Fields {
									str.WriteString(fmt.Sprintf("%s: %s,\n", structField.FieldName, file.convertServerResponse(structField, "val")))
								}
							}
						}
						return str.String()
					}
				}(),
			)
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
							return fmt.Sprintf("serverModel%s(%s)", GoType, fieldName)
						} else {
							str := strings.Builder{}
							str.WriteString(fmt.Sprintf("&%s{\n", field.ProtoType))
							for _, fileStruct := range file.Structs {
								if fileStruct.Name == field.ProtoType {
									for _, structField := range fileStruct.Fields {
										str.WriteString(fmt.Sprintf("%s: %s,\n", structField.FieldName, file.convertServerResponse(structField, fieldName)))
									}
								}
							}
							str.WriteString("}\n")
							return str.String()
						}
					}(),
				))

				return fmt.Sprintf(`func() *%s {
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
							return fmt.Sprintf("serverModel%s(%s)", GoType, fieldName)
						} else {
							str := strings.Builder{}
							str.WriteString(fmt.Sprintf("&%s{\n", field.ProtoType))
							for _, fileStruct := range file.Structs {
								if fileStruct.Name == field.ProtoType {
									for _, structField := range fileStruct.Fields {
										str.WriteString(fmt.Sprintf("%s: %s,\n", structField.FieldName, file.convertServerResponse(structField, fieldName)))
									}
								}
							}
							str.WriteString("}\n")
							return str.String()
						}
					}(),
				)
			} /* else {
				if field.IsRecursion {
					GoType := strings.Replace(field.GoType, "[", "", -1)
					GoType = strings.Replace(GoType, "]", "", -1)
					GoType = strings.Replace(GoType, "*", "", -1)
					index := strings.Index(GoType, ".")
					if index != -1 {
						GoType = GoType[index+1:]
					}
					return fmt.Sprintf("serverModel%s(val)", GoType)
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
								r += "\n" + structField.VariableCall + file.convertServerResponse(structField, fieldName)
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
