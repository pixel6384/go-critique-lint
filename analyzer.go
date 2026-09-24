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
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	// Track loop depth
	switch n.(type) {
	case *ast.ForStmt, *ast.RangeStmt:
		v.loopDepth++
	}

	for _, rule := range Rules {
		if rule.Check(n, v.loopDepth) {
			v.issues = append(v.issues, Issue{
				Pos:  v.fset.Position(n.Pos()),
				Rule: rule.Name,
				Msg:  rule.Description,
			})
		}
	}

	// Capture the current depth to decrement after children are visited
	currentDepth := v.loopDepth
	
	// To correctly decrement loopDepth, we need to wrap the traversal
	// or handle it in a way that happens after children. 
	// Since we return v, ast.Walk continues. To decrement, we must 
	// wrap the children visitation. A simpler way for this tool's scope
	// is to use a custom visitor for loops.

	if _, ok := n.(*ast.ForStmt); ok {
		return &loopDecrementer{v, currentDepth}
	}
	if _, ok := n.(*ast.RangeStmt); ok {
		return &loopDecrementer{v, currentDepth}
	}

	return v
}

type loopDecrementer struct {
	*visitor
	initialDepth int
}

func (ld *loopDecrementer) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		ld.visitor.loopDepth = ld.initialDepth
		return nil
	}
	return ld.visitor.Visit(n)
}

func AnalyzeFile(fset *token.FileSet, file *ast.File) []Issue {
	v := &visitor{
		fset:   fset,
		issues: []Issue{},
	}
	ast.Walk(v, file)
	return v.issues
}