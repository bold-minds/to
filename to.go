// Package to provides safe Go value conversion with ergonomic fallbacks.
//
// It offers outcome-named shortcuts for common target types (Str, Int, Bool,
// F64) along with matching fallback variants (StrOr, IntOr, etc.), a generic
// escape hatch (Type[T], TypeOr[T]) with typed error context, and safe
// pointer creation and dereferencing (Ptr, Val, ValOr).
//
// Every function is nil-safe and never panics. Float-to-integer conversions
// reject NaN, ±Inf, and out-of-range values rather than producing trap
// values. Float-to-float narrowing rejects values that would overflow to
// ±Inf.
//
// # Fallback semantics
//
// IntOr, BoolOr, F64Or, and TypeOr all return the supplied fallback when
// conversion fails — including the empty string, which is explicitly
// treated as "not a valid boolean/number" rather than as a valid zero
// value. This makes the package well-suited for env-var and JSON-config
// parsing, where an empty string overwhelmingly means "unset".
//
// StrOr is the exception: it returns the fallback only when v is nil.
// For any non-nil input (including arbitrary structs) StrOr returns the
// same result as Str — fmt/strconv-based formatting never fails, so the
// fallback is unreachable. If you want "only accept real strings"
// semantics, type-assert directly or call Type[string].
//
// For documentation and examples, see https://github.com/bold-minds/to.
package to

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ConversionError is returned by Type[T] when conversion fails. It carries
// structured context about the conversion attempt for debugging.
type ConversionError struct {
	From   string // Source type name (from fmt.Sprintf("%T", value))
	To     string // Target type name
	Value  any    // The original value
	Reason string // Human-readable explanation
	Cause  error  // Underlying error, if any (e.g., from strconv)
}

// Error implements the error interface.
func (e *ConversionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("cannot convert %s(%v) to %s: %s: %v",
			e.From, e.Value, e.To, e.Reason, e.Cause)
	}
	return fmt.Sprintf("cannot convert %s(%v) to %s: %s",
		e.From, e.Value, e.To, e.Reason)
}

// Unwrap returns the underlying cause for errors.Is and errors.As.
func (e *ConversionError) Unwrap() error {
	return e.Cause
}

// =============================================================================
// Outcome-named conversions — return zero value on failure
// =============================================================================

// Str converts v to a string. Returns the empty string for nil input.
// Common numeric and boolean types are formatted via strconv (allocation-free
// where possible); []byte is converted directly; error and fmt.Stringer get
// dedicated fast paths to avoid the reflection-heavy fmt.Sprintf fallback.
// Any other type falls back to fmt.Sprintf("%v", v).
func Str(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case bool:
		return strconv.FormatBool(x)
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	}
	// Interface-typed fast paths, checked after concrete cases so that
	// common numeric/bool types never reach here. error is checked before
	// fmt.Stringer because some types implement both (e.g., *url.Error).
	if err, ok := v.(error); ok {
		return err.Error()
	}
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", v)
}

// Int converts v to an int. Returns 0 on failure.
// Supports all numeric types, bool (true→1, false→0), and numeric strings.
// Float-to-int conversion truncates toward zero; NaN, ±Inf, and values
// outside the int range produce 0 (via a silenced error).
func Int(v any) int {
	result, _ := Type[int](v)
	return result
}

// Bool converts v to a bool. Returns false on failure.
// Supports bool, numeric types (non-zero→true), and strings. String
// parsing is case-insensitive for the alphabetic forms and accepts
// true/false, 1/0, yes/no, on/off. No whitespace trimming is performed;
// callers should trim inputs first. The empty string is not a valid
// boolean and causes failure (so BoolOr's fallback fires).
func Bool(v any) bool {
	result, _ := Type[bool](v)
	return result
}

// F64 converts v to a float64. Returns 0.0 on failure.
// Supports all numeric types, bool (true→1.0, false→0.0), and numeric strings.
func F64(v any) float64 {
	result, _ := Type[float64](v)
	return result
}

// =============================================================================
// Outcome-named conversions with fallback
// =============================================================================

