# Contributing to `to`

Thanks for your interest in contributing. This guide covers the operational process. For the **why** — the design principles every contribution is tested against — see **[bold-minds/oss/PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md)**.

## 🎯 Before You Start

Every contribution is measured against the four Bold Minds principles: **outcome naming**, **one way to do each thing**, **get out of the way**, and **non-goals explicit**. If your proposed change doesn't honor these, it will not be merged.

**Read [PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md) first.** It's the load-bearing document.

## 🔧 Development Setup

**Requirements:** Go 1.21 or later, Git, Bash.

```bash
git clone https://github.com/bold-minds/to.git
cd to
go test ./...              # unit tests
go test -race ./...        # race detection
go test -bench=. ./...     # benchmarks
./scripts/validate.sh      # full validation pipeline (local mode)
./scripts/validate.sh ci   # strict CI mode
```

Your contribution must pass `./scripts/validate.sh ci` before submitting.

## 📁 Project Structure

```
to/
├── to.go                  # Implementation (single file)
├── to_test.go             # Unit tests
├── bench_test.go          # Benchmarks
├── examples/              # Runnable examples
├── scripts/
│   └── validate.sh        # Validation pipeline
├── README.md
├── CONTRIBUTING.md        # This file
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── LICENSE
└── go.mod
```

Keep it flat. No `internal/` directory.

## 🎨 Code Style

### Naming
- Outcome naming per PRINCIPLES.md. For conversions, the function name IS the target type (`Int`, `Str`, `Bool`, `F64`).

### Error Handling
- Outcome-named functions (`Int`, `Str`, `Bool`, `F64`) **must not return errors**. They return the zero value on failure. Use `Or` variants for fallbacks or `Type[T]` for full error context.
- `Type[T]` returns `(T, error)` where the error is always a `*ConversionError` with populated From/To/Value/Reason fields.
- Never panic. All inputs (including nil) are handled gracefully.

### Documentation
- Every exported function has a doc comment.
- Every numeric conversion function documents which input types it supports.
- Package-level doc comment in `to.go`.

### Dependencies
- **Zero external dependencies.** `to` is pure stdlib.

## 🧪 Testing

**Coverage target: 100% of exported functions.** Conversion functions must cover every branch of their numeric type switches.

```bash
go test -v ./...
go test -race ./...
go test -cover ./...
go test -bench=. -benchmem ./...
```

## 📝 Pull Request Process

### PR Checklist

- [ ] **Outcome naming** — does the function name describe what the caller gets?
- [ ] **One way** — does any existing function already do this?
- [ ] **Get out of the way** — can a Go dev use this from the signature alone?
- [ ] **Non-goals** — does this violate any of the library's stated non-goals?
- [ ] Tests cover 100% of new code
- [ ] Benchmarks added for new exported functions
- [ ] README updated
- [ ] CHANGELOG.md updated
- [ ] `./scripts/validate.sh ci` passes locally

### PR Scope
- One function per PR when adding new functionality
- Bug fixes can be grouped if they share a root cause

## 🆕 Adding a New Function

`to` is deliberately small. New additions must clear a high bar:

1. Read the library's non-goals in [README.md](README.md) and [PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md).
2. Prove the stdlib gap. Go's `strconv` package covers more than you might think, and on Go 1.26+ the `new(v)` builtin replaces many pointer-creation helpers.
3. Show real-world evidence from a codebase.
4. For new target types, first ask: can `Type[T]` already handle this?

## 🏷️ Versioning and Releases

- Semantic versioning
- v0.x: API may change between minor versions
- v1.0+: breaking changes require major version bump
- Every release updates CHANGELOG.md

## 🙏 Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## 📄 License

By contributing, you agree your contributions are licensed under the MIT License.
