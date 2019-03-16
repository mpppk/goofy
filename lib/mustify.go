package lib

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/loader"

	goofyast "github.com/mpppk/goofy/ast"
)

func GenerateErrorWrappersFromProgram(filePath string) (*ast.File, []ast.Decl, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get abs file path")
	}

	prog, err := goofyast.NewProgram(filePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load program file")
	}

	pkg, file, ok := findPkgAndFileFromProgram(prog, absFilePath)
	if !ok {
		return nil, nil, errors.New("file not found: " + filePath)
	}

	var newDecls []ast.Decl
	importDecls := extractImportDeclsFromDecls(file.Decls)
	newDecls = append(newDecls, importDecls...)
	exportedFuncDecls := extractExportedFuncDeclsFromDecls(file.Decls)
	errorWrappers := funcDeclsToErrorFuncWrappers(exportedFuncDecls, pkg)
	newDecls = append(newDecls, errorWrappers...)
	return file, newDecls, nil
}

func extractImportDeclsFromDecls(decls []ast.Decl) (importDecls []ast.Decl) {
	for _, decl := range decls {
		if importDecl, ok := declToImportDecl(decl); ok {
			importDecls = append(importDecls, importDecl)
		}
	}
	return
}

func declToImportDecl(decl ast.Decl) (*ast.GenDecl, bool) {
	if genDecl, ok := decl.(*ast.GenDecl); ok {
		if genDecl.Tok == token.IMPORT {
			return genDecl, true
		}
	}
	return nil, false
}

func extractExportedFuncDeclsFromDecls(decls []ast.Decl) (funcDecls []*ast.FuncDecl) {
	for _, decl := range decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if ast.IsExported(funcDecl.Name.Name) {
				funcDecls = append(funcDecls, funcDecl)
			}
		}
	}
	return
}

func funcDeclsToErrorFuncWrappers(funcDecls []*ast.FuncDecl, pkg *loader.PackageInfo) (newDecls []ast.Decl) {
	for _, funcDecl := range funcDecls {
		//newDecl, ok := goofyast.ConvertErrorFuncToMustFunc(prog, pkg, funcDecl)
		if newDecl, ok := goofyast.GenerateErrorFuncWrapper(pkg, funcDecl); ok {
			newDecls = append(newDecls, newDecl)
		}
	}
	return
}

func findPkgAndFileFromProgram(prog *loader.Program, targetAbsFilePath string) (*loader.PackageInfo, *ast.File, bool) {
	for _, pkg := range prog.Created {
		for _, file := range pkg.Files {
			currentFilePath := prog.Fset.File(file.Pos()).Name()
			if absCurrentFilePath, err := filepath.Abs(currentFilePath); err == nil {
				if targetAbsFilePath == absCurrentFilePath {
					return pkg, file, true
				}
			}

		}
	}
	return nil, nil, false
}

func WriteAstFile(filePath string, file *ast.File) error {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to get abs file path: "+absFilePath)
	}

	f, err := os.Create(absFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to create file: "+absFilePath)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	if err := format.Node(f, token.NewFileSet(), file); err != nil {
		return errors.Wrap(err, "failed to write ast file to  "+absFilePath)
	}
	return nil
}
