package parse

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	log "github.com/cihub/seelog"
)

const ServiceClientTemplate = `// this file is generated from {{.PkgPath}} {{$pkg := .PkgPath | packageSort}} {{$Package := .Package}}
package {{.Package}}
import (
	{{range $importIndex, $import := .ImportA}}{{generateImport $importIndex $import}}
	{{end}}

	// _struct "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
)
{{range $interfaceIndex, $interface := .Interfaces}}
func New{{.Name}}Client_({{.Name}} {{.Name}}Client, opts ...grpc.CallOption) {{$pkg}}.{{.Name}} {
	return &{{.Name}}ClientImpl{
		{{.Name}}:  {{.Name}},
		opts:       opts,
	}
}

type {{.Name}}ClientImpl struct {
	{{.Name}}   {{.Name}}Client
	opts        []grpc.CallOption
}
{{range $funcsIndex, $func := .Funcs}} {{$ParamsLen := .Params | len | funcReduce}} {{$ResultsLen := .Results | len | funcReduce}}
func ({{$interface.Name}} *{{$interface.Name}}ClientImpl) {{.Name}}(ctx context.Context, {{range $paramsIndex, $param := .Params}}{{.Name}} {{.GoType}} {{if ne $paramsIndex  $ParamsLen }},{{end}} {{end}}) ({{range $resultsIndex, $result := .Results}} {{.Name}} {{.GoType}}{{if ne $resultsIndex $ResultsLen }},{{end}} {{end}}) {
	serviceRequest := &{{.Name}}Req_{
		{{range $paramsIndex, $param := .Params}}
	}
	serviceRequest.{{.Variable}}={{convertClientRequest . ""}}
	{{end}}
	serviceResponse, serviceError := {{$interface.Name}}.{{$interface.Name}}_.{{.Name}}(ctx, serviceRequest, {{$interface.Name}}.opts...)
	if serviceError != nil {
		log.Error(serviceError)
		panic(serviceError.Error())
	}
	{{range $resultsIndex, $result := .Results}}
	{{.Name}} = {{convertClientResponse . "serviceResponse"}}
	{{end}}
    {{if eq $ResultsLen -1}}_ = serviceResponse{{end}}
	return {{range $resultsIndex, $result := .Results}} {{.Name}}{{if ne $resultsIndex $ResultsLen }},{{end}} {{end}}
} 
{{end}}
{{end}}
{{range $structsIndex, $struct := .Structs}}
{{if $struct.IsRecursion }}
func clientModel{{.Name}}(value  *{{.Name}}) model.{{.Name}} {
	result := model.{{.Name}}{} 
	{{range $fieldsIndex, $field := .Fields}}
	result.{{.Name}} = {{convertClientRequest . "value"}}
	{{end}}
	return result
}
{{end}}
{{end}}
`

