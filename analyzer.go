package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Issue struct {
	Pos  token.Position
	Rule string
	Msg  string
}

func AnalyzeFile(fset *token.FileSet, file *ast.File) []Issue {
	var issues []Issue
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		for _, rule := range Rules {
			if rule.Check(n) {
				issues = append(issues, Issue{
					Pos:  fset.Position(n.Pos()),
					Rule: rule.Name,
					Msg:  rule.Description,
				})
			}
		}
		return true
	})
	return issues
}