// StrOr returns fallback only when v is nil. For any non-nil input StrOr
// returns the same result as Str — fmt/strconv-based formatting never fails,
// so the fallback is unreachable for non-nil values. If you need stricter
// "only accept real strings" semantics, type-assert directly or use
// Type[string].
func StrOr(v any, fallback string) string {
	if v == nil {
		return fallback
	}
	return Str(v)
}

// IntOr converts v to an int, returning fallback on failure.
func IntOr(v any, fallback int) int {
	return TypeOr(v, fallback)
}

// BoolOr converts v to a bool, returning fallback on failure.
func BoolOr(v any, fallback bool) bool {
	return TypeOr(v, fallback)
}

// F64Or converts v to a float64, returning fallback on failure.
func F64Or(v any, fallback float64) float64 {
	return TypeOr(v, fallback)
}

// =============================================================================
// Generic conversion — returns (value, error) for full context
// =============================================================================

// Type converts v to type T. It returns (value, nil) on success or
// (zero, *ConversionError) on failure.
//
// Use Type[T] for uncommon target types (time.Duration, int64, float32,
// user-defined named numeric types) or when you need structured error
// context. For common targets (string, int, bool, float64), the
// outcome-named shortcuts are shorter.
//
// Named numeric types (e.g., time.Duration, type Port int64) are supported
// via reflection and range-checked against their underlying kind.
// time.Duration additionally accepts Go duration strings like "5s" or
// "1h30m" via time.ParseDuration.
//
// When T is an interface type (including any), the fast path returns v
// unchanged if v already satisfies T. If v is nil and T is any, the zero
// value (nil) is returned with no error; for non-empty interface targets,
// nil input produces a ConversionError.
func Type[T any](v any) (T, error) {
	var zero T

	// Fast path: already the target type (also handles T=any for any v,
	// including nil).
	if result, ok := v.(T); ok {
		return result, nil
	}

	// Handle nil input. When T is an empty interface (any), nil is a valid
	// zero value and we return it without error; the fast-path v.(T) above
	// does not match because a nil untyped interface cannot be asserted to
	// another interface type.
	if v == nil {
		rt := reflect.TypeOf((*T)(nil)).Elem()
		if rt.Kind() == reflect.Interface && rt.NumMethod() == 0 {
			return zero, nil
		}
		return zero, &ConversionError{
			From:   "nil",
			To:     fmt.Sprintf("%T", zero),
			Value:  nil,
			Reason: "cannot convert nil",
		}
	}

	// Dispatch on target type. In each case the inner type assertion
	// (any(x).(T)) is algebraically guaranteed to succeed because T is
	// pinned by the matching case — so we discard the ok value.
	switch any(zero).(type) {
	case string:
		r, _ := any(Str(v)).(T)
		return r, nil
	case int:
		n, err := toInt64(v, "int")
		if err != nil {
			return zero, err
		}
		if n > math.MaxInt || n < math.MinInt {
			return zero, &ConversionError{
				From:   fmt.Sprintf("%T", v),
				To:     "int",
				Value:  v,
				Reason: "value exceeds int range",
			}
		}
		r, _ := any(int(n)).(T)
		return r, nil
	case int64:
		n, err := toInt64(v, "int64")
		if err != nil {
			return zero, err
		}
		r, _ := any(n).(T)
		return r, nil
	case float64:
		n, err := toFloat64(v, "float64")
		if err != nil {
			return zero, err
		}
		r, _ := any(n).(T)
		return r, nil
	case float32:
		n, err := toFloat64(v, "float32")
		if err != nil {
			return zero, err
		}
		if n > math.MaxFloat32 || n < -math.MaxFloat32 {
			return zero, &ConversionError{
				From:   fmt.Sprintf("%T", v),
				To:     "float32",
				Value:  v,
				Reason: "value exceeds float32 range",
			}
		}
		r, _ := any(float32(n)).(T)
		return r, nil
	case bool:
		b, err := toBool(v, "bool")
		if err != nil {
			return zero, err
		}
		r, _ := any(b).(T)
		return r, nil
	}

	// Fallback: named numeric types (time.Duration, type Port int64, …).
	if result, ok, err := convertNamed[T](v); ok {
		return result, err
	}

	return zero, &ConversionError{
		From:   fmt.Sprintf("%T", v),
		To:     fmt.Sprintf("%T", zero),
		Value:  v,
		Reason: "unsupported conversion",
	}
}

