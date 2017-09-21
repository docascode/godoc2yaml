// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godoc

import (
	"bytes"
	"go/doc"
	"go/printer"
)

type DocsPackage struct {
	IsMain      bool                  `json:"ismain"`
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	ImportPath  string                `json:"importPath"`
	Dir         string                `json:"dir"`
	Consts      []DocsValue           `json:"consts"`
	Types       []DocsType            `json:"types"`
	Vars        []DocsValue           `json:"vars"`
	Funcs       []DocsFunc            `json:"funcs"`
	Notes       map[string][]DocsNote `json:"notes"`
	Examples    []DocsExample         `json:"examples"`
	Dirs        []DocsDir             `json:"dirs"`
}

type DocsDir struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Summary string `json:"summary"`
	HasPkg  bool   `json:"haspkg"`
}

type DocsNote struct {
	UID         string `json:"uid"`
	Description string `json:"description"`
}

type DocsExample struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type DocsValue struct {
	Names       []string `json:"names"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
}

type DocsType struct {
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Code        string `json:"code"`

	Consts  []DocsValue `json:"consts"`
	Vars    []DocsValue `json:"vars"`
	Funcs   []DocsFunc  `json:"funcs"`
	Methods []DocsFunc  `json:"methods"`
}

type DocsFunc struct {
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

func toDocfx(p *Presentation, info *PageInfo) *DocsPackage {
	pkg := &DocsPackage{
		IsMain: info.IsMain,
		Dir:    info.Dirname,
		Notes:  toDocsNotes(info.Notes),
		Dirs:   toDocsDirs(info.Dirs),
	}
	if info.PDoc != nil {
		pkg.ImportPath = info.PDoc.ImportPath
		pkg.Summary = summary(info.PDoc.Doc)
		pkg.Description = description(info.PDoc.Doc)
		pkg.Examples = toDocsExamples(info.Examples, p, info)
		pkg.Consts = toDocsValues(info.PDoc.Consts, p, info)
		pkg.Vars = toDocsValues(info.PDoc.Vars, p, info)
		pkg.Funcs = toDocsFuncs(info.PDoc.Funcs, p, info)
		pkg.Types = toDocsTypes(info.PDoc.Types, p, info)
	}
	return pkg
}

func toDocsDirs(dirs *DirList) []DocsDir {
	if dirs == nil {
		return []DocsDir{}
	}

	arr := make([]DocsDir, len(dirs.List))
	for i, d := range dirs.List {
		arr[i] = DocsDir{
			Name:    d.Name,
			Path:    d.Path,
			Summary: d.Synopsis,
			HasPkg:  d.HasPkg,
		}
	}
	return arr
}

func toDocsTypes(types []*doc.Type, p *Presentation, info *PageInfo) []DocsType {
	arr := make([]DocsType, len(types))
	for i, t := range types {
		arr[i] = DocsType{
			Name:        t.Name,
			Summary:     summary(t.Doc),
			Description: description(t.Doc),
			Code:        p.nodeFunc(info, t.Decl),
			Consts:      toDocsValues(t.Consts, p, info),
			Vars:        toDocsValues(t.Vars, p, info),
			Funcs:       toDocsFuncs(t.Funcs, p, info),
			Methods:     toDocsFuncs(t.Methods, p, info),
		}
	}
	return arr
}

func toDocsFuncs(funcs []*doc.Func, p *Presentation, info *PageInfo) []DocsFunc {
	arr := make([]DocsFunc, len(funcs))
	for i, f := range funcs {
		arr[i] = DocsFunc{
			Name:        f.Name,
			Summary:     summary(f.Doc),
			Description: description(f.Doc),
			Code:        p.nodeFunc(info, f.Decl),
		}
	}
	return arr
}

func toDocsValues(values []*doc.Value, p *Presentation, info *PageInfo) []DocsValue {
	arr := make([]DocsValue, len(values))
	for i, v := range values {
		arr[i] = DocsValue{
			Names:       v.Names,
			Summary:     summary(v.Doc),
			Description: description(v.Doc),
			Code:        p.nodeFunc(info, v.Decl),
		}
	}
	return arr
}

func toDocsExamples(examples []*doc.Example, p *Presentation, info *PageInfo) []DocsExample {
	arr := make([]DocsExample, len(examples))
	for i, eg := range examples {
		cnode := &printer.CommentedNode{Node: eg.Code, Comments: eg.Comments}
		arr[i] = DocsExample{
			Name: eg.Name,
			Code: p.nodeFunc(info, cnode),
		}
	}
	return arr
}

func toDocsNotes(notes map[string][]*doc.Note) map[string][]DocsNote {
	m := map[string][]DocsNote{}
	for k, v := range notes {
		arr := make([]DocsNote, len(v))
		for i, n := range v {
			arr[i] = DocsNote{
				UID:         n.UID,
				Description: n.Body,
			}
		}
		m[k] = arr
	}
	return m
}

func summary(d string) string {
	return doc.Synopsis(d)
}

func description(d string) string {
	var buf bytes.Buffer
	doc.ToText(&buf, d, "", "    ", 999999)
	return buf.String()
}