func (file *File) GenClient(wr io.Writer) {
	log.Info("generating client file ...")
	t := template.New("pb.client.go")
	t.Funcs(template.FuncMap{
		"funcReduce": func(i int) int {
			return i - 1
		},
		"service": func(n string) string {
			if strings.HasSuffix(n, "Service") {
				return n
			}
			return n + "Service"
		},
		"packageSort":           packageSort,
		"generateImport":        generateImport,
		"convertClientRequest":  file.convertClientRequest,
		"convertClientResponse": file.convertClientResponse,
	})

	parsed, err := t.Parse(ServiceClientTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	parsed.Execute(wr, file)
}

func (file *File) convertClientRequest(field Field, expr string) string {
	if field.VariableCall == "" {
		field.VariableCall = fmt.Sprintf("%s.%s", "serviceRequest", field.Variable)
	}

	theType := field.ProtoType

	isRepeated := strings.Contains(field.ProtoType, "repeated")
	if isRepeated {
		field.ProtoType = strings.TrimPrefix(field.ProtoType, "repeated ")
	}

	fieldName := field.Name
	if expr != "" {
		fieldName = expr + "." + field.FieldName
	}

	theType1, ok := golangBaseType2ProtoBaseType[theType]
	if ok {
		theType = theType1
		if isRepeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(fieldName, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make([]")
			sb.WriteString(theType)
			sb.WriteString(",len(")
			sb.WriteString(fieldName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(fieldName)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString(theType)
			sb.WriteString("(v)\n")
			sb.WriteString("}\n")

			sb.WriteString(field.VariableCall + rightVal)

			return sb.String()
		} else {
			return theType + "(" + fieldName + ")"
		}
	}
	switch field.ProtoType {
	case "google.protobuf.Value":
		if isRepeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(fieldName, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make([]*structpb.Value")
			sb.WriteString(",len(")
			sb.WriteString(fieldName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(fieldName)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.EncodeToValue(v)\n")
			sb.WriteString("}\n")

			sb.WriteString(field.VariableCall + rightVal)

			return sb.String()
		} else {
			return "proto3.EncodeToValue(" + fieldName + ")"
		}
	case "google.protobuf.Struct":
		if field.GoType == "[]map[string]interface{}" {
			sb := strings.Builder{}
			sb.WriteString("make([]*structpb.Struct")
			sb.WriteString(",len(")
			sb.WriteString(fieldName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(fieldName)
			sb.WriteString("{\n")
			sb.WriteString(field.Variable)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.EncodeMapToStruct(v)\n")
			sb.WriteString("}\n")
			return sb.String()
		} else if field.GoType == "map[string]interface{}" {
			return "proto3.EncodeMapToStruct(" + fieldName + ")"
		} else if field.GoType == "[]error" {
			sb := strings.Builder{}
			sb.WriteString("make([]*structpb.Struct")
			sb.WriteString(",len(")
			sb.WriteString(fieldName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(fieldName)
			sb.WriteString("{\n")
			sb.WriteString(field.Variable)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.EncodeMapToStruct(proto3.ConvertStructToMap(weberror.BaseWebError{Code:-1,Err:errors.New(v.Error())}))\n")
			sb.WriteString("}\n")
			return sb.String()
		} else if field.GoType == "error" {
			sb := strings.Builder{}
			sb.WriteString("nil\n if ")
			sb.WriteString(fieldName)
			sb.WriteString("!=nil{\n")
			sb.WriteString("if _weberror, _ok := ")
			sb.WriteString(fieldName)
			sb.WriteString(".(weberror.BaseWebError); _ok {\n")
			sb.WriteString("_resp.")
			sb.WriteString(field.FieldName)
			sb.WriteString("=proto3.EncodeMapToStruct(proto3.ConvertStructToMap(_weberror))\n} else {\n")
			sb.WriteString("_resp.")
			sb.WriteString(field.FieldName)
			sb.WriteString("=proto3.EncodeMapToStruct(proto3.ConvertStructToMap(weberror.BaseWebError{Code:-1,Err:errors.New(")
			sb.WriteString(fieldName)
			sb.WriteString(".Error())}))\n}\n}\n")
			return sb.String()
		}
	default:
		if isRepeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(fieldName, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make([]*")
			sb.WriteString(theType)
			sb.WriteString(",len(")
			sb.WriteString(fieldName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(fieldName)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")

			if field.IsRecursion {
				if field.VariableCall == field.Name+"=" {
					field.VariableCall = "v0." + field.Name + "="
				}
				GoType := strings.Replace(field.GoType, "[", "", -1)
				GoType = strings.Replace(GoType, "]", "", -1)
				GoType = strings.Replace(GoType, "*", "", -1)
				index := strings.Index(GoType, ".")
				if index != -1 {
					GoType = GoType[index+1:]
				}
				sb.WriteString("to" + GoType + "MicroModelClient(v)")

			} else {
				r := "new(" + theType + ")\n"
				for _, v := range file.Structs {
					if v.Name == theType {
						for _, v1 := range v.Fields {
							v1.VariableCall = rightVal + "[i]." + v1.Variable + "="
							r += "\n" + v1.VariableCall + file.convertClientRequest(v1, "v")
						}
					}
				}
				sb.WriteString(r)

			}
			sb.WriteString("}\n")
			sb.WriteString(field.VariableCall + rightVal)

			return sb.String()
		} else {

			sb := strings.Builder{}

			if field.IsRecursion {
				GoType := strings.Replace(field.GoType, "[", "", -1)
				GoType = strings.Replace(GoType, "]", "", -1)
				GoType = strings.Replace(GoType, "*", "", -1)
				index := strings.Index(GoType, ".")
				if index != -1 {
					GoType = GoType[index+1:]
				}
				sb.WriteString("to" + GoType + "MicroModelClient(v)")
			} else {
				r := ""
				if strings.Contains(field.ProtoType, "map<") {
					r = fieldName + "\n"
				} else {
					r = "new(" + theType + ")\n"
				}

				for _, v := range file.Structs {
					if v.Name == theType {
						for range v.Fields {

						}
					}
				}
				sb.WriteString(r)
			}
			return sb.String()
		}
	}
	return field.Name
}

func (file *File) convertClientResponse(field Field, expr string) string {
	if field.VariableCall == "" {
		field.VariableCall = field.Name
	}

	repeated := strings.Contains(field.ProtoType, "repeated")
	if repeated {
		field.ProtoType = strings.Replace(field.ProtoType, "repeated", "", -1)
		field.ProtoType = strings.TrimSpace(field.ProtoType)
	}

	if expr != "" {
		field.Variable = expr + "." + field.Variable
	}

	name := field.Name
	if field.GoExpr != "" {
		name = field.GoExpr + "." + field.FieldName
	}

	_, ok := protoBaseType[field.ProtoType]
	if ok {
		if repeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(name, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make(")
			sb.WriteString(field.GoType)
			sb.WriteString(",len(")
			sb.WriteString(field.Variable)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(field.Variable)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString(field.ProtoType)
			sb.WriteString("(v)\n")
			sb.WriteString("}\n")

			sb.WriteString(field.VariableCall + rightVal)

			return sb.String()
		} else {
			return field.GoType + "(" + field.Variable + ")"
		}
	}
	switch field.ProtoType {
	case "google.protobuf.Value":
		if repeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(name, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make([]interface{}")
			sb.WriteString(",len(")
			sb.WriteString(field.Variable)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(field.Variable)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.DecodeValue(v)\n")
			sb.WriteString("}\n")

			sb.WriteString(field.VariableCall + rightVal)

			return sb.String()
		} else {
			return "proto3.DecodeValue(" + field.Variable + ")"
		}
	case "google.protobuf.Struct":
		if field.GoType == "[]map[string]interface{}" {
			sb := strings.Builder{}
			sb.WriteString("make([]map[string]interface{}")
			sb.WriteString(",len(")
			sb.WriteString(field.Variable)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(field.Variable)
			sb.WriteString("{\n")

			sb.WriteString(name)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.DecodeToMap(v)\n")
			sb.WriteString("}\n")
			return sb.String()
		} else if field.GoType == "map[string]interface{}" {
			return "proto3.DecodeToMap(" + field.Variable + ")"
		} else if field.GoType == "[]error" {
			sb := strings.Builder{}
			sb.WriteString("make([]weberror.BaseWebError,0,len(")
			sb.WriteString(field.Variable)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(field.Variable)
			sb.WriteString("{\n")
			sb.WriteString("error__:=weberror.BaseWebError{} \n ")
			sb.WriteString("proto3.ConvertMapToStruct(proto3.DecodeToMap(")
			sb.WriteString(field.Variable)
			sb.WriteString("),&error__)\n")
			sb.WriteString("if error__.Code != 0 && error__.Err != nil {\n")
			sb.WriteString(name)
			sb.WriteString("=append(")
			sb.WriteString(name)
			sb.WriteString(" ,error__)\n")
			sb.WriteString("}\n}\n")
			return sb.String()
		} else if field.GoType == "error" {
			sb := strings.Builder{}
			sb.WriteString("nil\n ")
			sb.WriteString("error__:=weberror.BaseWebError{} \n ")
			sb.WriteString("proto3.ConvertMapToStruct(proto3.DecodeToMap(")
			sb.WriteString(field.Variable)
			sb.WriteString("),&error__)\n")
			sb.WriteString("if error__.Code != 0 && error__.Err != nil {\n")
			sb.WriteString(field.Name)
			sb.WriteString(" =error__\n")
			sb.WriteString("}\n")
			return sb.String()
		}
	default:
		if repeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(name, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make(")
			sb.WriteString(field.GoType)
			sb.WriteString(",len(")
			sb.WriteString(field.Variable)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(field.Variable)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")

			if field.IsRecursion {
				if field.VariableCall == field.Name+"=" {
					field.VariableCall = "v0." + field.Name + "="
				}

				GoType := strings.Replace(field.GoType, "[", "", -1)
				GoType = strings.Replace(GoType, "]", "", -1)
				GoType = strings.Replace(GoType, "*", "", -1)
				index := strings.Index(GoType, ".")
				if index != -1 {
					GoType = GoType[index+1:]
				}
				sb.WriteString("to" + GoType + "GoModelClient(v)")
			} else {
				ty := strings.TrimPrefix(field.GoType, "[]")
				r := ty + "{}\n"
				r = strings.Replace(r, "*", "&", 1)
				for _, v := range file.Structs {
					if v.Name == field.ProtoType {
						for _, v1 := range v.Fields {

							v1.GoExpr = rightVal
							v1.VariableCall = rightVal + "[i]." + v1.FieldName + "="
							r += "\n" + v1.VariableCall + file.convertClientResponse(v1, "v")

						}
					}
				}
				sb.WriteString(r)
			}

			sb.WriteString("}\n")

			sb.WriteString(field.VariableCall + rightVal)

			return sb.String()
		} else {
			sb := strings.Builder{}
			if field.IsRecursion {
				GoType := strings.Replace(field.GoType, "[", "", -1)
				GoType = strings.Replace(GoType, "]", "", -1)
				GoType = strings.Replace(GoType, "*", "", -1)
				index := strings.Index(GoType, ".")
				if index != -1 {
					GoType = GoType[index+1:]
				}
				sb.WriteString(" to" + GoType + "GoModelClient(v)")
			} else {

				r := ""
				if strings.Contains(field.GoType, "*") {
					sb.WriteString("nil\n ")
					r = name + "=" + field.GoType + "{}\n"
				} else {
					sb.WriteString(field.GoType + "{}\n")
				}
				sb.WriteString("if " + field.Variable + "!=nil{\n")
				r = strings.Replace(r, "*", "&", 1)
				//r := ""
				for _, v := range file.Structs {
					if v.Name == field.ProtoType {
						for _, v1 := range v.Fields {
							v1.GoExpr = name
							v1.VariableCall = name + "." + v1.FieldName + "="
							r += "\n" + v1.VariableCall + file.convertClientResponse(v1, field.Variable)

						}
					}
				}
				sb.WriteString(r)
				sb.WriteString("}\n")
			}

			return sb.String()

		}
	}
	return field.Name
}
