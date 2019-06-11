package parse

import (
	"io"
	"text/template"

	log "github.com/cihub/seelog"
)

const ServiceClientTemplate = `//this file is generated from {{.PkgPath}}
{{$pkg := .PkgPath|pkgForSort}}
{{$Package := .Package}}
package {{.Package}}
import (
	"github.com/gogo/protobuf/types"
	{{range $i, $v := .ImportPkgs}}{{genImport $i $v}} {{end}}
)
{{range $i, $v := .Funcs}}
{{$ParamsLen := .Params|len|reduce1}}
{{$ResultsLen := .Results|len|reduce1}}
func ({{$iface.Name}} *{{$iface.Name}}MicroClientImpl) {{.Name}}({{declareContextParam $v}}{{range $i1, $v1 := .Params}}{{.Name}} {{.GoType}} {{if ne $i1  $ParamsLen }},{{end}} {{end}}) ({{range $i1, $v1 := .Results}} {{.Name}} {{.GoType}}{{if ne $i1 $ResultsLen }},{{end}}{{end}}) {
	_req := new({{.Name}}Req_)
	{{range $i1, $v1 := .Params}}
	_req.{{.MicroName}}={{convert2Micro . ""}}
	{{end}}
	_resp, _e := {{$iface.Name}}.{{$iface.Name}}_.{{.Name}}(ctx, _req,{{$iface.Name}}.opts...)
	if _e != nil {
		log.Error(_e)
		panic(_e.Error())
	}
	{{range $i1, $v1 := .Results}}
	{{.Name}} = {{convert2Go . "_resp"}}
	{{end}}
    {{if eq $ResultsLen -1}}_=_resp{{end}}
	return {{range $i1, $v1 := .Results}} {{.Name}}{{if ne $i1 $ResultsLen }},{{end}} {{end}}
} 
{{end}}
{{end}}
{{range $i, $v := .Structs}}
	{{if $v.IsRecursion }}
func to{{.Name}}GoModelClient(v *{{.Name}}) model.{{.Name}} {
	v0 := model.{{.Name}}{} 
	{{range $i1, $v1 := .Fields}}
	v0.{{.Name}} = {{convert2Go . "v"}}
	{{end}}
	return v0
}
	{{end}}
{{end}}
`

func (file *File) GenClient(wr io.Writer) {
	log.Info("generating client file ...")
	t := template.New("pb.client.go")
	t.Funcs(template.FuncMap{
		"reduce1": func(i int) int {
			return i - 1
		},
		"pkgForSort": pkgForSort,
		"genImport":  genImport,
	})

	parsed, err := t.Parse(ServiceClientTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	parsed.Execute(wr, file)
}

/*func (file *File) convert(field Field) string {
	if field.LeftVal == "" {
		field.LeftVal = field.Name + "="
	}
	repeated := strings.Contains(field.ProtoType, "repeated")
	if repeated {
		field.ProtoType = strings.Replace(field.ProtoType, "repeated", "", -1)
		field.ProtoType = strings.TrimSpace(field.ProtoType)
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
			sb.WriteString(f.MicroName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(f.MicroName)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString(field.ProtoType)
			sb.WriteString("(v)\n")
			sb.WriteString("}\n")
			sb.WriteString(field.LeftVal + rightVal)
			return sb.String()
		} else {
			return field.GoType + "(" + f.MicroName + ")"
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
			sb.WriteString(f.MicroName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(f.MicroName)
			sb.WriteString("{\n")
			sb.WriteString(rightVal)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.DecodeValue(v)\n")
			sb.WriteString("}\n")
			sb.WriteString(field.LeftVal + rightVal)
			return sb.String()
		} else {
			return "proto3.DecodeValue(" + f.MicroName + ")"
		}
	case "google.protobuf.Struct":
		if f.GoType == "[]map[string]interface{}" {
			sb := strings.Builder{}
			sb.WriteString("make([]map[string]interface{}")
			sb.WriteString(",len(")
			sb.WriteString(f.MicroName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(f.MicroName)
			sb.WriteString("{\n")
			sb.WriteString(name)
			sb.WriteString("[i]=")
			sb.WriteString("proto3.DecodeToMap(v)\n")
			sb.WriteString("}\n")
			return sb.String()
		} else if field.GoType == "map[string]interface{}" {
			return "proto3.DecodeToMap(" + f.MicroName + ")"
		} else if field.GoType == "[]error" {
			sb := strings.Builder{}
			sb.WriteString("make([]weberror.BaseWebError,0,len(")
			sb.WriteString(f.MicroName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(f.MicroName)
			sb.WriteString("{\n")
			sb.WriteString("error__:=weberror.BaseWebError{} \n ")
			sb.WriteString("proto3.ConvertMapToStruct(proto3.DecodeToMap(")
			sb.WriteString(f.MicroName)
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
			sb.WriteString(f.MicroName)
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
			sb.WriteString(f.MicroName)
			sb.WriteString("))\n")
			sb.WriteString("for i, v := range ")
			sb.WriteString(f.MicroName)
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
				sb.WriteString("to" + GoType + "GoModelClient(v)")
			} else {
				ty := strings.TrimPrefix(field.GoType, "[]")
				r := ty + "{}\n"
				r = strings.Replace(r, "*", "&", 1)
				for _, v := range file.Structs {
					if v.Name == f.ProtoType {
						for _, v1 := range v.Fields {

							v1.GoExpr = rightVal
							v1.LeftVal = rightVal + "[i]." + v1.FieldName + "="
							r += "\n" + v1.LeftVal + file.convert(v1)

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
				sb.WriteString(" to" + GoType + "GoModelClient(v)")
			} else {
				r := ""
				if strings.Contains(field.GoType, "*") {
					sb.WriteString("nil\n ")
					r = name + "=" + field.GoType + "{}\n"
				} else {
					sb.WriteString(field.GoType + "{}\n")
				}
				sb.WriteString("if " + f.MicroName + "!=nil{\n")
				r = strings.Replace(r, "*", "&", 1)
				for _, v := range file.Structs {
					if v.Name == field.ProtoType {
						for _, v1 := range v.Fields {
							v1.GoExpr = name
							v1.LeftVal = name + "." + v1.FieldName + "="
							r += "\n" + v1.LeftVal + file.convert(v1)

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
*/
