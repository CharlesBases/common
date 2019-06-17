package parse

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
)

func (file *File) ParseFile(astFile *ast.File) {
	file.ParseImport(astFile)
	inters := make([]Interface, 0)
	ast.Inspect(astFile, func(x ast.Node) bool {
		inter := Interface{}
		switch x.(type) {
		case *ast.FuncDecl:
			decl := x.(*ast.FuncDecl)
			if decl.Recv != nil {
				return true
			}
			if decl.Name.Name[0] != strings.ToUpper(decl.Name.Name)[0] {
				return true
			}
			inter.Name = decl.Name.Name + "Func"
			log.Info("find func:", inter.Name)
			funcType := decl.Type
			fun := file.ParseFunc(decl.Name.Name, funcType)
			inter.Funcs = []Func{fun}
			inter.IsFunc = true
			inters = append(inters, inter)
		case *ast.TypeSpec:
			typeSpec := x.(*ast.TypeSpec)
			inter.Name = typeSpec.Name.Name
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return true
			}
			log.Info("find interface:", inter.Name)
			inter.Funcs = make([]Func, len(interfaceType.Methods.List))
			for index, field := range interfaceType.Methods.List {
				fun := file.ParseFunc(field.Names[0].Name, field.Type.(*ast.FuncType))
				inter.Funcs[index] = fun
			}
			inters = append(inters, inter)
		default:
			return true
		}
		return false
	})
	file.Interfaces = inters
}

func (file *File) ParseImport(astFile *ast.File) {
	ImportA := make(map[string]string)
	ast.Inspect(astFile, func(x ast.Node) bool {
		switch x.(type) {
		case *ast.ImportSpec:
			importSpec := x.(*ast.ImportSpec)
			var key string
			val := importSpec.Path.Value
			val, _ = strconv.Unquote(val)
			if importSpec.Name != nil {
				key = importSpec.Name.Name
			} else {
				lastIndex := strings.LastIndex(val, "/")
				if lastIndex == -1 {
					key = val
				} else {
					key = val[lastIndex+1:]
				}
			}
			ImportA[key] = val
		default:
			return true
		}
		return false
	})
	file.ImportA = ImportA
	file.ImportB = make(map[string]string, 0)
	for key, val := range file.ImportA {
		file.ImportB[val] = key
	}
}

func (file *File) ParseStruct(name string, structType *ast.StructType) Struct {
	s := Struct{Fields: make([]Field, 0)}
	fields := file.ParseField(structType.Fields.List)
	for _, field := range fields {
		if strings.Title(field.Name) == field.Name {
			field.IsField = true
			s.Fields = append(s.Fields, field)
		}
	}
	s.Name = name
	s.Pkg = file.PkgPath
	return s
}

// 解析ast函数
func (file *File) ParseFunc(name string, funcType *ast.FuncType) Func {
	fun := Func{}
	// 函数名
	fun.Name = name
	// 函数参数

	if funcType.Params != nil {
		fun.Params = file.ParseField(funcType.Params.List)
	}
	// 函数返回值
	if funcType.Results != nil {
		fun.Results = file.ParseField(funcType.Results.List)
	}
	// 第一个参数是否是context
	if len(fun.Params) >= 1 {
		t := fun.Params[0].Package
		if t == "context" || t == "golang.org/x/net/context" {
			s := file.ImportB[t]
			if fun.Params[0].GoType == s+".Context" {
				fun.Params = fun.Params[1:]
			}
		}
	}
	return fun
}

