package parse

import (
	"go/parser"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
	"golang.org/x/tools/go/loader"
)

var protoBaseType = map[string]string{
	"float":  "float32",
	"double": "float64",
	"bytes":  "[]byte",
	"int32":  "int32",
	"int64":  "int64",
	"int8":   "int8",
	"int16":  "int16",
	"uint":   "uint",
	"uint32": "uint32",
	"uint64": "uint64",
	"uint8":  "uint8",
	"uint16": "uint16",
	"byte":   "byte",
	"bool":   "bool",
	"string": "string",
	"rune":   "rune",
}

var goBaseType = map[string]struct{}{
	"float32":                {},
	"float64":                {},
	"int":                    {},
	"int32":                  {},
	"int64":                  {},
	"int8":                   {},
	"int16":                  {},
	"uint":                   {},
	"uint32":                 {},
	"uint64":                 {},
	"uint8":                  {},
	"uint16":                 {},
	"byte":                   {},
	"bool":                   {},
	"string":                 {},
	"rune":                   {},
	"interface{}":            {},
	"map[string]interface{}": {},
	"error":                  {},
}

type File struct {
	Name          string
	Package       string
	PkgPath       string
	Structs       []Struct
	Interfaces    []Interface
	Imports       map[string]string
	Message       map[string]string
	StructMessage map[string][]Message
}

func NewFile(pkgname string, pkgpath string) File {
	return File{
		Package: pkgname,
		PkgPath: pkgpath,
	}
}

func (infor *File) ParsePkgStruct(root *Package) {
	infor.ParseStructs()
	infor.ParseStructMessage()
	structMessage := infor.StructMessage
	packages := make([]Package, 0)
	for key, value := range structMessage {
		log.Info("parse structs in package:", key)
		conf := loader.Config{ParserMode: parser.ParseComments}
		conf.Import(key)
		program, err := conf.Load()
		if err != nil {
			log.Error(err) // load error
			continue
		}
		files := program.Package(key).Files
		pkg := Package{root: root}
		pkg.PkgPath = key
		pkg.Files = make([]File, 0, len(files))
		for _, file := range files {
			file := pkg.ParseStruct(value, file)
			if file == nil {
				continue
			}
			file.PkgPath = pkg.PkgPath
			file.ParsePkgStruct(root)
			pkg.Files = append(pkg.Files, *file)
		}
		if len(pkg.Files) == 0 {
			continue
		}
		packages = append(packages, pkg)
	}
	if len(packages) == 0 {
		return
	}
	if len(infor.Structs) == 0 {
		infor.Structs = make([]Struct, 0)
	}
	// 合并结果

	for index := range packages {
		for _, file := range packages[index].Files {
			infor.Structs = append(infor.Structs, file.Structs...)
			var i int
			for key, value := range file.ImportPkgs {
				_, ok := infor.PkgImports[value]
				if !ok {
					_, ok_ := infor.ImportPkgs[key]
					if ok_ {
						key_ := key + strconv.Itoa(i)
						infor.ImportPkgs[key_] = value
						infor.PkgImports[value] = key_
						i++
					} else {
						infor.ImportPkgs[key] = value
						infor.PkgImports[value] = key
					}
				}
			}
		}
	}

	for structKey, structValue := range infor.Structs {
		for fieldKey, fieldValue := range structValue.Fields {
			goType := strings.Replace(fieldValue.GoType, "[", "", -1)
			goType = strings.Replace(goType, "]", "", -1)
			goType = strings.Replace(goType, "*", "", -1)
			if structValue.Name == goType {
				infor.Structs[structKey].IsRecursion = true
				infor.Structs[structKey].Fields[fieldKey].IsRecursion = true
			}
		}
	}
}

func (infor *File) ParseStructs() {
	for k, v := range infor.Structs {
		for k1, v := range v.Fields {
			value := infor.parseType(v.GoType)
			infor.Structs[k].Fields[k1].ProtoType = value
		}
	}
}

func (infor *File) ParseStructMessage() {
	message := make(map[string][]Message)
	for k, v := range infor.Message {
		imp := strings.TrimPrefix(v, "*")
		index := strings.Index(imp, ".")
		if index == -1 {
			_, ok := goBaseType[v]
			if !ok {
				s := infor.PkgPath
				if message[s] == nil {
					message[s] = make([]Message, 0)
				}
				m := Message{
					Name:     k,
					ExprName: v,
					FullName: s,
				}
				message[s] = append(message[s], m)
			}
		} else {
			impPrefix := imp[:index]
			s, ok := infor.ImportPkgs[impPrefix]
			if ok {
				if message[s] == nil {
					message[s] = make([]Message, 0)
				}
				message[s] = append(message[s], Message{
					Name:     k,
					ExprName: v,
					FullName: s,
				})
			}
		}
	}
	infor.StructMessage = message
}

func (infor *File) parseType(vType string) (value string) {
	if infor.Message == nil {
		infor.Message = map[string]string{}
	}
	var prefix string
	if !strings.HasPrefix(vType, "[]byte") {
		if strings.HasPrefix(vType, "[]") {
			prefix = "repeated "
			vType = strings.TrimPrefix(vType, "[]")
		}
	}

	if prefix == "" { // 不是数组
		if strings.HasPrefix(vType, "map") && vType != "map[string]interface{}" {
			vType1 := vType[4:]
			index := strings.LastIndex(vType1, "]")

			key1 := vType1[:index]
			value1 := vType1[index+1:]
			return "map<" + infor.parseType(key1) + ", " + infor.parseType(value1) + "> "
		}
	}

	value, ok := protoMap[vType]
	if !ok {
		value = strings.TrimPrefix(vType, "*")
		lastIndex := strings.LastIndex(value, ".")
		if lastIndex != -1 {
			value = value[lastIndex+1:]
		}
		if value != "Context" && vType != "context.Context" {
			infor.Message[value] = vType //非基本类型
		}
	}
	return prefix + value
}

func (file *File) GoTypeConfig() {
	for interfaceKey, interfaceValue := range file.Interfaces {
		for funcKey, funcValue := range interfaceValue.Funcs {
			for resultKey, resultVallue := range funcValue.Results {
				goType := file.typeConfig(&resultVallue)
				file.Interfaces[interfaceKey].Funcs[funcKey].Results[resultKey].GoType = goType
			}
			for paramKey, paramValue := range funcValue.Params {
				goType := file.typeConfig(&paramValue)
				file.Interfaces[interfaceKey].Funcs[funcKey].Params[paramKey].GoType = goType
			}
		}
	}
	for structKey, structValue := range file.Structs {
		for fieldKey, fieldValue := range structValue.Fields {
			goType := file.typeConfig(&fieldValue)
			file.Structs[structKey].Fields[fieldKey].GoType = goType
		}
	}
}

func (file *File) typeConfig(v *Field) string {
	if imports, ok := file.PkgImports[v.Package]; ok {
		i := 0
		if strings.Contains(v.GoType, "[]") {
			i = i + 2
		}
		if strings.Contains(v.GoType, "*") {
			i = i + 1
		}
		index := strings.Index(v.GoType, ".")
		if i == 0 {
			if index == -1 {
				return imports + "." + v.GoType
			}
			return imports + "." + v.GoType[index+1:]
		}
		if index == -1 {
			return v.GoType[:i] + imports + "." + v.GoType[i:]
		}
		return v.GoType[:i] + imports + "." + v.GoType[index+1:]
	}
	return v.GoType
}
