// Package parser provides a pseudo parser that parses arikawa api files and
// maps them to Methods.
// It does not actually parse the api files, but rather utilizes the output
// generated by go's built-in parser, hence pseudo parser.
package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/pkg/errors"

	"github.com/mavolin/dismock/tools/codegen/mock/config"
)

type Parser struct {
	fparsers []*fileParser
}

func New(apiPath string) (*Parser, error) {
	entries, err := os.ReadDir(apiPath)
	if err != nil {
		return nil, err
	}

	p := &Parser{fparsers: make([]*fileParser, 0, len(entries))}

	for _, entry := range entries {
		switch {
		case entry.IsDir():
		case entry.Name() == "api.go":
		case strings.HasSuffix(entry.Name(), "_test.go"):
		case strings.HasPrefix(entry.Name(), "_"):
		default:
			fp := newFileParser(apiPath + "/" + entry.Name())
			p.fparsers = append(p.fparsers, fp)

			i := sort.Search(len(p.fparsers), func(i int) bool {
				return p.fparsers[i].file >= fp.file
			})

			if i < len(p.fparsers)-1 {
				copy(p.fparsers[i+1:], p.fparsers[i:])
				p.fparsers[i] = fp
			}
		}
	}

	return p, nil
}

func (p *Parser) Parse() ([]File, error) {
	files := make([]File, 0, len(p.fparsers))

	for _, fp := range p.fparsers {
		f, err := fp.parse()
		if err != nil {
			return nil, err
		}

		if len(f.Methods) != 0 {
			files = append(files, *f)
		}
	}

	return files, nil
}

func newFileParser(file string) *fileParser {
	return &fileParser{file: file}
}

type fileParser struct {
	file string
}

func (p *fileParser) parse() (*File, error) {
	log.Printf("reading %s\n", p.file)

	src, err := parser.ParseFile(token.NewFileSet(), p.file, nil, 0)
	if err != nil {
		return nil, err
	}

	f := File{
		Name:    p.file[strings.LastIndex(p.file, "/")+1:],
		Methods: make([]Method, 0, len(src.Decls)),
	}

	for _, decl := range src.Decls {
		// filter out unrelated stuff

		fun, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// unexported
		if unicode.IsLower([]rune(fun.Name.Name)[0]) {
			continue
		}

		// not a method
		if fun.Recv == nil {
			continue
		}

		// not a method with pointer receiver
		recvType, ok := fun.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			continue
		}

		// there are multiple pointers
		recvElemType, ok := recvType.X.(*ast.Ident)
		if !ok {
			continue
		}

		// not a api.Client method
		if recvElemType.Name != "Client" {
			continue
		}

		// method is marked as excluded
		if funConfig := config.C[fun.Name.Name]; funConfig.Exclude {
			continue
		}

		mp := newMethodParser(fun)
		method, err := mp.Parse()
		if err != nil {
			return nil, errors.Wrap(err, f.Name)
		}

		f.Methods = append(f.Methods, *method)
	}

	return &f, nil
}