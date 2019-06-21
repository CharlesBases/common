package parse

import (
	"go/parser"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
	"golang.org/x/tools/go/loader"
)

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
		log.Info("parse structs in package: ", key)
		conf := loader.Config{ParserMode: parser.ParseComments}
		conf.Import(key)
		program, err := conf.Load()
		if err != nil {
			log.Error(err) // load error
			continue
		}
		astFiles := program.Package(key).Files
		Root := Package{root: root}
		Root.PkgPath = key
		Root.Files = make([]File, 0, len(astFiles))
		for _, astFile := range astFiles {
			structFile := Root.ParseStruct(value, astFile)
			if structFile == nil {
				continue
			}
			structFile.PkgPath = Root.PkgPath
			structFile.ParsePkgStruct(root)
			Root.Files = append(Root.Files, *structFile)
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
		for _, fieldValue := range structValue.Fields {
			goType := strings.Replace(strings.Replace(fieldValue.GoType, "[]", "", 1), "*", "", -1)
			if structValue.Name == goType {
				file.Structs[structKey].IsRecursion = true
			}
		}
	}
}

func (file *File) ParseStructs() {
	for structIndex, fileStruct := range file.Structs {
		for fieldIndex, field := range fileStruct.Fields {
			protoType := file.parseType(field.GoType)
			file.Structs[structIndex].Fields[fieldIndex].ProtoType = protoType
		}
	}
}

func (file *File) ParseStructMessage() {
	structMessage := make(map[string][]Message, 0)
	for key, val := range file.Message {
		imp := strings.TrimPrefix(val, "*")
		if index := strings.Index(imp, "."); index != -1 {
			impPrefix := imp[:index]
			imp, ok := file.ImportA[impPrefix]
			if ok {
				if structMessage[imp] == nil {
					structMessage[imp] = make([]Message, 0)
				}
				structMessage[imp] = append(structMessage[imp], Message{
					Name:     key,
					ExprName: val,
					FullName: imp,
				})
			}
		} else {
			_, ok := golangBaseType[val]
			if !ok {
				pkgpath := file.PkgPath
				if structMessage[pkgpath] == nil {
					structMessage[pkgpath] = make([]Message, 0)
				}
				message := Message{
					Name:     key,
					ExprName: val,
					FullName: pkgpath,
				}
				structMessage[pkgpath] = append(structMessage[pkgpath], message)
			}
		}
	}
	file.StructMessage = structMessage
}

func (file *File) parseType(golangType string) string {
	if file.Message == nil {
		file.Message = make(map[string]string, 0)
	}
	builder := strings.Builder{}
	if strings.HasPrefix(golangType, "[]") {
		if strings.HasPrefix(golangType, "[]byte") {
			return "bytes"
		} else {
			builder.WriteString("repeated ")
			golangType = strings.TrimPrefix(golangType, "[]")
		}
	}
	if protoType2RPCType, ok := golangBaseType2ProtoBaseType[golangType]; ok {
		builder.WriteString(protoType2RPCType)
	} else {
		if protoType, ok := golangType2ProtoType[golangType]; ok {
			builder.WriteString(protoType)
		} else {
			if strings.HasPrefix(golangType, "map") {
				if index := strings.Index(golangType, "]"); index != -1 {
					builder.WriteString("map<string, google.protobuf.Value>")
				}
			} else {
				protoType = strings.TrimPrefix(golangType, "*")
				if index := strings.LastIndex(protoType, "."); index != -1 {
					protoType = protoType[index+1:]
					builder.WriteString(protoType)
				} else {
					builder.WriteString(protoType)
				}
			}
			if golangType != "context.Context" {
				file.Message[protoType] = golangType
			}
		}
	}

	return builder.String()
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
