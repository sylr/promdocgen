package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
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
	switch decl.(type) {
	case []ast.Decl:
		decls := decl.([]ast.Decl)
		for _, d := range decls {
			fmt.Println("--------")
			printDeclRecursive(d, level, true && level > 0)
		}
	case []ast.Expr:
		exprs := decl.([]ast.Expr)
		for _, e := range exprs {
			printDeclRecursive(e, level, true && level > 0)
		}
	case []interface{}:
		exprs := decl.([]interface{})
		for _, e := range exprs {
			printDeclRecursive(e, level, true && level > 0)
		}
	case *ast.ArrayType:
		arr := decl.(*ast.ArrayType)
		fprint(composite, level, "%T\n", arr)
		fprint(composite, level+1, "%T Elt\n", arr)
		printDeclRecursive(arr.Elt, level+2, composite)
	case *ast.FuncDecl:
		fun := decl.(*ast.FuncDecl)
		fprint(composite, level, "%T name=%s\n", fun, fun.Name)
		fprint(composite, level+1, "%T Body:\n", fun)
		printDeclRecursive(fun.Body, level+2, false)
	case *ast.BlockStmt:
		block := decl.(*ast.BlockStmt)
		fmt.Printf("%*s%T[]:\n", level*2, "", block)
		for _, l := range block.List {
			printDeclRecursive(l, level+1, true)
		}
	case *ast.ReturnStmt:
		ret := decl.(*ast.ReturnStmt)
		fprint(composite, level, "%T\n", ret)
		fprint(false, level+1, "%T Results[]:\n", ret)
		for _, r := range ret.Results {
			printDeclRecursive(r, level+2, true)
		}
	case *ast.ExprStmt:
		expr := decl.(*ast.ExprStmt)
		fprint(composite, level, "%+#v\n", expr)
		fprint(false, level+1, "ast.ExprStmt X\n")
		printDeclRecursive(expr.X, level+2, false)
	case *ast.CallExpr:
		cexpr := decl.(*ast.CallExpr)
		fprint(composite, level, "%T\n", cexpr)
		fprint(false, level+1, "%T Fun:\n", cexpr)
		printDeclRecursive(cexpr.Fun, level+2, false)
		fprint(false, level+1, "%T Args[]:\n", cexpr)
		printDeclRecursive(cexpr.Args, level+2, true)
	case *ast.SelectorExpr:
		sexpr := decl.(*ast.SelectorExpr)

		x, okx := sexpr.X.(*ast.Ident)
		if okx && x.Obj == nil && sexpr.Sel.Obj == nil {
			//fmt.Printf("%*sast.SelectorExpr %s.%s\n", level*2, "", x.Name, sexpr.Sel.Name)
			fprint(composite, level, "%T %s.%s\n", sexpr, x.Name, sexpr.Sel.Name)
		} else {
			fprint(composite, level, "%T\n", sexpr)
			fprint(composite, level+1, "%T X:\n", sexpr)
			printDeclRecursive(sexpr.X, level+2, false)
			fprint(composite, level+1, "%T Sel:\n", sexpr)
			printDeclRecursive(sexpr.Sel, level+2, false)
		}
	case *ast.KeyValueExpr:
		kvexpr := decl.(*ast.KeyValueExpr)

		key, ok := kvexpr.Key.(*ast.Ident)
		value, ok2 := kvexpr.Value.(*ast.BasicLit)
		if ok && ok2 {
			fprint(composite, level, "ast.KeyValueExpr key=%+#v value=%v\n", key.Name, value.Value)
		} else {
			fprint(composite, level, "%T\n", kvexpr)
			fprint(false, level+1, "%T Key:\n", kvexpr)
			printDeclRecursive(kvexpr.Key, level+2, false)
			fprint(false, level+1, "%T Value:\n", kvexpr)
			printDeclRecursive(kvexpr.Value, level+2, false)
		}
	case *ast.CompositeLit:
		compol := decl.(*ast.CompositeLit)
		fprint(composite, level, "ast.CompositeLit\n")
		fprint(false, level+1, "ast.CompositeLit Type:\n")
		printDeclRecursive(compol.Type, level+2, false)
		fprint(false, level+1, "ast.CompositeLit Elts[]:\n")
		printDeclRecursive(compol.Elts, level+2, true)
	case *ast.BasicLit:
		basicl := decl.(*ast.BasicLit)
		fprint(composite, level, "ast.BasicLit kind=%+#v value=%v\n", basicl.Kind.String(), basicl.Value)
	case *ast.Ident:
		ident := decl.(*ast.Ident)
		if ident.Obj != nil {
			fprint(composite, level, "%T name=%s\n", ident, ident.Name)
			printDeclRecursive(ident.Obj, level+1, false)
		} else {
			fprint(composite, level, "%+#v\n", ident)
		}
	case *ast.Object:
		object := decl.(*ast.Object)
		fprint(composite, level, "%+#v\n", object)
	case ast.Stmt:
		stat := decl.(ast.Stmt)
		fprint(composite, level, "%+#v\n", stat)
	case ast.Expr:
		expr := decl.(ast.Expr)
		fprint(composite, level, "%T %+#v\n", expr, expr)
	case ast.Decl:
		decl := decl.(ast.Decl)
		fprint(composite, level, "%T %+#v\n", decl, decl)
	default:
		fprint(composite, level, "%T %+#v\n", decl, decl)
	}
}
