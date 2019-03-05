// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package compiler has a set of code generation tools used in DOT
package compiler

import (
	"bytes"
	"fmt"
	"go/format"
	"golang.org/x/tools/imports"
	"sort"
	"strings"
)

// Info contians all the info needed to generate code
type Info struct {
	Package  string
	Imports  [][2]string
	Streams  []StreamInfo
	Contexts []ContextInfo
}

// StreamInfo holds the information to generate a single stream type
type StreamInfo struct {
	StreamType       string
	ValueType        string
	Fields           []FieldInfo
	EntryStreamType  string
	EntryConstructor string
}

// Generate generates the code needed to deal with a stream
func (s *StreamInfo) Generate() string {
	var result bytes.Buffer
	must(streamTpl.Execute(&result, s))
	for _, f := range s.Fields {
		var data struct {
			*StreamInfo
			*FieldInfo
		}
		data.StreamInfo = s
		data.FieldInfo = &f
		must(fieldTpl.Execute(&result, data))
	}

	if s.EntryStreamType != "" {
		must(entryTpl.Execute(&result, s))
	}

	return result.String()
}

// FieldInfo holds info on individual substream fields of the base stream
type FieldInfo struct {
	Field            string
	FieldType        string
	FieldStreamType  string
	FieldConstructor string
	FieldSubstream   string
}

// Generate returns the source code generated from the provided info
func Generate(info Info) string {
	var result bytes.Buffer
	must(headerTpl.Execute(&result, info))
	r := result.String()
	for _, s := range info.Streams {
		r += "\n" + s.Generate()
	}

	for _, c := range info.Contexts {
		r += "\n" + c.Generate()
	}

	p, err := format.Source([]byte(r))
	must(err)

	p, err = imports.Process("compiled.go", p, nil)
	must(err)

	return string(p)
}

// ParamInfo has info about arguments
type ParamInfo struct {
	Name, Type string
}

// ResultInfo has info about return values
type ResultInfo struct {
	Name, Type string
}

// ContextInfo has info about the context
type ContextInfo struct {
	ContextType   string
	Function      string
	Subcomponents []string
	Params        []ParamInfo
	Results       []ResultInfo

	Component         string
	ComponentComments string
	Method            string
	MethodComments    string
}

// Args helper
func (c *ContextInfo) Args() interface{} {
	result := []map[string]interface{}{}
	for _, p := range c.Params {
		if strings.HasSuffix(p.Name, "State") {
			continue
		}

		pType := p.Type
		if strings.HasPrefix(pType, "...") {
			pType = "[]" + pType[3:]
		}

		result = append(result, map[string]interface{}{
			"Name":    p.Name,
			"Type":    pType,
			"IsArray": strings.HasPrefix(pType, "[]"),
			"IsLast":  false,
		})
	}
	result[len(result)-1]["IsLast"] = true
	return result
}

// StateArgs helper
func (c *ContextInfo) StateArgs() interface{} {
	result := []interface{}{}
	seen := map[string]bool{}
	for _, p := range c.Params {
		if !strings.HasSuffix(p.Name, "State") {
			continue
		}

		resultName := ""
		for kk, r := range c.Results {
			rName := r.Name
			if rName == "" {
				rName = fmt.Sprintf("result%d", kk+1)
			}

			if !seen[rName] && r.Type == p.Type {
				resultName = rName
				break
			}
		}
		if resultName == "" {
			panic("Could not match result with state")
		}
		seen[resultName] = true

		result = append(result, map[string]interface{}{
			"Name":       p.Name,
			"Type":       p.Type,
			"ResultName": resultName,
		})
	}
	return result
}

// HasEllipsis helper
func (c *ContextInfo) HasEllipsis() bool {
	return strings.HasPrefix(c.Params[len(c.Params)-1].Type, "...")
}

