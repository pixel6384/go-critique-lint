package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-critique-lint <file.go>")
		os.Exit(1)
	}

	fset := token.NewFileSet()
	path := os.Args[1]
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	issues := AnalyzeFile(fset, node)
	if len(issues) == 0 {
		fmt.Println("No issues found!")
		return
	}

	fmt.Printf("Found %d issues in %s:\n", len(issues), path)
	for _, issue := range issues {
		fmt.Printf("%s:%d: [%s] %s\n", issue.Pos.Filename, issue.Pos.Line, issue.Rule, issue.Msg)
	}
}