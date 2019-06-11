package parse

import (
	"io"
	"strings"
	"text/template"

	log "github.com/cihub/seelog"
)

const ServiceServerTemplate = `{{$pkg := .PkgPath|pkgForSort}}
{{$Package := .Package}}
package {{.Package}}
import (
	"{{.PkgPath}}"
	{{range $i, $v := .ImportPkgs}}{{genImport $i $v}}
{{end}}
	"github.com/gogo/protobuf/types"
	"runtime"
)
{{range $index, $iface := .Interfaces}}
func New{{.Name}}Server({{.Name}}0 {{$pkg}}.{{.Name}}) {{.Name}}Handler {
	return &{{.Name}}ServerImpl{
		{{.Name}}_:{{.Name}}0,
	}
}
type {{.Name}}ServerImpl struct {
	{{.Name}}_ {{$pkg}}.{{.Name}}
}
{{range $i, $v := .Funcs}}
{{$ParamsLen := .Params|len|reduce1}}
{{$ResultsLen := .Results|len|reduce1}}
func ({{$iface.Name}} *{{$iface.Name}}ServerImpl) {{.Name}} (ctx context.Context, request *{{.Name}}Req_, respone *{{.Name}}Resp_) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("rpc-server error: %v \n%s", e, debug.Stack()))
		}
	}()
	{{range $i1, $v1 := .Params}}
	{{end}}
{{if ne $ResultsLen -1}}{{range $i1, $v1 := .Results}} {{.Name}} {{if ne $i1 $ResultsLen }},{{end}} {{end}}  := {{end}}{{$iface.Name}}.{{$iface.Name}}_.{{.Name}}(ctx, {{range $i1, $v1 := .Params}}{{.Name}} {{if ne $i1 $ParamsLen }},{{end}}{{end}} )
	{{range $i1, $v1 := .Results}}
	{{end}}
	return 
} 
{{end}}
{{end}}
{{range $i, $v := .Structs}}
	{{if $v.IsRecursion }}
func to{{.Name}}MicroModelServer(v  model.{{.Name}}) *{{.Name}} {
	_resp := &{{.Name}}{} 
	{{range $i1, $v1 := .Fields}}
	{{end}}
	return _resp
}
{{end}}
{{end}}
`

func (file *File) GenServer(wr io.Writer) {
	log.Info("generating server file ...")
	t := template.New("pb.server.go")
	t.Funcs(template.FuncMap{
		"reduce1": func(i int) int {
			return i - 1
		},
		"pkgForSort": func(p string) string {
			index := strings.LastIndex(p, "/")
			pkgForSort := p
			if index != -1 {
				pkgForSort = p[index+1:]
			}
			return pkgForSort
		},
		"ServerImpl": func(n string) string {
			if strings.HasSuffix(n, "Service") {
				return n
			}
			return n + "Service"
		},
		"genImport": genImport,
	})

	parsed, err := t.Parse(ServiceServerTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	parsed.Execute(wr, file)
}
