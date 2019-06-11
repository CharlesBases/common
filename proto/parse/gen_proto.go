package parse

import (
	"html/template"
	"io"

	log "github.com/cihub/seelog"
)

const ProtoTemplate = `syntax = "proto3";

package {{.Package}};

import "google/protobuf/struct.proto";
{{range $index, $iface := .Interfaces}}
service {{.Name}} {
{{range $i0, $v := .Funcs}}	rpc {{.Name}} ({{.Name}}Req_) returns ({{.Name}}Resp_) {} 
{{end}}}
{{range $i, $v := .Funcs}} 
message {{.Name}}Req_ {
{{range $i1, $v1 := .Params}}	{{.ProtoType}} {{.Name}} = {{$i1 | add1}};
{{end}}}

message {{$v.Name}}Resp_ {
{{range $i1, $v1 := .Results}}	{{.ProtoType}} {{.Name}} = {{$i1 | add1}};
{{end}}}{{end}}{{end}}
{{range $k, $v := .Structs}}
message {{$v.Name}} {
{{range $i1, $v1 := .Fields}}	{{.ProtoType}} {{.Name}} = {{$i1 | add1}};
{{end}}}
{{end}}
`

func (file *File) GenProtoFile(wr io.Writer) {
	log.Info("generating .proto files ...")
	temp := template.New("pb.proto")
	temp.Funcs(template.FuncMap{
		"add1": func(i int) int {
			return i + 1
		},
	})
	protoTemplate, err := temp.Parse(ProtoTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	protoTemplate.Execute(wr, file)
}
