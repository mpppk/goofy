package lib

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/loader"

	goofyast "github.com/mpppk/goofy/ast"
)

func DeclToImportDecl(decl ast.Decl) (*ast.GenDecl, bool) {
	if genDecl, ok := decl.(*ast.GenDecl); ok {
		if genDecl.Tok == token.IMPORT {
			return genDecl, true
		}
	}
	return nil, false
}

func GenerateErrorFuncWrappers(funcDecls []*ast.FuncDecl, pkg *loader.PackageInfo) (newDecls []ast.Decl) {
	for _, funcDecl := range funcDecls {
		if !ast.IsExported(funcDecl.Name.Name) {
			continue
		}

		newDecl, ok := goofyast.GenerateErrorFuncWrapper(pkg, funcDecl)
		//newDecl, ok := goofyast.ConvertErrorFuncToMustFunc(prog, pkg, funcDecl)
		if !ok {
			continue
		}
		newDecls = append(newDecls, newDecl)
	}
	return
}
