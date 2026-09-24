package main

import "go/ast"

// Rule defines a pattern to look for in the AST and the suggested fix.
type Rule struct {
	Name        string
	Description string
	Check       func(ast.Node) bool
}

var Rules = []Rule{
	{
		Name:        "AvoidMagicNumbers",
		Description: "Avoid using magic numbers; use named constants instead.",
		Check: func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != 2 { // 2 is token.INT
				return false
		}
			// Simple check: ignore 0, 1, -1
			val := lit.Value
			return val != "0" && val != "1" && val != "-1"
		},
	},
	{
		Name:        "SlicePreallocation",
		Description: "Use make([]T, 0, capacity) when the final size is known to avoid multiple allocations.",
		Check: func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return false
		}
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "make" {
				if len(call.Args) == 1 {
					return true
				}
			}
			return false
		},
	},
}