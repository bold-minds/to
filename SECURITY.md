# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability, please follow these steps:

### 1. **Do Not** Create a Public Issue

Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.

### 2. Report Privately

Send an email to **security@boldminds.tech** with:

- **Subject**: Security Vulnerability in bold-minds/to
- **Description**: Detailed description of the vulnerability
- **Steps to Reproduce**: Clear steps to reproduce the issue
- **Impact**: Potential impact and severity assessment
- **Suggested Fix**: If you have ideas for a fix (optional)

### 3. Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution**: Varies based on complexity, typically within 30 days

## Security Considerations

`to` is a pure-computation library with a very small attack surface:

- **No network I/O.** `to` does not make network calls.
- **No file I/O.** `to` does not read or write files.
- **No reflection on the happy path.** Type dispatch uses concrete type switches.
- **No external dependencies.** Pure Go stdlib.
- **Immutable.** `to` never modifies input values.
- **Nil-safe.** All functions handle nil inputs without panicking.

### Input validation

`to` accepts arbitrary `any` values as input. Well-formed Go values are handled safely, including numeric types, strings, booleans, and `nil`. Malformed values (e.g., interfaces over channels or functions for numeric conversion targets) return a `ConversionError` rather than panicking.

### Known limitations

- `to.Int` and `to.IntOr` parse strings via `strconv.Atoi`, which follows Go's standard parsing rules. Extreme values (exceeding int range) return a conversion error on the `Type[T]` path or the zero/fallback value on the outcome-named paths.
- `to.Str(v)` uses `fmt.Sprintf("%v", v)` for non-string inputs. This can produce unexpected output for types with custom `String()` methods that panic — such panics would propagate out of `to`. Since stdlib `fmt` itself handles this case, `to` does not add protection beyond what `fmt` provides.

## Security Updates

Security updates will be released as patch versions (e.g., 0.1.1), documented in CHANGELOG.md, and announced through GitHub releases.

## Acknowledgments

We appreciate responsible disclosure and will acknowledge security researchers who help improve the security of this project.
