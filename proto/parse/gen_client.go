package parse

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	log "github.com/cihub/seelog"
)

const ServiceClientTemplate = `// this file is generated from {{.PkgPath}} {{$pkg := .PkgPath|pkgForSort}} {{$Package := .Package}}
package {{.Package}}
import (
	{{range $importIndex, $import := .ImportA}}{{funcImport $importIndex $import}}
	{{end}}

	// _struct "github.com/golang/protobuf/ptypes/struct"
)
{{range $interfaceIndex, $interface := .Interfaces}}
func New{{.Name}}Client({{.Name}}_ {{.Name | service}}) {{$pkg}}.{{.Name}} {
	return &{{.Name}}ClientImpl{
		{{.Name}}: {{.Name}}_,
	}
}
type {{.Name}}ClientImpl struct {
	{{.Name}} {{.Name | microName}}
}
{{range $funcsIndex, $func := .Funcs}} {{$ParamsLen := .Params|len|funcReduce}} {{$ResultsLen := .Results|len|funcReduce}}
func ({{$interface.Name}} *{{$interface.Name}}ClientImpl) {{.Name}}(ctx context.Context{{range $paramsIndex, $param := .Params}}{{.Name}} {{.GoType}} {{if ne $paramsIndex  $ParamsLen }},{{end}} {{end}}) ({{range $resultsIndex, $result := .Results}} {{.Name}} {{.GoType}}{{if ne $resultsIndex $ResultsLen }},{{end}} {{end}}) {
	serviceRequest := &{{.Name}}Req_{
		{{range $paramsIndex, $param := .Params}}
	}
	serviceRequest.{{.Variable}}={{convertClientRequest . ""}}
	{{end}}
	serviceResponse, serviceError := {{$interface.Name}}.{{$interface.Name}}_.{{.Name}}({{useContextParam $v}} serviceRequest,{{$interface.Name}}.opts...)
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
		theType = strings.Replace(theType, "repeated", "", -1)
		theType = strings.TrimSpace(theType)
	}

	name := field.Name
	if expr != "" {
		name = expr + "." + field.FieldName
	}

	theType1, ok := protoBaseType[theType]
	if ok {
		theType = theType1
		if isRepeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(name, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make([]")
			sb.WriteString(theType)
			sb.WriteString(",len(")
			sb.WriteString(name)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(name)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString(theType)
			sb.WriteString("(v)\n")
			sb.WriteString("}\n")

			sb.WriteString(field.LeftVal + rightVal)

			return sb.String()
		} else {
			return theType + "(" + name + ")"
		}
	}
	switch field.ProtoType {
	case "google.protobuf.Value":
		if isRepeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(name, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make([]*structpb.Value")
			sb.WriteString(",len(")
			sb.WriteString(name)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(name)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.EncodeToValue(v)\n")
			sb.WriteString("}\n")

			sb.WriteString(field.LeftVal + rightVal)

			return sb.String()
		} else {
			return "proto3.EncodeToValue(" + name + ")"
		}
	case "google.protobuf.Struct":
		if field.GoType == "[]map[string]interface{}" {
			sb := strings.Builder{}
			sb.WriteString("make([]*structpb.Struct")
			sb.WriteString(",len(")
			sb.WriteString(name)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(name)
			sb.WriteString("{\n")
			sb.WriteString(field.MicroExpr)
			sb.WriteString(field.MicroName)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.EncodeMapToStruct(v)\n")
			sb.WriteString("}\n")
			return sb.String()
		} else if field.GoType == "map[string]interface{}" {
			return "proto3.EncodeMapToStruct(" + name + ")"
		} else if field.GoType == "[]error" {
			sb := strings.Builder{}
			sb.WriteString("make([]*structpb.Struct")
			sb.WriteString(",len(")
			sb.WriteString(name)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(name)
			sb.WriteString("{\n")
			sb.WriteString(field.MicroExpr)
			sb.WriteString(field.MicroName)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.EncodeMapToStruct(proto3.ConvertStructToMap(weberror.BaseWebError{Code:-1,Err:errors.New(v.Error())}))\n")
			sb.WriteString("}\n")
			return sb.String()
		} else if field.GoType == "error" {
			sb := strings.Builder{}
			sb.WriteString("nil\n if ")
			sb.WriteString(name)
			sb.WriteString("!=nil{\n")
			sb.WriteString("if _weberror, _ok := ")
			sb.WriteString(name)
			sb.WriteString(".(weberror.BaseWebError); _ok {\n")
			sb.WriteString("_resp.")
			sb.WriteString(field.FieldName)
			sb.WriteString("=proto3.EncodeMapToStruct(proto3.ConvertStructToMap(_weberror))\n} else {\n")
			sb.WriteString("_resp.")
			sb.WriteString(field.FieldName)
			sb.WriteString("=proto3.EncodeMapToStruct(proto3.ConvertStructToMap(weberror.BaseWebError{Code:-1,Err:errors.New(")
			sb.WriteString(name)
			sb.WriteString(".Error())}))\n}\n}\n")
			return sb.String()
		}
	default:
		if isRepeated {
			sb := strings.Builder{}

			rightVal := strings.Replace(name, ".", "", -1)
			rightVal = strings.Replace(rightVal, "[", "", -1)
			rightVal = strings.Replace(rightVal, "]", "", -1)
			rightVal = "_slice_" + rightVal
			sb.WriteString("nil\n" + rightVal + ":=")

			sb.WriteString("make([]*")
			sb.WriteString(theType)
			sb.WriteString(",len(")
			sb.WriteString(name)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(name)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")

			if field.IsRecursion {
				if field.LeftVal == field.Name+"=" {
					field.LeftVal = "v0." + field.Name + "="
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
				for _, v := range gf.Structs {
					if v.Name == theType {
						for _, v1 := range v.Fields {
							v1.MicroExpr = rightVal + "."
							v1.LeftVal = rightVal + "[i]." + v1.MicroName + "="
							r += "\n" + v1.LeftVal + gf.convert2MicroClient(v1, "v")
						}
					}
				}
				sb.WriteString(r)

			}
			sb.WriteString("}\n")
			sb.WriteString(field.LeftVal + rightVal)

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
					r = name + "\n"
				} else {
					r = "new(" + theType + ")\n"
				}

				for _, v := range gf.Structs {
					if v.Name == theType {
						for _, v1 := range v.Fields {

							v1.MicroExpr = field.MicroExpr + field.MicroName + "."
							v1.LeftVal = field.MicroExpr + field.MicroName + "." + v1.MicroName + "="
							r += "\n" + v1.LeftVal + gf.convert2MicroClient(v1, name)

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

func (file *File) convertClientResponse() {

}
