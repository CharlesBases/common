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
	ImportA       map[string]string
	ImportB       map[string]string
	Message       map[string]string
	StructMessage map[string][]Message
}

func NewFile(pkgname string, pkgpath string) File {
	return File{
		Package: pkgname,
		PkgPath: pkgpath,
	}
}

func (file *File) ParsePkgStruct(root *Package) {
	file.ParseStructs()
	file.ParseStructMessage()
	structMessage := file.StructMessage
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
		Root := Package{root: root}
		Root.PkgPath = key
		Root.Files = make([]File, 0, len(files))
		for _, file := range files {
			file := Root.ParseStruct(value, file)
			if file == nil {
				continue
			}
			file.PkgPath = Root.PkgPath
			file.ParsePkgStruct(root)
			Root.Files = append(Root.Files, *file)
		}
		if len(Root.Files) == 0 {
			continue
		}
		packages = append(packages, Root)
	}
	if len(packages) == 0 {
		return
	}
	if len(file.Structs) == 0 {
		file.Structs = make([]Struct, 0)
	}

	// 合并结果
	for _, packageValue := range packages {
		for _, fileValue := range packageValue.Files {
			file.Structs = append(file.Structs, fileValue.Structs...)
			var i int
			for key, val := range fileValue.ImportA {
				_, okB := file.ImportB[val]
				if !okB { // 未导入包
					_, okA := file.ImportA[val]
					if okA { //包名会冲突
						keyIndex := key + strconv.Itoa(i)
						file.ImportA[keyIndex] = val
						file.ImportB[val] = keyIndex
						i++
					} else { // 包名不会冲突
						file.ImportA[key] = val
						file.ImportB[val] = key
					}
				}
			}
		}
	}

	for structKey, structValue := range file.Structs {
		for fieldKey, fieldValue := range structValue.Fields {
			goType := strings.Replace(fieldValue.GoType, "[", "", -1)
			goType = strings.Replace(goType, "]", "", -1)
			goType = strings.Replace(goType, "*", "", -1)
			if structValue.Name == goType {
				file.Structs[structKey].IsRecursion = true
				file.Structs[structKey].Fields[fieldKey].IsRecursion = true
			}
		}
	}
}

func (file *File) ParseStructs() {
	for k, v := range file.Structs {
		for k1, v := range v.Fields {
			value := file.parseType(v.GoType)
			file.Structs[k].Fields[k1].ProtoType = value
		}
	}
}

func (file *File) ParseStructMessage() {
	message := make(map[string][]Message)
	for key, val := range file.Message {
		imp := strings.TrimPrefix(val, "*")
		index := strings.Index(imp, ".")
		if index == -1 {
			_, ok := goBaseType[val]
			if !ok {
				pkgpath := file.PkgPath
				if message[pkgpath] == nil {
					message[pkgpath] = make([]Message, 0)
				}
				mg := Message{
					Name:     key,
					ExprName: val,
					FullName: pkgpath,
				}
				message[pkgpath] = append(message[pkgpath], mg)
			}
		} else {
			impPrefix := imp[:index]
			imp, ok := file.ImportA[impPrefix]
			if ok {
				if message[imp] == nil {
					message[imp] = make([]Message, 0)
				}
				message[imp] = append(message[imp], Message{
					Name:     key,
					ExprName: val,
					FullName: imp,
				})
			}
		}
	}
	file.StructMessage = message
}

func (file *File) parseType(vType string) (value string) {
	if file.Message == nil {
		file.Message = map[string]string{}
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
			return "map<" + file.parseType(key1) + ", " + file.parseType(value1) + "> "
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
			file.Message[value] = vType
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

func (file *File) typeConfig(field *Field) string {
	if ImportA, ok := file.ImportA[field.Package]; ok {
		i := 0
		if strings.Contains(field.GoType, "[]") {
			i = i + 2
		}
		if strings.Contains(field.GoType, "*") {
			i = i + 1
		}
		index := strings.Index(field.GoType, ".")
		if i == 0 {
			if index == -1 {
				return ImportA + "." + field.GoType
			}
			return ImportA + "." + field.GoType[index+1:]
		}
		if index == -1 {
			return field.GoType[:i] + ImportA + "." + field.GoType[i:]
		}
		return field.GoType[:i] + ImportA + "." + field.GoType[index+1:]
	}
	return field.GoType
}
