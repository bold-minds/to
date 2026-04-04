# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- **No `Ptr(v)` function** — Go 1.26's `new(v)` builtin replaces it at the language level
- **No per-type variants** for every numeric width (no `Int32`, `Int8`, `Uint`, etc.) — use `Type[T]` for uncommon targets
- **No `Must*` panicking variants** — use the outcome-named functions or `Or` fallback variants

### Requires
- Go 1.26 or later
