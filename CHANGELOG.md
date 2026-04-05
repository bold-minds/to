# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `PtrOr[T any](v T, fallback *T) *T` — value-to-pointer counterpart of `ValOr`. Returns a pointer to v when v is not the zero value of T, and returns fallback otherwise. Zero-value detection uses `reflect.Value.IsZero`, which is defined for every Go type (primitives, structs, arrays, nil pointers/maps/slices/channels/funcs/interfaces). The primary use case is building JSON patch/update payloads where unset fields must be omitted via `omitempty`. Completes the symmetry with the existing `Val`/`ValOr` pair: `Ptr`/`PtrOr` handles the T → *T direction with a caller-supplied fallback for the zero case.
- `Ptr[T any](v T) *T` — ergonomic pointer creation. Returns a non-nil pointer to a fresh copy of v. Equivalent to Go 1.26+'s `new(v)` builtin, but available on all supported Go versions of this package.

### Changed
- Lowered Go floor from 1.26 to 1.21 (`go.mod`, README, CHANGELOG, workflow matrix, validate.sh).
- README `Pointer creation` section rewritten to document `to.Ptr` instead of directing users to Go 1.26's `new(v)`.
- CI workflow split into `test` matrix (read-only permissions) and `badges` job (write, main-only) for least-privilege.
- `golangci-lint` pinned to `v2.0.2` instead of `@latest` in CI and validate.sh.
- `validate.sh`: portable `grep -oE | sed` instead of GNU-only `grep -oP`; SC2155 fix; numeric sanitization on alert-count API responses.
- `.golangci.yml`: dropped dead settings blocks for linters not in enable list.
- `codeql.yml`: enabled push/PR/weekly schedule triggers.

## [0.1.0] — Initial release

### Added
- **Outcome-named conversions** (return zero on failure):
  - `Str(v any) string`
  - `Int(v any) int`
  - `Bool(v any) bool`
  - `F64(v any) float64`
- **Fallback variants** (return fallback on failure):
  - `StrOr(v any, fallback string) string`
  - `IntOr(v any, fallback int) int`
  - `BoolOr(v any, fallback bool) bool`
  - `F64Or(v any, fallback float64) float64`
- **Generic conversion** with structured error context:
  - `Type[T any](v any) (T, error)`
  - `TypeOr[T any](v any, fallback T) T`
- **Pointer dereferencing** (safe nil handling):
  - `Val[T any](p *T) T` — zero value on nil
  - `ValOr[T any](p *T, fallback T) T` — custom default on nil
- **`ConversionError` type** with From/To/Value/Reason/Cause fields and `Unwrap()` support for `errors.Is`/`errors.As`
- Support for converting between all numeric types (`int`, `int8`/16/32/64, `uint`, `uint8`/16/32/64, `float32`, `float64`), `bool`, and `string` — including round-tripping through strings for the standard forms (`true`/`1`/`yes`/`on`, `false`/`0`/`no`/`off`)
- 97.7% test coverage
- Zero allocations on happy paths (sub-5 ns/op for common conversions)
- Zero external dependencies — pure stdlib

### Deliberate non-goals
- **No per-type variants** for every numeric width (no `Int32`, `Int8`, `Uint`, etc.) — use `Type[T]` for uncommon targets
- **No `Must*` panicking variants** — use the outcome-named functions or `Or` fallback variants

### Requires
- Go 1.21 or later

*(Note: the `Ptr(v)` function was added in the Unreleased section after the initial 0.1.0 draft, when the family's minimum Go version was lowered from 1.26 to 1.21. On Go 1.26+, the `new(v)` builtin is equivalent.)*
