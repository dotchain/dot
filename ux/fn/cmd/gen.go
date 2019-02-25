// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package main generates go context spec based on specific functional components
package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

func main() {
	structName := ""

	args := os.Args[1:]
	if len(args) >= 3 {
		structName = args[0]
		args = args[1:]
	}

	fbytes, err := ioutil.ReadFile(args[1])
	if err != nil {
		panic(err)
	}

	output := process(readCode(structName, args[0], args[1], fbytes))
	dir := path.Dir(args[1])
	err = ioutil.WriteFile(
		path.Join(dir, structName+args[0]+"Cache.go"),
		[]byte(output),
		os.ModePerm,
	)
	if err != nil {
		panic(err)
	}
}

func readCode(structName, typeName, fileName string, fbytes []byte) map[string]interface{} {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		panic(err)
	}

	decl := getComponentDecl(f, structName, typeName)
	ctxtName, ctxtType := getContextNameAndType(decl, fbytes)
	args, argTypes := namesAndTypes(decl.Type.Params, fbytes)
	_, resultTypes := namesAndTypes(decl.Type.Results, fbytes)

	argsAndTypes := ""
	for kk, arg := range args {
		if kk == 0 {
			continue
		}
		if kk > 1 {
			argsAndTypes += ", "
		}
		argsAndTypes += arg + " " + argTypes[kk]
	}

	ellipsis := ""
	if strings.HasPrefix(argTypes[len(argTypes)-1], "...") {
		ellipsis = "..."
	}

	raw, sub := getSubPackages(decl, ctxtName, fbytes)
	return map[string]interface{}{
		"stateStructName":   structName,
		"component":         typeName,
		"cache":             typeName + "Cache",
		"package":           f.Name.Name,
		"imports":           getImports(f),
		"contextStructName": ctxtType,
		"argsAndTypes":      argsAndTypes,
		"args":              strings.Join(args[1:], ",") + ellipsis,
		"resultTypes":       strings.Join(resultTypes, ","),
		"rawSubComponents":  raw,
		"packages":          sub,
	}
}

func process(data map[string]interface{}) string {
	t, err := template.New("template").Parse(code)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		panic(err)
	}
	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	p, err = imports.Process("test.go", p, nil)
	if err != nil {
		panic(err)
	}

	return string(p)
}

func getImports(f *ast.File) [][]string {
	imports := [][]string{}
	for _, spec := range f.Imports {
		im := []string{"", ""}
		if spec.Name != nil {
			im[0] = spec.Name.Name
		}
		im[1] = spec.Path.Value
		imports = append(imports, im)
	}
	return imports
}

func getComponentDecl(f *ast.File, structName string, typeName string) *ast.FuncDecl {
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)

		if !ok || fn.Name == nil || fn.Name.Name != typeName {
			continue
		}

		if structName == "" && fn.Recv == nil {
			return fn
		}
		if structName != "" && fn.Recv != nil {
			firstType, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				panic("expected receiver to be pointer: " + structName + "." + typeName)
			}

			id, ok := firstType.X.(*ast.Ident)
			if !ok {
				panic(firstType.X)
			}

			if id.Name == structName {
				return fn
			}
		}
	}
	panic("function not found")
}

func getNameAndType(f *ast.Field, fbytes []byte) (string, string) {
	name := ""
	if f.Names != nil {
		name = f.Names[0].Name
	}

	return name, string(fbytes[f.Type.Pos()-1 : f.Type.End()-1])
}

func getContextNameAndType(d *ast.FuncDecl, fbytes []byte) (string, string) {
	n, t := getNameAndType(d.Type.Params.List[0], fbytes)
	if t[:1] != "*" {
		panic("expect pointer type for context: " + t)
	}
	return n, t[1:]
}

func namesAndTypes(p *ast.FieldList, fbytes []byte) ([]string, []string) {
	if p == nil {
		return nil, nil
	}
	names, types := []string{}, []string{}
	for _, f := range p.List {
		n, t := getNameAndType(f, fbytes)
		names = append(names, n)
		types = append(types, t)
	}
	return names, types
}

func getSubPackages(f *ast.FuncDecl, ctxtName string, fbytes []byte) ([]string, []interface{}) {
	visitor := &subVisitor{ctxtName, map[[2]string]bool{}}
	ast.Walk(visitor, f.Body)
	rawPackages := []string{}
	packages := map[string][]string{}
	for pair := range visitor.subs {
		if pair[0] == "" {
			rawPackages = append(rawPackages, pair[1])
		} else {
			packages[pair[0]] = append(packages[pair[0]], pair[1])
		}
	}
	subPackages := []interface{}{}
	for name, sub := range packages {
		subPackages = append(subPackages, map[string]interface{}{
			"name":          name,
			"subComponents": sub,
		})
	}
	return rawPackages, subPackages
}

type subVisitor struct {
	ctxtName string
	subs     map[[2]string]bool
}

func (s *subVisitor) Visit(n ast.Node) ast.Visitor {
	if call, ok := n.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			inner := sel.Sel.Name
			if s.isContext(sel.X) {
				s.subs[[2]string{"", inner}] = true
			}
			if sel, ok := sel.X.(*ast.SelectorExpr); ok {
				outer := sel.Sel.Name
				if s.isContext(sel.X) {
					s.subs[[2]string{outer, inner}] = true
				}
			}
		}
	}
	return s
}

func (s *subVisitor) isContext(e ast.Expr) bool {
	if ident, ok := e.(*ast.Ident); ok && ident.Name == s.ctxtName {
		return true
	}
	return false
}

const code = `
// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.
//
//
// This code is generated by github.com/dotchain/dot/ux/fn/cmd/gen.go

package {{.package}}

import (
{{- range $m := .imports}}
	{{index $m 0}} {{index $m 1}}
{{- end}}
)

// {{.contextStructName}} is the context struct needed for {{.component}}
type {{.contextStructName}} struct {
	{{.stateStructName}}
	{{range $sub := .rawSubComponents}}{{$sub}}Cache
	{{end}}{{range $pkg := .packages}}
	{{$pkg.name}} struct { {{range $sub := $pkg.subComponents}}
		{{$pkg.name}}.{{$sub}}Cache{{end}}
        }{{end}}
}

// {{.cache}} implements a cache of {{.component}} controls
type {{.cache}} struct {
	old, current map[interface{}]*{{.contextStructName}}
}

// Begin starts the round
func (c *{{.cache}}) Begin() {
	c.old, c.current = c.current, map[interface{}]*{{.contextStructName}}{}
}

// End ends the round
func (c *{{.cache}}) End() {
	// TODO: deliver Close() handlers if they exist
	c.old = nil
}

// {{.component}} implements the cache create or fetch method
func (c *{{.cache}}) {{.component}}(key interface{}, {{.argsAndTypes}}) ({{.resultTypes}}) {
	ctx, ok := c.old[key]

	if ok {
		delete(c.old, key)
	} else {
		ctx = &{{.contextStructName}}{}
	}
	c.current[key] = ctx

	{{range $sub := .rawSubComponents}}ctx.{{$sub}}Cache.Begin()
	defer ctx.{{$sub}}Cache.End()
	{{end}}{{range $pkg := .packages}}{{range $sub := $pkg.subComponents}}ctx.{{$pkg.name}}.{{$sub}}Cache.Begin()
	defer ctx.{{$pkg.name}}.{{$sub}}Cache.End(){{end}}{{end}}

	{{if eq .stateStructName ""}}return {{.component}}(ctx, {{.args}})
	{{else}}return ctx.{{.stateStructName}}.{{.component}}(ctx, {{.args}}){{end}}
}
`
