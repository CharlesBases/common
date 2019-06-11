package parse

import (
	"go/ast"
	"strconv"

	log "github.com/cihub/seelog"
)

var (
	// golang type 对应的proto3 type
	protoMap = map[string]string{
		"error":                  "google.protobuf.Struct",
		"interface{}":            "google.protobuf.Value",
		"map[string]interface{}": "google.protobuf.Struct",
		"float64":                "double",
		"float32":                "float",
		"int":                    "int64",
		"int8":                   "int32",
		"int16":                  "int32",
		"uint":                   "uint64",
		"uint8":                  "uint32",
		"uint16":                 "uint32",
		"[]byte":                 "bytes",
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
	Name        string // 字段名 原参数名或返回值名或struct中的字段名
	FieldName   string // 原参数名或返回值名的可导出形式
	GoType      string // 正常go类型
	ProtoType   string // proto类型
	GoExpr      string // go类型的引用前缀
	Package     string // go类型定义的所在包
	IsField     bool   // struct中的字段
	LeftVal     string // 被赋值变量
	IsRecursion bool   // 递归应用类型
}

func (pkg *Package) ParseStruct(message []Message, astFile *ast.File) *File {
	file := File{}
	file.PkgPath = pkg.PkgPath

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
			if pkg.root.MessageTypes == nil {
				pkg.root.MessageTypes = make(map[string][]string, 0)
				isContainsB = false
			} else {
				messageType, ok := pkg.root.MessageTypes[pkg.PkgPath]
				if ok {
					for _, v := range messageType {
						if v == spec.Name.Name {
							isContainsB = true
						}
					}
				} else {
					pkg.root.MessageTypes[pkg.PkgPath] = make([]string, 0)
				}
			}
			if isContainsA && !isContainsB {
				s := file.ParseStruct(spec.Name.Name, structType)
				log.Info("find struct: ", spec.Name.Name)
				structs = append(structs, s)
				pkg.root.MessageTypes[pkg.PkgPath] = append(pkg.root.MessageTypes[pkg.PkgPath], spec.Name.Name)
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
func (p *Package) Summary() File {
	log.Info("SummaryGofiles...")
	if len(p.Files) == 0 {
		return File{}
	}
	var gf = p.Files[0]

	if len(p.Files) > 1 {
		for k, v := range p.Files[1:] {
			if k == 0 {
				continue
			}
			gf.Structs = append(gf.Structs, v.Structs...)
			gf.Interfaces = append(gf.Interfaces, v.Interfaces...)
			//fmt.Print(v.PkgImports, v.ImportPkgs)
			var i int
			for k, v := range v.ImportPkgs {
				_, o1 := gf.PkgImports[v]
				if !o1 { // 未导入包
					_, o := gf.ImportPkgs[k]
					if o { //包名会冲突
						kk := k + strconv.Itoa(i)
						gf.ImportPkgs[kk] = v
						gf.PkgImports[v] = kk
						i++
					} else { // 包名不会冲突
						gf.ImportPkgs[k] = v
						gf.PkgImports[v] = k
					}
				}
			}
		}
	}

	return gf
}
