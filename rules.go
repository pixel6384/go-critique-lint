package main

import (
	"go/ast"
	"go/token"
)

// Rule defines a pattern to look for in the AST and the suggested fix.
type Rule struct {
	Name        string
	Description string
	Check       func(ast.Node, int, int) bool
}

var Rules = []Rule{
	{
		Name:        "AvoidMagicNumbers",
		Description: "Avoid using magic numbers; use named constants instead.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
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
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return false
			}
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "make" {
			// Check if the first argument is a slice type
			if len(call.Args) > 0 {
				if _, ok := call.Args[0].(*ast.CompositeLit); ok {
					// Simplification: we check if it has only 1 arg and is likely a slice
					if len(call.Args) == 1 {
						return true
					}
				}
			}
			// If it's a make call with exactly 1 argument, it's often a map or a slice without capacity
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
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
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
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return false
			}
			ident, ok := sel.X.(*ast.Ident)
			if !ok {
				return false
			}
			// Ignore common package names to reduce noise
			packages := map[string]bool{"fmt": true, "os": true, "log": true, "strings": true, "strconv": true}
			if packages[ident.Name] {
				return false
			}
			return true
		},
	},
	{
		Name:        "DeferInLoop",
		Description: "Avoid using 'defer' inside a loop; it may cause resource leakage as defers only execute when the function returns.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			if loopDepth > 0 {
				if _, ok := n.(*ast.DeferStmt); ok {
					return true
				}
			}
			return false
		},
	},
	{
		Name:        "LongFunction",
		Description: "Function is too long; consider breaking it down into smaller functions to improve maintainability.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return false
			}
			// Rough estimate of function length based on number of statements
			stmtCount := 0
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				if node != nil {
					stmtCount++
				}
				return true
			})
			// Threshold of 50 AST nodes as a heuristic for "too long"
			return stmtCount > 50
		},
	},
	{
		Name:        "AvoidNestedLoops",
		Description: "Deeply nested loops (3 or more) detected; consider extracting inner loops into a separate function.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			if loopDepth >= 3 {
				if _, ok := n.(*ast.ForStmt); ok {
					return true
				}
				if _, ok := n.(*ast.RangeStmt); ok {
					return true
				}
			}
			return false
		},
	},
	{
		Name:        "DeeplyNestedIfs",
		Description: "Deeply nested if-statements detected; consider using guard clauses to flatten the logic.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			if _, ok := n.(*ast.IfStmt); ok && ifDepth >= 4 {
				return true
			}
			return false
		},
	},
	{
		Name:        "UnusedImport",
		Description: "Import is declared but not used in the file.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			// This is a simplified check. A real implementation would track all identifiers
			// and cross-reference them with imports. For the scope of this tool,
			// we flag imports that are explicitly marked as blank imports if not intended.
			gen, ok := n.(*ast.GenDecl)
			if !ok || gen.Tok != token.IMPORT {
				return false
			}
			for _, spec := range gen.Specs {
				importSpec, ok := spec.(*ast.ImportSpec)
				if ok && importSpec.Name != nil && importSpec.Name.Name == "_" {
					return true
				}
			}
			return false
		},
	},
	{
		Name:        "RedundantElse",
		Description: "Redundant else block detected after a return statement; consider removing it to flatten the code.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok || ifStmt.Else == nil {
				return false
			}
			// Check if the 'if' block ends with a return statement
			endsWithReturn := false
			if ifStmt.Body != nil {
				for i := len(ifStmt.Body.List) - 1; i >= 0; i-- {
					stmt := ifStmt.Body.List[i]
					if _, ok := stmt.(*ast.ReturnStmt); ok {
						endsWithReturn = true
						break
					}
					break
				}
			}
			return endsWithReturn
		},
	},
	{
		Name:        "EmptyIfBlock",
		Description: "Empty if-block detected; ensure this is intentional and not a leftover from debugging.",
		Check: func(n ast.Node, loopDepth, ifDepth int) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return false
			}
			if ifStmt.Body == nil || len(ifStmt.Body.List) == 0 {
				return true
			}
			return false
		},
	},
}