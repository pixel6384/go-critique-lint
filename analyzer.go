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

type visitor struct {
	fset      *token.FileSet
	issues    []Issue
	loopDepth int
	ifDepth   int
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	// Track depths
	switch n.(type) {
	case *ast.ForStmt, *ast.RangeStmt:
		v.loopDepth++
	case *ast.IfStmt:
		v.ifDepth++
	}

	for _, rule := range Rules {
		if rule.Check(n, v.loopDepth, v.ifDepth) {
			v.issues = append(v.issues, Issue{
				Pos:  v.fset.Position(n.Pos()),
				Rule: rule.Name,
				Msg:  rule.Description,
			})
		}
	}

	// Handle depth decrementing via wrapping visitors
	if _, ok := n.(*ast.ForStmt); ok {
		return &depthDecrementer{v, "loop"}
	}
	if _, ok := n.(*ast.RangeStmt); ok {
		return &depthDecrementer{v, "loop"}
	}
	if _, ok := n.(*ast.IfStmt); ok {
		return &depthDecrementer{v, "if"}
	}

	return v
}

type depthDecrementer struct {
	*visitor
	typeOfDepth string
}

func (dd *depthDecrementer) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		if dd.typeOfDepth == "loop" {
			dd.visitor.loopDepth--
		} else {
			dd.visitor.ifDepth--
		}
		return nil
	}
	return dd.visitor.Visit(n)
}

func AnalyzeFile(fset *token.FileSet, file *ast.File) []Issue {
	v := &visitor{
		fset:   fset,
		issues: []Issue{},
	}
	ast.Walk(v, file)
	return v.issues
}