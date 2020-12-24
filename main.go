package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	go func() {

	}()

	fset := token.NewFileSet()

	src, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	printDeclRecursive(f.Decls, 0, false)
}

func fprint(composite bool, level int, format string, args ...interface{}) {
	if composite {
		fmt.Printf("%*s- ", (level*2)-2, "")
		fmt.Printf(format, args...)
	} else {
		fmt.Printf("%*s", level*2, "")
		fmt.Printf(format, args...)
	}
}

func printDeclRecursive(decl interface{}, level int, composite bool) {
	switch typeddecl := decl.(type) {
	case []ast.Decl:
		for _, d := range typeddecl {
			fmt.Println("--------")
			printDeclRecursive(d, level, true && level > 0)
		}
	case []ast.Expr:
		for _, e := range typeddecl {
			printDeclRecursive(e, level, true && level > 0)
		}
	case []ast.Stmt:
		for _, e := range typeddecl {
			printDeclRecursive(e, level, true && level > 0)
		}
	case []ast.Spec:
		for _, e := range typeddecl {
			printDeclRecursive(e, level, true && level > 0)
		}
	case []*ast.Ident:
		for _, e := range typeddecl {
			printDeclRecursive(e, level, true && level > 0)
		}
	case []interface{}:
		for _, e := range typeddecl {
			printDeclRecursive(e, level, true && level > 0)
		}

	case *ast.DeclStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Decl:\n", typeddecl)
		printDeclRecursive(typeddecl.Decl, level+2, false)
	case *ast.GenDecl:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Specs:\n", typeddecl)
		printDeclRecursive(typeddecl.Specs, level+2, false)
	case *ast.GoStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Call:\n", typeddecl)
		printDeclRecursive(typeddecl.Call, level+2, false)
	case *ast.ArrayType:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(composite, level+1, "%T Elt:\n", typeddecl)
		printDeclRecursive(typeddecl.Elt, level+2, composite)
	case *ast.FuncDecl:
		fprint(composite, level, "%T name=%s\n", typeddecl, typeddecl.Name)
		fprint(composite, level+1, "%T Body:\n", typeddecl)
		printDeclRecursive(typeddecl.Body, level+2, false)
	case *ast.BlockStmt:
		fprint(composite, level, "%T:\n", typeddecl)
		fprint(composite, level+1, "%T List[]:\n", typeddecl)
		for _, l := range typeddecl.List {
			printDeclRecursive(l, level+2, true)
		}
	case *ast.ReturnStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Results[]:\n", typeddecl)
		for _, r := range typeddecl.Results {
			printDeclRecursive(r, level+2, true)
		}
	case *ast.ExprStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T X\n", typeddecl)
		printDeclRecursive(typeddecl.X, level+2, false)
	case *ast.AssignStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Lhs[]:\n", typeddecl)
		printDeclRecursive(typeddecl.Lhs, level+2, true)
		fprint(false, level+1, "%T Rhs[]:\n", typeddecl)
		printDeclRecursive(typeddecl.Rhs, level+2, true)
	case *ast.CallExpr:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Fun:\n", typeddecl)
		printDeclRecursive(typeddecl.Fun, level+2, false)
		if len(typeddecl.Args) > 0 {
			fprint(false, level+1, "%T Args[]:\n", typeddecl)
			printDeclRecursive(typeddecl.Args, level+2, true)
		}
	case *ast.SelectorExpr:
		x, okx := typeddecl.X.(*ast.Ident)
		if okx && x.Obj == nil && typeddecl.Sel.Obj == nil {
			fprint(composite, level, "%T %s.%s\n", typeddecl, x.Name, typeddecl.Sel.Name)
		} else {
			fprint(composite, level, "%T\n", typeddecl)
			fprint(false, level+1, "%T X:\n", typeddecl)
			printDeclRecursive(typeddecl.X, level+2, false)
			fprint(false, level+1, "%T Sel:\n", typeddecl)
			printDeclRecursive(typeddecl.Sel, level+2, false)
		}
	case *ast.KeyValueExpr:
		key, ok := typeddecl.Key.(*ast.Ident)
		value, ok2 := typeddecl.Value.(*ast.BasicLit)
		if ok && ok2 {
			fprint(composite, level, "ast.KeyValueExpr key=%+#v value=%v\n", key.Name, value.Value)
		} else {
			fprint(composite, level, "%T\n", typeddecl)
			fprint(false, level+1, "%T Key:\n", typeddecl)
			printDeclRecursive(typeddecl.Key, level+2, false)
			fprint(false, level+1, "%T Value:\n", typeddecl)
			printDeclRecursive(typeddecl.Value, level+2, false)
		}
	case *ast.CompositeLit:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Type:\n", typeddecl)
		printDeclRecursive(typeddecl.Type, level+2, false)
		if len(typeddecl.Elts) > 0 {
			fprint(false, level+1, "%T Elts[]:\n", typeddecl)
			printDeclRecursive(typeddecl.Elts, level+2, true)
		}
	case *ast.BasicLit:
		fprint(composite, level, "ast.BasicLit kind=%+#v value=%v\n", typeddecl.Kind.String(), typeddecl.Value)
	case *ast.TypeSwitchStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Init:\n", typeddecl)
		printDeclRecursive(typeddecl.Init, level+2, false)
		fprint(false, level+1, "%T Body:\n", typeddecl)
		printDeclRecursive(typeddecl.Body, level+2, false)
	case *ast.CaseClause:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Body:\n", typeddecl)
		printDeclRecursive(typeddecl.Body, level+2, false)
		if len(typeddecl.List) > 0 {
			fprint(false, level+1, "%T List[]:\n", typeddecl)
			printDeclRecursive(typeddecl.List, level+2, true)
		}

	// X
	case *ast.StarExpr:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T X:\n", typeddecl)
		printDeclRecursive(typeddecl.X, level+2, false)
	case *ast.ParenExpr:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T X:\n", typeddecl)
		printDeclRecursive(typeddecl.X, level+2, false)
	case *ast.IndexExpr:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T X:\n", typeddecl)
		printDeclRecursive(typeddecl.X, level+2, false)
	case *ast.TypeAssertExpr:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T X:\n", typeddecl)
		printDeclRecursive(typeddecl.X, level+2, false)

	case *ast.RangeStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Value:\n", typeddecl)
		printDeclRecursive(typeddecl.Value, level+2, false)
		fprint(false, level+1, "%T Body:\n", typeddecl)
		printDeclRecursive(typeddecl.Body, level+2, false)
	case *ast.IfStmt:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T Cond:\n", typeddecl)
		printDeclRecursive(typeddecl.Cond, level+2, false)
		fprint(false, level+1, "%T Body:\n", typeddecl)
		printDeclRecursive(typeddecl.Body, level+2, false)
		fprint(false, level+1, "%T Else:\n", typeddecl)
		printDeclRecursive(typeddecl.Else, level+2, false)
	case *ast.BinaryExpr:
		fprint(composite, level, "%T\n", typeddecl)
		fprint(false, level+1, "%T X:\n", typeddecl)
		printDeclRecursive(typeddecl.X, level+2, false)
		fprint(false, level+1, "%T Op: %s\n", typeddecl, typeddecl.Op.String())
		fprint(false, level+1, "%T Y:\n", typeddecl)
		printDeclRecursive(typeddecl.Y, level+2, false)
	case *ast.Ident:
		fprint(composite, level, "%T name=%s\n", typeddecl, typeddecl.Name)
		if typeddecl.Obj != nil {
			//fprint(false, level+1, "%T Object: %#+v\n", typeddecl, typeddecl.Obj)
		}
	case *ast.ImportSpec:
		fprint(composite, level, "%T name=%s value=%v\n", typeddecl, typeddecl.Name.String(), typeddecl.Path.Value)
		if typeddecl.Comment != nil {
			fprint(false, level+1, "%T Comment:\n", typeddecl)
			printDeclRecursive(typeddecl.Comment, level+2, false)
		}
	case *ast.ValueSpec:
		fprint(composite, level, "%T type=%s\n", typeddecl, typeddecl.Type)
		if len(typeddecl.Names) > 0 {
			fprint(false, level+1, "%T Names:\n", typeddecl)
			printDeclRecursive(typeddecl.Names, level+2, false)
		}
		if len(typeddecl.Values) > 0 {
			fprint(false, level+1, "%T Values:\n", typeddecl)
			printDeclRecursive(typeddecl.Values, level+2, false)
		}
	case *ast.Object:
		fprint(composite, level, "%T name=%s type=%v\n", typeddecl, typeddecl.Name, typeddecl.Type)
		if typeddecl.Decl != nil {
			fprint(false, level+1, "%T Decl: %+#v\n", typeddecl, typeddecl.Decl)
			printDeclRecursive(typeddecl.Decl, level+2, false)
		}
		if typeddecl.Data != nil {
			fprint(false, level+1, "%T Data: %+#v\n", typeddecl, typeddecl.Data)
		}
	case ast.Stmt:
		fprint(composite, level, "%+#v\n", typeddecl)
	case ast.Expr:
		fprint(composite, level, "%T %+#v\n", typeddecl, typeddecl)
	case ast.Decl:
		fprint(composite, level, "%T %+#v\n", typeddecl, typeddecl)
	default:
		fprint(composite, level, "%+#v\n", typeddecl)
	}
}