// TypeOr converts v to type T, returning fallback on any failure.
func TypeOr[T any](v any, fallback T) T {
	if result, err := Type[T](v); err == nil {
		return result
	}
	return fallback
}

// =============================================================================
// Pointer creation and dereferencing
// =============================================================================

// Ptr returns a pointer to a fresh copy of v. It is the ergonomic
// one-call alternative to the two-line temp-variable pattern:
//
//	x := v
//	p := &x
//
// Equivalent to Go 1.26+'s new(v) builtin, but works on all supported
// Go versions of this package.
func Ptr[T any](v T) *T {
	return &v
}

// Val returns the value pointed to by p, or the zero value of T if p is nil.
// Val returns a copy; for large struct types, prefer reading fields directly
// to avoid the copy cost.
func Val[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// ValOr returns the value pointed to by p, or fallback if p is nil.
// ValOr returns a copy; for large struct types, prefer reading fields directly
// to avoid the copy cost.
func ValOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// =============================================================================
// Internal conversion helpers
// =============================================================================

// convertNamed handles target types whose underlying kind is numeric but
// which aren't matched by Type[T]'s concrete type switch — e.g., time.Duration
// or user-defined `type Port int64`. Returns (value, true, nil) on success,
// (_, true, err) if the kind was supported but conversion failed, and
// (_, false, _) if the kind is not numeric.
func convertNamed[T any](v any) (T, bool, error) {
	var zero T
	rt := reflect.TypeOf(zero)
	if rt == nil {
		return zero, false, nil
	}
	targetName := rt.String()
	fromName := fmt.Sprintf("%T", v)

	// Special case: time.Duration can parse Go duration strings like "5s".
	if rt == reflect.TypeOf(time.Duration(0)) {
		if s, isString := v.(string); isString {
			if d, perr := time.ParseDuration(s); perr == nil {
				if r, assigned := any(d).(T); assigned {
					return r, true, nil
				}
			}
			// fall through: allow pure numeric duration strings via toInt64
		}
	}

	rv := reflect.New(rt).Elem()
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := toInt64(v, targetName)
		if err != nil {
			return zero, true, err
		}
		if rv.OverflowInt(n) {
			return zero, true, numericOverflowError(fromName, targetName, v)
		}
		rv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := toUint64(v, targetName)
		if err != nil {
			return zero, true, err
		}
		if rv.OverflowUint(n) {
			return zero, true, numericOverflowError(fromName, targetName, v)
		}
		rv.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := toFloat64(v, targetName)
		if err != nil {
			return zero, true, err
		}
		if rv.OverflowFloat(n) {
			return zero, true, numericOverflowError(fromName, targetName, v)
		}
		rv.SetFloat(n)
	default:
		return zero, false, nil
	}

	if r, ok := rv.Interface().(T); ok {
		return r, true, nil
	}
	// Kind was numeric but the final assertion to T failed — surface a
	// generic error rather than falling back to "unsupported conversion".
	return zero, true, &ConversionError{
		From: fromName, To: targetName, Value: v,
		Reason: "internal type mismatch",
	}
}

func numericOverflowError(fromName, targetName string, v any) error {
	return &ConversionError{
		From: fromName, To: targetName, Value: v,
		Reason: "value exceeds " + targetName + " range",
	}
}

