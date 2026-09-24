# go-critique-lint

A lightweight static analysis tool for Go that identifies non-idiomatic patterns and "code smells".

## Features
- Detects magic numbers in source code.
- Suggests slice preallocation for efficiency.
- AST-based analysis using Go's native `go/ast` package.

## Usage
```bash
go run . path/to/your/file.go
```