// PkgSubcomps helper
func (c *ContextInfo) PkgSubcomps() interface{} {
	comps := map[string]map[string]bool{}
	add := func(pkg string, comp string) {
		if comps[pkg] == nil {
			comps[pkg] = map[string]bool{comp: true}
		} else {
			comps[pkg][comp] = true
		}
	}

	for _, comp := range c.Subcomponents {
		pairs := strings.SplitN(comp, ".", 2)
		if len(pairs) > 1 {
			add(pairs[0], comp)
		} else {
			add("", comp)
		}
	}
	add("", "initialized bool")
	add("", "stateHandler streams.Handler")
	for kk, arg := range c.Params {
		if kk > 0 {
			aType := arg.Type
			if strings.HasPrefix(aType, "...") {
				aType = "[]" + aType[3:]
			}
			add("memoized", arg.Name+" "+aType)
		}
	}

	for kk, r := range c.Results {
		n := r.Name
		if n == "" {
			n = fmt.Sprintf("result%d", kk+1)
		}
		add("memoized", n+" "+r.Type)
	}

	packages := []string{}
	for k := range comps {
		packages = append(packages, k)
	}
	sort.Strings(packages)
	result := []interface{}{}
	for _, pkg := range packages {
		sub := []string{}
		for k := range comps[pkg] {
			sub = append(sub, k)
		}
		sort.Strings(sub)
		result = append(result, map[string]interface{}{"Pkg": pkg, "Subcomps": sub})
	}
	return result
}

// MethodDecl helper
func (c *ContextInfo) MethodDecl() string {
	result := []string{}
	for kk, a := range c.Params {
		if kk == 0 {
			result = append(result, a.Name+"Key interface{}")
		} else if !strings.HasSuffix(a.Name, "State") {
			result = append(result, a.Name+" "+a.Type)
		}
	}
	return strings.Join(result, ", ")
}

// NonContextArgs helper
func (c *ContextInfo) NonContextArgs() string {
	result := []string{}
	for kk, a := range c.Params {
		if kk > 0 && !strings.HasSuffix(a.Name, "State") {
			result = append(result, a.Name)
		}
	}
	return strings.Join(result, ", ")
}

// AllArgs helper
func (c *ContextInfo) AllArgs() string {
	result := []string{}
	for _, a := range c.Params {
		if strings.HasSuffix(a.Name, "State") {
			result = append(result, c.Params[0].Name+".memoized."+a.Name)
		} else if strings.HasPrefix(a.Type, "...") {
			result = append(result, a.Name+"...")
		} else {
			result = append(result, a.Name)
		}
	}
	return strings.Join(result, ", ")
}

// MemoizedNonContextArgs helper
func (c *ContextInfo) MemoizedNonContextArgs() string {
	result := []string{}
	ctx := c.Params[0].Name
	for kk, a := range c.Params {
		if kk > 0 && !strings.HasSuffix(a.Name, "State") {
			result = append(result, ctx+".memoized."+a.Name)
		}
	}
	return strings.Join(result, ", ")
}

// NonContextArgsDecl helper
func (c *ContextInfo) NonContextArgsDecl() string {
	result := []string{}
	for kk, a := range c.Params {
		if kk > 0 && !strings.HasSuffix(a.Name, "State") {
			aType := a.Type
			if strings.HasPrefix(aType, "...") {
				aType = "[]" + aType[3:]
			}
			result = append(result, a.Name+" "+aType)
		}
	}
	return strings.Join(result, ", ")
}

// ResultsDecl helper
func (c *ContextInfo) ResultsDecl() string {
	seen := map[string]bool{}
	for _, a := range c.StateArgs().([]interface{}) {
		seen[a.(map[string]interface{})["ResultName"].(string)] = true
	}

	result := []string{}
	for kk, r := range c.Results {
		n := r.Name
		if n == "" {
			n = fmt.Sprintf("result%d", kk+1)
		}
		if seen[n] {
			continue
		}
		decl := n + " " + r.Type
		result = append(result, decl)
	}
	return strings.Join(result, ", ")
}

// MemoizedResults helper
func (c *ContextInfo) MemoizedResults() string {
	ctx := c.Params[0].Name
	result := []string{}
	for kk, r := range c.Results {
		n := r.Name
		if n == "" {
			n = fmt.Sprintf("result%d", kk+1)
		}
		result = append(result, ctx+".memoized."+n)
	}
	return strings.Join(result, ", ")
}

// MemoizedNonStateResults helper
func (c *ContextInfo) MemoizedNonStateResults() string {
	seen := map[string]bool{}
	for _, a := range c.StateArgs().([]interface{}) {
		seen[a.(map[string]interface{})["ResultName"].(string)] = true
	}

	ctx := c.Params[0].Name
	result := []string{}
	for kk, r := range c.Results {
		n := r.Name
		if n == "" {
			n = fmt.Sprintf("result%d", kk+1)
		}
		if seen[n] {
			continue
		}
		result = append(result, ctx+".memoized."+n)
	}
	return strings.Join(result, ", ")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// Generate generates the code needed to deal with a context
func (c *ContextInfo) Generate() string {
	var result bytes.Buffer
	must(contextTpl.Execute(&result, c))
	return result.String()
}