// 解析ast方法声明中的表达式
func ParseExpr(expr ast.Expr) (fieldType string) {
	switch expr.(type) {
	case *ast.StarExpr:
		starExpr := expr.(*ast.StarExpr)
		return "*" + ParseExpr(starExpr.X)
	case *ast.SelectorExpr:
		selectorExpr := expr.(*ast.SelectorExpr)
		return ParseExpr(selectorExpr.X) + "." + selectorExpr.Sel.Name
	case *ast.ArrayType:
		arrayType := expr.(*ast.ArrayType)
		return "[]" + ParseExpr(arrayType.Elt)
	case *ast.MapType:
		mapType := expr.(*ast.MapType)
		return "map[" + ParseExpr(mapType.Key) + "]" + ParseExpr(mapType.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ident:
		ident := expr.(*ast.Ident)
		return ident.Name
	default:
		return fieldType
	}
}

// 解析ast字段列表
func (file *File) ParseField(astField []*ast.Field) []Field {
	fields := make([]Field, 0, len(astField))
	for _, field := range astField {
		fieldType := ParseExpr(field.Type)
		protoType := file.parseType(fieldType)

		name := strings.Replace(fieldType, "[]", "", -1)
		name = strings.Replace(name, "*", "", -1)

		var variableType = name
		var pkgSort string
		if strings.Contains(fieldType, ".") {
			index := strings.Index(fieldType, name)
			str := fieldType[0:index]
			lastIndex := strings.Index(name, ".")
			if lastIndex != -1 {
				variableType = "#" + name[lastIndex:]
				pkgSort = name[:lastIndex]
			}
			variableType = str + variableType
		}
		imp, ok := file.ImportA[pkgSort]
		if !ok || imp == "" {
			imp = file.PkgPath
		}
		_, ok = goBaseType[name]
		if ok {
			imp = ""
		}

		if field.Names == nil || len(field.Names) == 0 {
			var name string
			name = strings.Replace(fieldType, "[", "", -1)
			name = strings.Replace(name, "]", "", -1)
			name = strings.Replace(name, "*", "", -1)
			name = strings.Replace(name, ".", "", -1)
			name = strings.Replace(name, "{", "", -1)
			name = strings.Replace(name, "}", "", -1)
			name += "0"
			fieldname := title(name)
			field := Field{
				Name:         name,
				FieldName:    fieldname,
				Variable:     fieldname,
				VariableType: variableType,
				GoType:       fieldType,
				Package:      imp,
				ProtoType:    protoType,
			}
			fields = append(fields, field)
		} else {
			for _, name := range field.Names {
				fieldname := title(name.Name)
				field := Field{
					Name:         name.Name,
					FieldName:    name.Name,
					Variable:     fieldname,
					VariableType: variableType,
					GoType:       fieldType,
					Package:      imp,
					ProtoType:    protoType,
				}
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// 解析ast字段列表
// func ParseFieleList(astField []*ast.Field) []Field {
// 	fields := make([]Field, 0, len(astField))
// 	for _, field := range astField {
// 		fieldType := ParseExpr(field.Type)
// 		if field.Names == nil || len(field.Names) == 0 {
// 			var name string
// 			name = strings.Replace(fieldType, "[", "", -1)
// 			name = strings.Replace(name, "]", "", -1)
// 			name = strings.Replace(name, "*", "", -1)
// 			name = strings.Replace(name, ".", "", -1)
// 			name = strings.Replace(name, "{", "", -1)
// 			name = strings.Replace(name, "}", "", -1)
// 			name += "0"
// 			fieldname := title(name)
// 			p := Field{
// 				Name:      name,
// 				FieldName: fieldname,
// 				GoType:    fieldType,
// 			}
// 			fields = append(fields, p)
// 		} else {
// 			for _, name := range field.Names {
// 				fieldname := title(name.Name)
// 				p := Field{
// 					Name:      name.Name,
// 					FieldName: fieldname,
// 					GoType:    fieldType,
// 				}
// 				fields = append(fields, p)
// 			}
// 		}
// 	}
// 	return fields
// }

func title(name string) string {
	builder := strings.Builder{}
	for _, val := range strings.Split(name, "_") {
		builder.WriteString(strings.Title(val))
	}
	return builder.String()
}

// 生成导入包
func genImport(key, val string) string {
	sort := packageSort(val)
	if key == sort {
		return fmt.Sprintf(`"%s"`, val)
	} else {
		return fmt.Sprintf(`%s "%s"`, key, val)
	}
}

func packageSort(Package string) string {
	index := strings.LastIndex(Package, "/")
	packageSort := Package
	if index != -1 {
		packageSort = Package[index+1:]
	}
	return packageSort
}
