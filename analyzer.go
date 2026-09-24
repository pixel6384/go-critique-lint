package main

import (
	"go/ast"
	"go/token"
)

type Issue struct {
	Pos  token.Position
	Rule string
	Msg  string
}

func AnalyzeFile(fset *token.FileSet, file *ast.File) []Issue {
	var issues []Issue
	loopDepth := 0
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		// Track loop depth
		switch n.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			loopDepth++
		}

		for _, rule := range Rules {
			if rule.Check(n, loopDepth) {
				issues = append(issues, Issue{
					Pos:  fset.Position(n.Pos()),
					Rule: rule.Name,
					Msg:  rule.Description,
				})
			}
		}

		// We need to decrement loopDepth AFTER visiting children.
		// Since ast.Inspect is a simple traversal, we can't easily decrement here
		// without custom traversal. Let's use a manual stack-based approach or
		// refine the Inspect logic. 

		return true
	})
	return issues
}