# Agent Notes for gosmarty

## Overview
- Go implementation of the Smarty template engine.
- Core workflow: input template -> lexer -> parser -> AST -> evaluator using environment and modifier registries.

## Repository Layout
- `gosmarty.go`: public entry points (`GoSmarty`, `Parse`, `Template.Execute`, modifier registration).
- `environment.go`: environment builder with helper options for registering variables.
- `lexer/`: stateful lexer switching between text and tag modes.
- `parser/`: recursive-descent parser building AST nodes for text, actions, if/elseif/else, pipelines, field and index access.
- `ast/`: node definitions (`ListNode`, `TextNode`, `ActionNode`, `IfNode`, `ElseIfNode`, `PipeNode`, etc.).
- `eval.go`: tree evaluator returning concatenated strings and supporting truthiness checks across object types.
- `object/`: runtime value types (`String`, `Number`, `Bool`, `Map`, `Array`, `Optional`, `Null`) plus conversion helpers (`NewObjectFromAny`).
- `modifier/`: built-in variable modifiers (`upper`, `lower`, `nl2br`, `number_format`, `wordwrap`) with registry helpers.
- `token/`: token definitions and keyword lookup.
- Tests: `gosmarty_test.go` (integration-focused), `object/object_test.go`, helper utilities in `helpers_test.go`.

## External Dependencies
- `github.com/szks-repo/go-php-functions` for modifier helpers (string utilities).
- `golang.org/x/exp` for optional types (used inside `object/`).
- `golang.org/x/text` indirectly.
- Internet access (or pre-populated `GOMODCACHE`) is required for `go test` / `go build` because modules are not vendored.

## Build & Test
- Go toolchain version in `go.mod`: 1.25 (ensure Go ≥1.21 supports this module, or adjust if needed).
- Run tests: `go test ./...` (set `GOCACHE` and `GOMODCACHE` locally if sandboxed without global cache).
- Formatting: `gofmt -w <files>`; no additional linters configured.

## Key Behaviors & Constraints
- `{if}` supports chained `{elseif}` blocks terminating with optional `{else}`; conditions currently limited to variable/identifier truthiness (no comparison operators yet).
- Variable modifiers parsed as pipelines (`{$value|upper|lower}`) with no argument parsing implemented.
- Field access (`{$user.name}`) and index access (`{$users[0]}`) supported.
- Truthiness: strings empty=false, numbers zero=false, arrays/maps empty=false, optionals unwrap recursively, null=false.
- Comments `{* ... *}` skipped entirely.

## Known Gaps / TODO Signals
- README roadmap lists `{foreach}`, `{assign}`, advanced expressions (`$num > 5`) as not yet implemented.
- Parser lacks expression grammar beyond identifiers/numbers; extending conditions will require significant work (lexer already emits comparison tokens but parser ignores them).
- Modifier arguments & broader builtin coverage pending (see comments in `modifier/modifier.go`).
- Error handling mostly accumulates parser errors and returns joined string; evaluator often returns `NULL` silently for missing bindings.

## Development Tips
- Use `WithVariable` helpers to seed environments in tests.
- Register custom modifiers via `gosmarty.RegisterModifier` in setup/test code.
- When adding syntax, update lexer `token` enums, parser, AST structs, evaluator, and expand tests in `gosmarty_test.go`.
- Keep new files ASCII encoded and run `gofmt` before committing.

## Useful Commands
- `go test ./...`
- `gofmt -w <path>`
- `go list ./...` (sanity check module paths)
- `go env GOPATH` (if you need to inspect module cache locations)

## Reference Materials
- Smarty docs: https://www.smarty.net/docs/en/ (useful when planning feature parity).
- Upstream roadmap embedded in README for prioritizing next tasks.
