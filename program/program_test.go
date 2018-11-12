package program

import (
	"fmt"
	"go/ast"
	"strings"
	"testing"

	"github.com/logrusorgru/aurora"
)

func PrintCommentText(name interface{}, comments []*ast.CommentGroup) {
	str := strings.Replace(GetTextFromCommentGroup(comments), "\n", "//", -1)
	fmt.Printf("%s: %s\n", name, aurora.Sprintf(aurora.Red(str)))
}

const FIXTURES = "github.com/wusphinx/gin-swagger/program/fixtures"

func TestProgram_CommentGroupFor(t *testing.T) {
	pkgComments := FIXTURES + "/comments"
	p := NewProgram(pkgComments, false)
	pkgInfo := p.Program.Package(pkgComments)

	for _, file := range pkgInfo.Files {
		PrintCommentText(file.Name, p.CommentGroupFor(file))

		for _, decl := range file.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				genDecl := decl.(*ast.GenDecl)
				PrintCommentText(genDecl.Tok, p.CommentGroupFor(genDecl))

				for _, spec := range genDecl.Specs {
					switch spec.(type) {
					case *ast.ValueSpec:
						valueSpec := spec.(*ast.ValueSpec)
						PrintCommentText(valueSpec.Names, p.CommentGroupFor(valueSpec))
					case *ast.ImportSpec:
						importSpec := spec.(*ast.ImportSpec)
						PrintCommentText(importSpec.Path.Value, p.CommentGroupFor(importSpec))
					case *ast.TypeSpec:
						typeSpec := spec.(*ast.TypeSpec)
						PrintCommentText(typeSpec.Name, p.CommentGroupFor(typeSpec))

						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							for _, field := range structType.Fields.List {
								PrintCommentText(field.Names, p.CommentGroupFor(field))
							}
						}
					}
				}
			case *ast.FuncDecl:
				funcDecl := decl.(*ast.FuncDecl)
				PrintCommentText(funcDecl.Name, p.CommentGroupFor(funcDecl))

				for _, stmt := range funcDecl.Body.List {
					PrintCommentText(stmt, p.CommentGroupFor(stmt))

					if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
						PrintCommentText("assign Lhs", p.CommentGroupFor(assignStmt.Lhs[0]))
						PrintCommentText("assign Rhs", p.CommentGroupFor(assignStmt.Rhs[0]))
					}

					if ifStmt, ok := stmt.(*ast.IfStmt); ok {
						for _, stmt := range ifStmt.Body.List {
							PrintCommentText(stmt, p.CommentGroupFor(stmt))
						}
					}
				}
			}

		}
	}
}

func indirectSelectorX(selectorExpr *ast.SelectorExpr) *ast.Ident {
	if s, ok := selectorExpr.X.(*ast.SelectorExpr); ok {
		return indirectSelectorX(s)
	}
	return selectorExpr.X.(*ast.Ident)
}
