package parse

import (
	"go/ast"

	log "github.com/cihub/seelog"
)

var (
	golangBaseType2ProtoBaseType = map[string]string{
		"bool":    "bool",
		"string":  "string",
		"int":     "sint64",
		"int32":   "sint64",
		"int64":   "sint64",
		"uint":    "uint64",
		"uint32":  "uint64",
		"uint64":  "uint64",
		"float32": "double",
		"float64": "double",
	}
	golangType2ProtoType = map[string]string{
		"error":       "google.protobuf.Value",
		"interface{}": "google.protobuf.Value",
	}
	golangBaseType = map[string]struct{}{
		"byte":    {},
		"bool":    {},
		"string":  {},
		"int":     {},
		"int32":   {},
		"int64":   {},
		"uint":    {},
		"uint32":  {},
		"uint64":  {},
		"float32": {},
		"float64": {},

		"error":       {},
		"interface{}": {},
	}
	protoBaseType = map[string]string{
		"bool":   "bool",
		"bytes":  "[]byte",
		"string": "string",
		"sint64": "int64",
		"uint64": "uint64",
		"double": "float64",
	}
)

type Package struct {
	Name         string
	Path         string
	PkgPath      string
	Files        []File
	MessageTypes map[string][]string
	root         *Package
}

type Message struct {
	Name     string //struct名字
	ExprName string //调用名 （pager.PagerListResp）
	FullName string // 全名 （带包名）
}

type Interface struct {
	Funcs  []Func
	Name   string
	IsFunc bool
}

type Struct struct {
	Name        string
	Fields      []Field
	Pkg         string // go类型定义的所在包
	IsRecursion bool   // 递归应用类型
}

type InterfaceImpl struct {
	Methods []Method
	Name    string
}

// Method represents a method signature.
type Method struct {
	Recv string
	Func
}

// Func represents a function signature.
type Func struct {
	Name    string
	Params  []Field
	Results []Field
}

// Field represents a parameter in a function or method signature.
type Field struct {
	Name         string // 字段名 原参数名或返回值名或struct中的字段名
	FieldName    string // 原参数名或返回值名的可导出形式
	GoType       string // 正常go类型
	ProtoType    string // proto类型
	GoExpr       string // go类型的引用前缀
	Package      string // go类型定义的所在包
	IsField      bool   // struct中的字段
	Variable     string // 被赋值变量
	VariableType string // 变量类型
	VariableCall string // 变量调用名
	IsRecursion  bool   // 递归应用类型
}

func (root *Package) ParseStruct(message []Message, astFile *ast.File) *File {
	file := File{}
	file.PkgPath = root.PkgPath

	file.ParseImport(astFile)

	structs := make([]Struct, 0, 1)
	ast.Inspect(astFile, func(x ast.Node) bool {
		switch x.(type) {
		case *ast.TypeSpec:
			spec := x.(*ast.TypeSpec)
			structType, ok := spec.Type.(*ast.StructType)
			if !ok {
				return true
			}
			var (
				isContainsA bool
				isContainsB bool
			)
			if message == nil {
				isContainsA = true
			} else {
				for _, v := range message {
					if v.Name == spec.Name.Name {
						isContainsA = true
					}
				}
			}
			if root.root.MessageTypes == nil {
				root.root.MessageTypes = make(map[string][]string, 0)
				isContainsB = false
			} else {
				messageType, ok := root.root.MessageTypes[root.PkgPath]
				if ok {
					for _, v := range messageType {
						if v == spec.Name.Name {
							isContainsB = true
						}
					}
				} else {
					root.root.MessageTypes[root.PkgPath] = make([]string, 0)
				}
			}
			if isContainsA && !isContainsB {
				s := file.ParseStruct(spec.Name.Name, structType)
				log.Info("find struct: ", spec.Name.Name)
				structs = append(structs, s)
				root.root.MessageTypes[root.PkgPath] = append(root.root.MessageTypes[root.PkgPath], spec.Name.Name)
			}
		default:
			return true
		}
		return false
	})
	file.Structs = structs
	return &file
}

//把gofiles 汇总到 一个gofile
func (root *Package) Summary() File {
	log.Info("SummaryGofiles...")
	if len(root.Files) == 0 {
		return File{}
	}
	var file = root.Files[0]

	if len(root.Files) > 1 {
		for key, val := range root.Files[1:] {
			if key == 0 {
				continue
			}
			file.Structs = append(file.Structs, val.Structs...)
			file.Interfaces = append(file.Interfaces, val.Interfaces...)
			// var i int
			// for imp, quote := range val.Imports {
			// 	_, ok1 := file.PkgImports[quote]
			// 	if !ok1 { // 未导入包
			// 		_, ok2 := file.ImportPkgs[imp]
			// 		if ok2 { //包名会冲突
			// 			kk := imp + strconv.Itoa(i)
			// 			file.ImportPkgs[kk] = quote
			// 			file.PkgImports[quote] = kk
			// 			i++
			// 		} else { // 包名不会冲突
			// 			file.ImportPkgs[imp] = quote
			// 			file.PkgImports[quote] = imp
			// 		}
			// 	}
			// }
		}
	}
	return file
}
