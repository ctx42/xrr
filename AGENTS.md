# Project Rules for xrr

Instructions for AI coding assistants.

## Commands

- All tests: `go test -v -race ./...`
- Single: `go test -v -race -run TestName ./pkg/xrr/`
- Subpackage: `go test -v -race ./pkg/xrr/xrrtest/`
- Lint: `golangci-lint run ./...`
- Regenerate README doc examples: `gomake :project:doc-eg`

CI runs full suite on master. Min Go 1.26+.

## Architecture

All source under `pkg/xrr/`. Only external dep: `github.com/ctx42/testing`.

### Core Types

Generics with `Domain` constraint (`comparable`):

- `GenericError[T Domain]`: main error (msg, code, meta, wrapped). Implements
  error, Coder, Metadater, json.Marshaler/Unmarshaler.
- `GenericFields[T Domain]`: `map[string]error` for field validation.
  Implements Fielder + json.Marshaler.
- `ErrorFunc[T]()` / `ErrorfFunc[T]()`: domain-scoped constructors.
- `WrapUsing[T](err, opts...)`: annotate without changing message (prefer
  New+WithCause publicly).
- Default domain `EDXrr` powers `Error`, `Fields`, `New`, `NewFieldError`
  aliases in xrr.go.

User domains should use an unexported marker (e.g. `edPayment struct{}`) for
compile-time isolation. The package default `EDXrr` is exported.

### Interfaces

- `Coder`: `ErrorCode() string`
- `Fielder`: `ErrorFields() map[string]error`
- `Metadater`: `MetaAll() map[string]any`

### Traversal

`inspect.go` provides `walk`/`walkReverse` (BFS) over chains, Fielder fields,
and joined errors. `Get*` helpers use them. GetMeta walks reverse (root wins).

### Key Files

- `meta.go`: Metadata builder + `.Option()`
- `options.go`: WithCode / WithMeta / WithMetaFrom / WithCause
- `envelope.go`: API JSON envelope (cause + optional lead)
- `errors.go`, `sync_errors.go`: []error collections (thread-safe variant)
- `helpers.go`: Split/Join/IsJoined, DefaultCode, AssertDomain[T]
- `xrrtest/`: Assert* helpers on ctx42/testing framework
- `masked.go`: Mask/Masked (safe public error + internal cause chain)

### Sentinels

`sentinel.go`: ErrInvJSON (ECInvJSON), ErrInvJSONError, ErrFields (ECFields).
ECGeneric is the default for errors without an explicit code.

### Test Support

`all_test.go`: TErrorFields, TErrMarshalJSON, TMetaAll + tree builders
(TstTreeCase1..5, TstTreeMeta).

## Coding Conventions

### Headers

Every .go file starts with:
```go
// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT
```

### Tests

- `_test.go` matches its production file.
- **Test funcs appear in same order** as the methods they cover in prod code.
- Use `github.com/ctx42/testing/pkg/assert` (order: t, want, have).
- Preferred subtest (exact markers):
```go
t.Run("desc", func(t *testing.T) {
    // --- Given ---
    // --- When ---
    // --- Then ---
})
```

### TODO / FIXME

Never delete `// TODO`, `// FIXME`, `// HACK`, or `// OPTIMIZE` without
completing the task. godox linter flags them.

### Markdown Tables (docs/comments)

Align columns with spaces. Separator dashes match widest cell. Pad cells to
column width. No trailing spaces after final `|`.

### Style

- Receivers: 1-3 letters, consistent per type.
- Prefer `any`.
- Functional options (`With*`).
- Compile-time checks: `var _ Coder = (*T)(nil)`

### Editor Configuration

This project uses a minimal `.editorconfig` for basic cross-editor
consistency:

- UTF-8, LF line endings, final newline required, no trailing whitespace.
- Go: tabs (actual formatting controlled by `gofmt`/`gofumpt` + linters).
- YAML: 2-space indent.
- Markdown (including `AGENTS.md`): 80-column guidance.

See the root `.editorconfig` for details.

## Linting & Tooling

Key linters: errcheck, errorlint, exhaustive, forcetypeassert, gocritic, godox,
gosec.

forbidigo blocks `fmt.Print*` / `print` (allowed in examples + tests).

## New Domains

```go
type edPayment struct{} // unexported marker

type PaymentError = xrr.GenericError[edPayment]
var NewPaymentError = xrr.ErrorFunc[edPayment]()
```
Check with `xrr.IsDomain[edPayment](err)` or `errors.As`.

User-defined domains should use unexported marker types. The built-in
default domain `EDXrr` is exported for internal aliases.

## Documentation

- Run `gomake :project:doc-eg` to sync README examples with `Example*` funcs
  in examples_test.go (uses gmdoceg markers).
- Update CHANGELOG.md for every user-visible change.

## Editing This File

Keep concise and actionable. Specific commands/rules > general advice.
Subdirectory `AGENTS.md` files take precedence.
