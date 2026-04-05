# to

[![Go Reference](https://pkg.go.dev/badge/github.com/bold-minds/to.svg)](https://pkg.go.dev/github.com/bold-minds/to)
[![Build](https://img.shields.io/github/actions/workflow/status/bold-minds/to/test.yaml?branch=main&label=tests)](https://github.com/bold-minds/to/actions/workflows/test.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bold-minds/to)](go.mod)
[![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/bold-admin/d173f3c3709c1ac53f0dd6d38d2cbac0/raw/coverage.json)](https://github.com/bold-minds/to/actions/workflows/test.yaml)

**Safe Go value conversion — any value, any type, one call.**

Go's type system is strict for good reasons, but real-world values arrive untyped or wrong-typed: environment variables are strings, JSON numbers are `float64`, config fields are `any`. `to` gives you one-line conversion to the type you actually need, with fallbacks when conversion fails.

```go
// Before
port := 8080
if s := os.Getenv("PORT"); s != "" {
    if p, err := strconv.Atoi(s); err == nil {
        port = p
    }
}

// After
port := to.IntOr(os.Getenv("PORT"), 8080)
```

## ✨ Why to?

- 🎯 **Outcome-named** — `to.Int(x)` reads like English, no `[T]` at the call site for common types
- 🔁 **Consistent fallback pattern** — every conversion has a plain form and an `Or` form with a default
- 🧭 **Generic escape hatch** — `to.Type[T](x)` handles any target type when the outcome-named shortcuts don't cover it
- 🧨 **Rich error context** — `ConversionError` tells you the source type, target type, value, and reason
- 🛡️ **Nil-safe everywhere** — `nil` inputs convert predictably, no panics
- 🪶 **Tiny** — 13 functions, one file, zero dependencies
- 🔗 **Pairs with [`bold-minds/dig`](https://github.com/bold-minds/dig)** — dig into nested data, then convert the leaf

## 📦 Installation

```bash
go get github.com/bold-minds/to
```

Requires Go 1.21 or later.

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "os"

    "github.com/bold-minds/to"
)

func main() {
    // Env vars are always strings — convert with a fallback
    port := to.IntOr(os.Getenv("PORT"), 8080)
    host := to.StrOr(os.Getenv("HOST"), "localhost")
    debug := to.BoolOr(os.Getenv("DEBUG"), false)

    fmt.Println(port, host, debug)

    // Config values are often any
    cfg := map[string]any{"timeout": 30.0, "name": "alice"}
    timeout := to.IntOr(cfg["timeout"], 60) // float64 → int, fallback 60
    name := to.StrOr(cfg["name"], "unknown")

    fmt.Println(timeout, name)

    // Generic escape hatch for uncommon types
    bytes, err := to.Type[int64]("4096")
    if err != nil {
        // handle conversion failure
    }
    _ = bytes

    // Safe pointer dereferencing
    var namePtr *string
    display := to.ValOr(namePtr, "anonymous")
    fmt.Println(display)
}
```

## 🔧 Core Features

### Outcome-named conversions

Four common target types get outcome-named shortcuts. No `[T]` at the call site — the function name IS the target type.

```go
s := to.Str(x)    // any → string
i := to.Int(x)    // any → int
b := to.Bool(x)   // any → bool
f := to.F64(x)    // any → float64
```

Each returns the zero value on failure. If you want a different fallback, use the `Or` variant:

```go
s := to.StrOr(x, "default")
i := to.IntOr(x, 0)
b := to.BoolOr(x, false)
f := to.F64Or(x, 0.0)
```

> ⚠️ **`StrOr` is nil-only.** `IntOr`, `BoolOr`, and `F64Or` fall back whenever conversion fails (including empty strings). `StrOr` is the exception: it returns the fallback **only when `v` is nil**, because `Str` never fails on non-nil input. If you reach for `to.StrOr(cfg["name"], "anonymous")` expecting an empty string to trigger `"anonymous"`, it won't — use an explicit check:
>
> ```go
> name := to.Str(cfg["name"])
> if name == "" {
>     name = "anonymous"
> }
> ```

### Generic conversion

For target types not covered by the outcome-named shortcuts — `int64`, `float32`, `time.Duration`, user-defined named numeric types, etc. — reach for `Type[T]`:

```go
id, err := to.Type[int64](userMap["id"])
if err != nil { /* handle */ }

bytes, err := to.Type[int64]("4096")
if err != nil { /* handle */ }

// time.Duration accepts Go duration strings ("5s", "1h30m") or numeric sources.
timeout, err := to.Type[time.Duration]("5s")

// User-defined named numeric types are supported via the reflect fallback,
// with range-checking against the underlying kind.
type Port uint16
p, err := to.Type[Port]("8080")
```

`Type[T]` returns `(T, error)` so you get full error context on failure. For the same "fallback on failure" ergonomics as the outcome-named funcs, use `TypeOr[T]`:

```go
id := to.TypeOr[int64](userMap["id"], 0)
timeout := to.TypeOr[int64]("4096", 60)
```

### Rich error context with `ConversionError`

When conversion fails and you need to know *why*, `Type[T]` returns a `*ConversionError` with structured details:

```go
_, err := to.Type[int]("not a number")
if err != nil {
    var ce *to.ConversionError
    if errors.As(err, &ce) {
        fmt.Printf("could not convert %s(%v) to %s: %s\n",
            ce.From, ce.Value, ce.To, ce.Reason)
        // could not convert string(not a number) to int: invalid numeric string
    }
}
```

Fields on `ConversionError`: `From` (source type name), `To` (target type name), `Value` (the original value), `Reason` (human-readable cause), `Cause` (underlying error, if any). Implements `error` and `Unwrap()`, so `errors.Is` and `errors.As` work as expected.

### Pointer dereferencing

Go's strict nil-safety means dereferencing a nullable pointer requires boilerplate. `to.Val` and `to.ValOr` collapse it:

```go
// Without
var name string
if namePtr != nil {
    name = *namePtr
}

// With
name := to.Val(namePtr)                 // zero value on nil
name := to.ValOr(namePtr, "anonymous")  // custom default on nil
```

### Pointer creation

Go's strict nil-safety means creating a pointer to a literal or freshly computed value requires a temp variable:

```go
// Without
x := 8080
p := &x

// With
p := to.Ptr(8080)
```

`to.Ptr` returns a pointer to a fresh copy of its argument. It works for any type — primitives, strings, slices, maps, structs, and interfaces.

```go
port := to.Ptr(8080)          // *int
name := to.Ptr("alice")       // *string
now  := to.Ptr(time.Now())    // *time.Time
user := to.Ptr(User{ID: 1})   // *User
```

Each call returns a distinct pointer — mutating `*p` does not affect the source. `Ptr` is never nil, even for zero values: `to.Ptr("")` returns a non-nil `*string` pointing at an empty string.

On Go 1.26+, the `new(v)` builtin provides the same functionality at the language level. `to.Ptr` exists to give Bold Minds code the same ergonomic on Go 1.21–1.25 without depending on a third-party package.

## 🛡️ Safety guarantees

- **Never panics.** `nil` inputs, wrong types, unparseable strings, and nil pointers all return zero values, fallbacks, or `ConversionError` — never a panic.
- **No float→int trap values.** `NaN`, `±Inf`, and out-of-range floats are rejected rather than silently becoming `INT64_MIN` (the default Go behavior for `int(math.NaN())` on amd64).
- **No silent float32 overflow.** `Type[float32]` of a value exceeding `math.MaxFloat32` returns a `ConversionError` rather than `±Inf`.
- **32-bit-safe.** `Type[int]` range-checks `int64` and `uint32` sources against `math.MaxInt` / `math.MinInt`, so conversions behave correctly on `GOARCH=386`/`arm`/`wasm` as well as 64-bit platforms.
- **Immutable.** `to` never modifies input values.
- **Zero dependencies.** Pure stdlib.
- **No reflection on the happy path.** Outcome-named conversions and the common `Type[T]` targets use concrete type switches; reflection is used only for the named-numeric-type fallback (`time.Duration`, `type Port int64`, …).

## 🏎️ Performance

Measured on Go 1.26 (Intel Ultra 9 275HX; library targets Go 1.21+). Happy paths are sub-5 nanosecond with zero allocations; pointer creation and dereferencing are effectively free.

```
BenchmarkStr_String-24         920353114    0.66 ns/op    0 B/op    0 allocs/op
BenchmarkStr_Int-24             20948056   27.17 ns/op    2 B/op    1 allocs/op
BenchmarkInt_FromString-24     149880273    3.82 ns/op    0 B/op    0 allocs/op
BenchmarkInt_FromFloat-24      263979162    2.07 ns/op    0 B/op    0 allocs/op
BenchmarkInt_FromInt-24        817980475    0.79 ns/op    0 B/op    0 allocs/op
BenchmarkBool_FromString-24    226217947    2.39 ns/op    0 B/op    0 allocs/op
BenchmarkBool_FromInt-24       261220177    2.00 ns/op    0 B/op    0 allocs/op
BenchmarkF64_FromString-24      44144103   13.63 ns/op    0 B/op    0 allocs/op
BenchmarkF64_FromInt-24        260099806    2.14 ns/op    0 B/op    0 allocs/op
BenchmarkIntOr_Success-24      132154460    4.24 ns/op    0 B/op    0 allocs/op
BenchmarkType_Int-24           165019142    3.77 ns/op    0 B/op    0 allocs/op
BenchmarkVal-24               1000000000    0.13 ns/op    0 B/op    0 allocs/op
BenchmarkValOr_NonNil-24      1000000000    0.13 ns/op    0 B/op    0 allocs/op
BenchmarkValOr_Nil-24         1000000000    0.13 ns/op    0 B/op    0 allocs/op
```

Error paths allocate (the `ConversionError` struct plus `fmt.Sprintf` for the error message). If you care about error-path performance, prefer the outcome-named functions or their `Or` variants — they skip error construction entirely.

## 🧪 Testing

```bash
go test ./...                      # unit tests
go test -race ./...                # race detection
go test -bench=. -benchmem ./...   # benchmarks
```

Current coverage: 97.7%.

## 📚 API Reference

```go
// Outcome-named conversions — return zero on failure.
func Str(v any) string
func Int(v any) int
func Bool(v any) bool
func F64(v any) float64

// Outcome-named conversions with fallback.
func StrOr(v any, fallback string) string
func IntOr(v any, fallback int) int
func BoolOr(v any, fallback bool) bool
func F64Or(v any, fallback float64) float64

// Generic conversion. Returns (value, nil) on success,
// (zero, *ConversionError) on failure.
func Type[T any](v any) (T, error)

// Generic conversion with fallback. Returns fallback on any failure.
func TypeOr[T any](v any, fallback T) T

// Ptr returns a pointer to a fresh copy of v. Ergonomic one-call
// alternative to the two-line "x := v; &x" pattern.
func Ptr[T any](v T) *T

// Safe pointer dereferencing. Val returns zero on nil pointer;
// ValOr returns fallback on nil pointer.
func Val[T any](p *T) T
func ValOr[T any](p *T, fallback T) T

// ConversionError is returned by Type[T] when conversion fails.
// Implements error and Unwrap() for errors.Is/errors.As compatibility.
type ConversionError struct {
    From   string // Source type name
    To     string // Target type name
    Value  any    // Original value
    Reason string // Human-readable reason
    Cause  error  // Underlying error, if any
}
```

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bold Minds Go libraries follow a shared set of design principles; read [PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md) before opening a PR.

## 📄 License

MIT. See [LICENSE](LICENSE).

## 🔗 Related Projects

- [`bold-minds/dig`](https://github.com/bold-minds/dig) — nested data navigation. Pairs naturally with `to`: `to.Int(dig.At(data, "count"))` when you need to dig into `any` trees and convert the leaf.
- [`spf13/cast`](https://github.com/spf13/cast) — the established Go conversion library. Predates generics; uses `ToInt`/`ToIntE`-style naming. `to` is a leaner, generics-first, outcome-named take on the same problem space.
- Go 1.26 `new(v)` builtin — the language-level equivalent of `to.Ptr(v)` on 1.26+. Both work, pick whichever matches your minimum Go version.
- [`k8s.io/utils/ptr`](https://pkg.go.dev/k8s.io/utils/ptr) — the kubernetes ecosystem's widely-used pointer helper (`ptr.To`, `ptr.Deref`). `to.Ptr`/`to.Val`/`to.ValOr` cover the same ground with an outcome-naming convention consistent across Bold Minds libraries.
- Go standard library `strconv` — the low-level primitive for string ↔ numeric conversions. `to` wraps it with outcome naming and `any` input handling.
