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
	"runtime"

	proto "github.com/CharlesBases/common/proto/parse"
	"{{.PkgPath}}"
	{{range $i, $v := .ImportA}}{{funcImport $i $v}}
	{{end}}
	log "github.com/cihub/seelog"
	"github.com/gogo/protobuf/types"
)
{{range $index, $iface := .Interfaces}}
func New{{.Name}}Server({{.Name}} {{$Package}}.{{.Name}}) {{.Name}}Server {
	return &{{.Name}}ServerImpl{
		{{.Name}}:{{.Name}},
	}
}

type {{.Name}}ServerImpl struct {
	{{.Name}} {{$Package}}.{{.Name}}
}
{{range $i, $v := .Funcs}}{{$ParamsLen := .Params|len|funcReduce}}{{$ResultsLen := .Results|len|funcReduce}}
func ({{$iface.Name}} *{{$iface.Name}}ServerImpl) {{.Name}} (ctx context.Context, request_ *{{.Name}}Req_, respone *{{.Name}}Resp_) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("rpc-server error: %v \n%s", e, debug.Stack()))
		}
	}()
	{{range $i1, $v1 := .Params}}
		{{.Name}} := {{convertGo $v1 "request_"}}
	{{end}}

	{{if ne $ResultsLen -1}}
		{{range $i1, $v1 := .Results}}{{.Name}}{{if ne $i1 $ResultsLen }},{{end}} {{end}} := {{end}}{{$iface.Name}}.{{$iface.Name}}.{{.Name}}(ctx, {{range $i1, $v1 := .Params}}{{.Name}}{{if ne $i1 $ParamsLen }},{{end}}{{end}})
	{{range $i1, $v1 := .Results}}
	{{end}}
	return 
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
		"convertGo":  file.serverConvertGo,
		"funcImport": genImport,
	})

	parsed, err := t.Parse(ServiceServerTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	parsed.Execute(wr, file)
}

func (file *File) serverConvertGo(field Field, expr string) string {
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
		field.VariableCall = fmt.Sprintf("%s.%s", expr, field.FieldName)
	}
	if _, ok := protoBaseType[field.ProtoType]; ok {
		if isRepeated {
			builder := strings.Builder{}

			variable := strings.Replace(fieldName, ".", "", -1)
			variable = strings.Replace(variable, "[", "", -1)
			variable = strings.Replace(variable, "]", "", -1)
			variable = "slice" + strings.Title(variable)

			builder.WriteString(fmt.Sprintf("make(%s, 0)\n", field.GoType))
			builder.WriteString(fmt.Sprintf("%s := make(%s, len(%s))\n", variable, field.GoType, field.VariableCall))
			builder.WriteString(fmt.Sprintf("for key, val := range %s {\n", field.VariableCall))
			builder.WriteString(fmt.Sprintf("%s[key] = %s(val)\n}\n", variable, strings.TrimPrefix(field.GoType, "[]")))
			builder.WriteString(fmt.Sprintf("%s = %s", field.Variable, variable))

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

			variable := strings.Replace(fieldName, ".", "", -1)
			variable = strings.Replace(variable, "[", "", -1)
			variable = strings.Replace(variable, "]", "", -1)
			variable = "slice" + strings.Title(variable)

			builder.WriteString("make([]interface{},0)\n")
			builder.WriteString(fmt.Sprintf("%s := make([]interface{}, len(%s))\n", variable, field.VariableCall))
			builder.WriteString(fmt.Sprintf("for key, val := range %s {\n", field.VariableCall))
			builder.WriteString(fmt.Sprintf("%s[key] = proto.DecodeProtoStruct2Interface(val)\n}\n", variable))
			builder.WriteString(fmt.Sprintf("%s = %s", field.Variable, variable))

			return builder.String()
		} else {
			return "proto.DecodeProtoStruct2Interface(" + field.VariableCall + ")"
		}
	case "google.protobuf.Struct":
		if /*field.GoType == "error" {
			return "nil\n error__:=weberror.BaseWebError{} \n  proto3.ConvertMapToStruct(proto.DecodeProtoStruct2Map(" + field.VariableCall + "),&error__)\n" +
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
			builder.WriteString("error__:=weberror.BaseWebError{} \n  proto3.ConvertMapToStruct(proto.DecodeProtoStruct2Map(fileStruct),&error__)\n")
			builder.WriteString("[i]=error__\n")
			builder.WriteString("}\n")

			return builder.String()
		} else if*/field.GoType == "map[string]interface{}" {
			return "proto.DecodeProtoStruct2Map(" + field.VariableCall + ")"
		} else if field.GoType == "[]map[string]interface{}" {
			builder := strings.Builder{}

			builder.WriteString(fmt.Sprintf("make([]map[string]interface{}, len(%s))\n", field.VariableCall))
			builder.WriteString(fmt.Sprintf("for key, val := range %s {\n", field.VariableCall))
			builder.WriteString(fmt.Sprintf("%s[key] = proto.DecodeProtoStruct2Map(val)\n}\n", fieldName))

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
									str += "\n" + structField.Variable + file.serverConvertGo(structField, "fileStruct")
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
			variableValue.WriteString(strings.Replace(fmt.Sprintf("%s{}\n", field.GoType), "*", "&", 1))
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
							variableValue.WriteString(fmt.Sprintf("%s.%s = to%sGoModelServer(val)\n", fieldName, structField.FieldName, GoType))
						} else {
							variableValue.WriteString(fmt.Sprintf("%s.%s = %s\n", fieldName, structField.FieldName, file.serverConvertGo(structField, field.VariableCall)))
						}
					}
				}
			}
			return variableValue.String()
		}
	}
	return field.Name
}