// toInt64 converts a numeric-ish value to int64. The target parameter is
// embedded in error messages so callers get accurate target naming
// ("int", "int64", "time.Duration", …) without the caller mutating the
// returned error.
func toInt64(v any, target string) (int64, error) {
	switch x := v.(type) {
	case int:
		return int64(x), nil
	case int8:
		return int64(x), nil
	case int16:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case int64:
		return x, nil
	case uint:
		if uint64(x) > math.MaxInt64 {
			return 0, &ConversionError{
				From: "uint", To: target, Value: v,
				Reason: "value exceeds int64 range",
			}
		}
		return int64(x), nil
	case uint8:
		return int64(x), nil
	case uint16:
		return int64(x), nil
	case uint32:
		return int64(x), nil
	case uint64:
		if x > math.MaxInt64 {
			return 0, &ConversionError{
				From: "uint64", To: target, Value: v,
				Reason: "value exceeds int64 range",
			}
		}
		return int64(x), nil
	case float32:
		return floatToInt64(float64(x), "float32", target, v)
	case float64:
		return floatToInt64(x, "float64", target, v)
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	case string:
		n, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return 0, &ConversionError{
				From: "string", To: target, Value: v,
				Reason: "invalid numeric string", Cause: err,
			}
		}
		return n, nil
	}
	return 0, &ConversionError{
		From: fmt.Sprintf("%T", v), To: target, Value: v,
		Reason: "unsupported source type",
	}
}

// maxInt64AsFloat is the smallest float64 that is out of int64 range.
//
// math.MaxInt64 (= 2^63 - 1) cannot be represented exactly as a float64 —
// the nearest float64 rounds up to 2^63. If we compared x > math.MaxInt64,
// the untyped constant would be promoted to float64 as 2^63, letting the
// exact value x = 2^63 slip through to int64(x), whose Go-spec behavior
// for an out-of-range float is implementation-defined (on amd64 it yields
// INT64_MIN). So we reject anything >= 2^63.
const maxInt64AsFloat = float64(1 << 63)

// floatToInt64 truncates x toward zero after rejecting NaN, ±Inf, and values
// outside the int64 range. A library whose tagline is "safe" must not return
// INT64_MIN for NaN, which is what a plain int64(x) conversion does on amd64.
func floatToInt64(x float64, from, target string, orig any) (int64, error) {
	if math.IsNaN(x) {
		return 0, &ConversionError{
			From: from, To: target, Value: orig,
			Reason: "cannot convert NaN to integer",
		}
	}
	if math.IsInf(x, 0) {
		return 0, &ConversionError{
			From: from, To: target, Value: orig,
			Reason: "cannot convert Inf to integer",
		}
	}
	// Note: x < math.MinInt64 is fine as-is because -2^63 is exactly
	// representable in float64, so the strict comparison correctly
	// accepts the boundary. The upper bound must be >= maxInt64AsFloat
	// because 2^63 - 1 is not representable; see maxInt64AsFloat.
	if x >= maxInt64AsFloat || x < math.MinInt64 {
		return 0, &ConversionError{
			From: from, To: target, Value: orig,
			Reason: "value exceeds int64 range",
		}
	}
	return int64(x), nil
}

// toUint64 converts a numeric-ish value to uint64, mirroring toInt64 but
// preserving the full uint64 range for unsigned named targets (e.g.,
// type Size uint64). Negative sources are rejected.
func toUint64(v any, target string) (uint64, error) {
	switch x := v.(type) {
	case uint:
		return uint64(x), nil
	case uint8:
		return uint64(x), nil
	case uint16:
		return uint64(x), nil
	case uint32:
		return uint64(x), nil
	case uint64:
		return x, nil
	case uintptr:
		return uint64(x), nil
	case int:
		if x < 0 {
			return 0, negativeToUnsignedError(v, target)
		}
		return uint64(x), nil
	case int8:
		if x < 0 {
			return 0, negativeToUnsignedError(v, target)
		}
		return uint64(x), nil
	case int16:
		if x < 0 {
			return 0, negativeToUnsignedError(v, target)
		}
		return uint64(x), nil
	case int32:
		if x < 0 {
			return 0, negativeToUnsignedError(v, target)
		}
		return uint64(x), nil
	case int64:
		if x < 0 {
			return 0, negativeToUnsignedError(v, target)
		}
		return uint64(x), nil
	case float32:
		return floatToUint64(float64(x), "float32", target, v)
	case float64:
		return floatToUint64(x, "float64", target, v)
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	case string:
		n, err := strconv.ParseUint(x, 10, 64)
		if err != nil {
			return 0, &ConversionError{
				From: "string", To: target, Value: v,
				Reason: "invalid numeric string", Cause: err,
			}
		}
		return n, nil
	}
	return 0, &ConversionError{
		From: fmt.Sprintf("%T", v), To: target, Value: v,
		Reason: "unsupported source type",
	}
}

