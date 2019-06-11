package parse

import (
	"go/ast"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
)

func (infor *File) ParseFile(file *ast.File) {
	infor.ParseImport(file)

	inters := make([]Interface, 0)
	ast.Inspect(file, func(x ast.Node) bool {
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
			f := infor.ParseFunc(decl.Name.Name, funcType)
			inter.Funcs = []Func{f}
			inter.IsFunc = true
			inters = append(inters, inter)
		case *ast.TypeSpec:
			s := x.(*ast.TypeSpec)
			inter.Name = s.Name.Name
			s1, ok1 := s.Type.(*ast.InterfaceType)
			if !ok1 {
				return true
			}
			log.Info("find interface:", inter.Name)
			inter.Funcs = make([]Func, len(s1.Methods.List))
			for i, field := range s1.Methods.List {
				f := infor.ParseFunc(field.Names[0].Name, field.Type.(*ast.FuncType))
				inter.Funcs[i] = f
			}
			inters = append(inters, inter)
		default:
			return true
		}
		return false
	})
	infor.Interfaces = inters

}

func (infor *File) ParseImport(file *ast.File) {
	imports := make(map[string]string)
	ast.Inspect(file, func(x ast.Node) bool {
		switch x.(type) {
		case *ast.ImportSpec:
			importSpec := x.(*ast.ImportSpec)
			var key string
			value := importSpec.Path.Value
			value, _ = strconv.Unquote(value)
			if importSpec.Name != nil {
				key = importSpec.Name.Name
			} else {
				lastIndex := strings.LastIndex(value, "/")
				if lastIndex == -1 {
					key = value
				} else {
					key = value[lastIndex+1:]
				}
			}
			imports[key] = value
		default:
			return true
		}
		return false
	})
	infor.Imports = imports
}

func (infor *File) ParseStruct(name string, st *ast.StructType) Struct {
	s := Struct{Fields: make([]Field, 0)}
	fields := infor.ParseField(st.Fields.List)
	for _, field := range fields {
		if strings.Title(field.Name) == field.Name {
			field.IsField = true
			s.Fields = append(s.Fields, field)
		}
	}
	s.Name = name
	s.Pkg = infor.PkgPath
	return s
}

// 解析ast函数
func (infor *File) ParseFunc(name string, funcType *ast.FuncType) Func {
	fun := Func{}
	// 函数名
	fun.Name = name
	// 函数参数

	if funcType.Params != nil {
		fun.Params = infor.ParseField(funcType.Params.List)
	}
	// 函数返回值
	if funcType.Results != nil {
		fun.Results = infor.ParseField(funcType.Results.List)
	}
	return fun
}

// 解析ast方法声明中的表达式
func ParseExpr(expr ast.Expr) (typeName string) {
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
		return typeName
	}
}

// 解析ast字段列表
func (infor *File) ParseField(fieldList []*ast.Field) []Field {
	ls := make([]Field, 0, len(fieldList))
	for _, field := range fieldList {
		typeName := ParseExpr(field.Type)
		protoType := infor.parseType(typeName)

		var name string
		name = strings.Replace(typeName, "[]", "", -1)
		name = strings.Replace(name, "*", "", -1)

		var mtype = name
		var pkgSort string
		if strings.Contains(typeName, ".") {
			index := strings.Index(typeName, name)
			s := typeName[0:index]
			lastIndex := strings.Index(name, ".")
			if lastIndex != -1 {
				mtype = "#" + name[lastIndex:]
				pkgSort = name[:lastIndex]
			}
			mtype = s + mtype
		}
		pkg, ok := infor.Imports[pkgSort]
		if !ok || pkg == "" {
			pkg = infor.PkgPath
		}
		_, o := goBaseType[name]
		if o {
			pkg = ""
		}

		if field.Names == nil || len(field.Names) == 0 {
			var name string
			name = strings.Replace(typeName, "[", "", -1)
			name = strings.Replace(name, "]", "", -1)
			name = strings.Replace(name, "*", "", -1)
			name = strings.Replace(name, ".", "", -1)
			name = strings.Replace(name, "{", "", -1)
			name = strings.Replace(name, "}", "", -1)
			name += "0"
			fieldname := toExportField(name)
			p := Field{
				Name:      name,
				FieldName: fieldname,
				GoType:    typeName,
				Package:   pkg,
				ProtoType: protoType,
			}
			ls = append(ls, p)
		} else {
			for _, name := range field.Names {
				fieldname := toExportField(name.Name)
				p := Field{
					Name:      name.Name,
					FieldName: fieldname,
					GoType:    typeName,
					Package:   pkg,
					ProtoType: protoType,
				}
				ls = append(ls, p)
			}
		}
	}
	return ls
}

// 解析ast字段列表
func FieldListParse(fieldList []*ast.Field) []Field {
	ls := make([]Field, 0, len(fieldList))
	for _, field := range fieldList {
		typeName := ParseExpr(field.Type)
		if field.Names == nil || len(field.Names) == 0 {
			var name string
			name = strings.Replace(typeName, "[", "", -1)
			name = strings.Replace(name, "]", "", -1)
			name = strings.Replace(name, "*", "", -1)
			name = strings.Replace(name, ".", "", -1)
			name = strings.Replace(name, "{", "", -1)
			name = strings.Replace(name, "}", "", -1)
			name += "0"
			fieldname := toExportField(name)
			p := Field{
				Name:      name,
				FieldName: fieldname,
				GoType:    typeName,
			}
			ls = append(ls, p)
		} else {
			for _, name := range field.Names {
				fieldname := toExportField(name.Name)
				p := Field{
					Name:      name.Name,
					FieldName: fieldname,
					GoType:    typeName,
				}
				ls = append(ls, p)
			}
		}
	}
	return ls
}

func toExportField(name string) string {
	split := strings.Split(name, "_")
	titles := make([]string, len(split))
	for k, v := range split {
		titles[k] = strings.Title(v)
	}
	name = strings.Join(titles, "")
	return name
}

// 生成导入包
func genImport(k, v string) string {
	sort := pkgForSort(v)
	if k == sort {
		return `"` + v + `"`
	} else {
		return k + ` "` + v + `"`
	}

}

func pkgForSort(p string) string {
	index := strings.LastIndex(p, "/")
	pkgForSort := p
	if index != -1 {
		pkgForSort = p[index+1:]
	}
	return pkgForSort
}
