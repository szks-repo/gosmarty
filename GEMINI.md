# GEMINI Project Overview: gosmarty

## Project Overview

This project is an interpreter for a Smarty-like template engine, written in Go. The module path is `github.com/szks-repo/gosmarty`.

The interpreter is composed of the following key components:

*   **Lexer (`lexer/lexer.go`):** A stateful lexer that tokenizes the input template into `TEXT` and `TAG` parts. It handles basic Smarty syntax like `{...}`, variables (`$name`), pipes (`|`), and comments (`{*...*}`).
*   **AST (`ast/ast.go`):** Defines the Abstract Syntax Tree structure. Key nodes include `TextNode`, `ActionNode` (for `{$...}`), `IfNode`, and `PipeNode`. The root of the tree is `ast.Tree`.
*   **Parser (`parser/parser.go`):** A top-down recursive descent parser that builds the AST from the token stream provided by the lexer. It handles variable tags, if/else blocks, and modifier pipelines.
*   **Object System (`object/object.go`):** A simple object system to represent evaluated values within the template, such as `String`, `Boolean`, and `Null`.
*   **Evaluator (`eval.go`):** Traverses the AST to evaluate the template. It manages a symbol table (`object.Environment`) for variables and executes built-in modifier functions.

The main entry point for using the interpreter is the `gosmarty.go` file, which provides a `New()` function to create a new `GoSmarty` instance, a `Parse()` method to parse a template string, and an `Exec()` method to evaluate the parsed template with a given environment.

## Building and Running

This is a standard Go project.

### Running Tests

The project includes a comprehensive test suite. To run all tests, use the standard Go test command from the root directory:

```sh
go test ./...
```

The tests in `gosmarty_test.go` cover variable evaluation, `if` statements, and modifier pipelines.

### Building the Project

To build the project, use the standard Go build command:

```sh
go build
```

## Development Conventions

*   **Structure:** The code is organized into distinct packages (`ast`, `lexer`, `object`, `parser`), promoting separation of concerns.
*   **Testing:** Tests are written using the standard `testing` package and are located in `_test.go` files. The existing tests make use of table-driven tests and `t.Parallel()` for efficiency.
*   **Error Handling:** The parser collects parsing errors into a slice, which can be retrieved via the `Errors()` method. This allows the parser to report multiple errors at once instead of stopping at the first one.
*   **Dependencies:** The project has a single external dependency (`golang.org/x/exp`) managed via `go.mod`.