func negativeToUnsignedError(v any, target string) error {
	return &ConversionError{
		From: fmt.Sprintf("%T", v), To: target, Value: v,
		Reason: "negative value cannot be converted to unsigned type",
	}
}

// maxUint64AsFloat is the smallest float64 strictly greater than math.MaxUint64.
// math.MaxUint64 = 2^64 - 1 is not exactly representable as float64; the
// nearest float64 rounds to 2^64. Any x >= 2^64 is out of uint64 range.
const maxUint64AsFloat = float64(1 << 64)

func floatToUint64(x float64, from, target string, orig any) (uint64, error) {
	if math.IsNaN(x) {
		return 0, &ConversionError{
			From: from, To: target, Value: orig,
			Reason: "cannot convert NaN to integer",
		}
	}
	if math.IsInf(x, 0) {
		return 0, &ConversionError{
			From: from, To: target, Value: orig,
			Reason: "cannot convert Inf to integer",
		}
	}
	if x < 0 {
		return 0, negativeToUnsignedError(orig, target)
	}
	if x >= maxUint64AsFloat {
		return 0, &ConversionError{
			From: from, To: target, Value: orig,
			Reason: "value exceeds uint64 range",
		}
	}
	return uint64(x), nil
}

func toFloat64(v any, target string) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int8:
		return float64(x), nil
	case int16:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case uint8:
		return float64(x), nil
	case uint16:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case bool:
		if x {
			return 1.0, nil
		}
		return 0.0, nil
	case string:
		n, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, &ConversionError{
				From: "string", To: target, Value: v,
				Reason: "invalid numeric string", Cause: err,
			}
		}
		return n, nil
	}
	return 0, &ConversionError{
		From: fmt.Sprintf("%T", v), To: target, Value: v,
		Reason: "unsupported source type",
	}
}

func toBool(v any, target string) (bool, error) {
	switch x := v.(type) {
	case bool:
		return x, nil
	case int:
		return x != 0, nil
	case int8:
		return x != 0, nil
	case int16:
		return x != 0, nil
	case int32:
		return x != 0, nil
	case int64:
		return x != 0, nil
	case uint:
		return x != 0, nil
	case uint8:
		return x != 0, nil
	case uint16:
		return x != 0, nil
	case uint32:
		return x != 0, nil
	case uint64:
		return x != 0, nil
	case float32:
		return x != 0, nil
	case float64:
		return x != 0, nil
	case string:
		return parseBoolString(x, v, target)
	}
	return false, &ConversionError{
		From: fmt.Sprintf("%T", v), To: target, Value: v,
		Reason: "unsupported source type",
	}
}

// parseBoolString accepts true/false, 1/0, yes/no, on/off (all
// case-insensitive for the alphabetic forms). No whitespace trimming is
// performed; callers should trim first. The empty string is NOT treated
// as a valid boolean: it returns a ConversionError so that BoolOr's
// fallback fires for unset env vars and missing config values.
func parseBoolString(x string, orig any, target string) (bool, error) {
	if x == "" {
		return false, &ConversionError{
			From: "string", To: target, Value: orig,
			Reason: "empty string is not a valid boolean",
		}
	}
	switch x {
	case "1":
		return true, nil
	case "0":
		return false, nil
	}
	if strings.EqualFold(x, "true") || strings.EqualFold(x, "yes") || strings.EqualFold(x, "on") {
		return true, nil
	}
	if strings.EqualFold(x, "false") || strings.EqualFold(x, "no") || strings.EqualFold(x, "off") {
		return false, nil
	}
	return false, &ConversionError{
		From: "string", To: target, Value: orig,
		Reason: "unrecognized boolean string (expected true/false/1/0/yes/no/on/off)",
	}
}

// Compile-time assertions that *ConversionError satisfies error and exposes
// Unwrap for errors.Is/errors.As.
var (
	_ error                       = (*ConversionError)(nil)
	_ interface{ Unwrap() error } = (*ConversionError)(nil)
)
