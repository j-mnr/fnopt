package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	withFieldFuncFmt = `func With%[1]s(v %[2]s) func(*%[3]s) {
	return func(t *%[3]s) {
		t.%[1]s = v
	}
}

`
	newXFuncFmt = `func New%[1]s(options ...func(*%[2]s)) *%[2]s {
	v := &%[2]s{}
	for _, o := range options {
		o(v)
	}
	return v
}

`
)

var (
	typeName = flag.String("type", "", "type name; must be set")
	fileName = flag.String("file", "", "full relative path to file; must be set")
)

type Generator struct{ buf bytes.Buffer }

func (g *Generator) Printf(format string, args ...any) (int, error) {
	return fmt.Fprintf(&g.buf, format, args...)
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage of fnopt:
	fnopt -type T -file path

Example usage:
	fnopt -type Actor -file ./actor/actor.go

Flags:
`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("fnopts: ")
	flag.Usage = usage
	flag.Parse()
	if len(*typeName) == 0 || len(*fileName) == 0 {
		flag.Usage()
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *fileName, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	var fields []*ast.Field
	ast.Inspect(f, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}

		for _, spec := range decl.Specs {
			if tspec, ok := spec.(*ast.TypeSpec); ok &&
				tspec.Type != nil && tspec.Name.Name == *typeName {
				if stype, ok := tspec.Type.(*ast.StructType); ok {
					for _, field := range stype.Fields.List {
						i := 0
						for _, n := range field.Names {
							if !n.IsExported() {
								continue
							}
							field.Names[i] = n
							i++
						}
						field.Names = field.Names[:i]
						if len(field.Names) == 0 {
							continue
						}
						fields = append(fields, field)
					}
				}
				return false
			}
		}
		return false
	})
	if len(fields) == 0 {
		log.SetOutput(os.Stderr)
		log.Fatalf("error: type does not exist or has no exported fields\n")
	}

	var g Generator
	g.Printf("package " + f.Name.Name + "\n\n")
	g.Printf(newXFuncFmt, strings.Title(*typeName), *typeName)
	for _, fld := range fields {
		for _, name := range fld.Names {
			g.Printf(withFieldFuncFmt, strings.Title(name.Name), types.ExprString(fld.Type), *typeName)
		}
	}

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Fatal(err, "\n\n", g.buf.String())
	}
	if err := os.WriteFile(path.Dir(*fileName)+"/"+toSnakeCase(*typeName)+"_fnopt.go",
		src, 0o0644); err != nil {
		log.Fatal(err)
	}
}

func toSnakeCase[S ~string](s S) S {
	r, n := utf8.DecodeRuneInString(string(s))
	var sb strings.Builder
	sb.WriteRune(unicode.ToLower(r))
	for _, r := range s[n:] {
		if unicode.IsUpper(r) {
			sb.WriteByte('_')
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return S(sb.String())
}
