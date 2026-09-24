package main

import (
	"go/ast"
	"go/token"
)

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
			if !ok || lit.Kind != token.INT {
				return false
			}
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
	{
		Name:        "UseAnyInsteadOfEmptyInterface",
		Description: "Use 'any' instead of 'interface{}' for better readability (Go 1.18+).",
		Check: func(n ast.Node) bool {
			iface, ok := n.(*ast.InterfaceType)
			if !ok {
				return false
			}
			return iface.Methods == nil || len(iface.Methods.List) == 0
		},
	},
	{
		Name:        "PotentialNilDereference",
		Description: "Variable is used in a method call without a preceding nil check in the current block.",
		Check: func(n ast.Node) bool {
			// This is a simplified check: look for selector expressions (x.Method())
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return false
			}
			// We check if the X part is an identifier
			_, ok = sel.X.(*ast.Ident)
			return ok
			// Note: A full implementation would require data-flow analysis
			// For this lint tool, we flag potential points for manual review
		},
	},
}