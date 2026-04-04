// Package to provides safe Go value conversion with ergonomic fallbacks.
//
// It offers outcome-named shortcuts for common target types (Str, Int, Bool,
// F64) along with matching fallback variants (StrOr, IntOr, etc.), a generic
// escape hatch (Type[T], TypeOr[T]) with typed error context, and safe
// pointer dereferencing (Val, ValOr).
//
// Every function is nil-safe and never panics. For pointer creation, use
// Go 1.26's new(v) builtin directly — this package does not provide Ptr(v).
//
// For documentation and examples, see https://github.com/bold-minds/to.
package to

import (
	"fmt"
	"math"
	"strconv"
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

// Str converts v to a string. Returns the empty string on failure.
// For non-string inputs, Str uses fmt.Sprintf("%v", v).
func Str(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// Int converts v to an int. Returns 0 on failure.
// Supports all numeric types, bool (true→1, false→0), and numeric strings.
func Int(v any) int {
	result, _ := Type[int](v)
	return result
}

// Bool converts v to a bool. Returns false on failure.
// Supports bool, numeric types (non-zero→true), and strings
// ("true"/"1"/"yes"/"on" → true, "false"/"0"/"no"/"off" → false).
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

// StrOr converts v to a string, returning fallback on failure.
func StrOr(v any, fallback string) string {
	if v == nil {
		return fallback
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
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
// etc.) or when you need structured error context. For common targets
// (string, int, bool, float64), the outcome-named shortcuts are shorter.
func Type[T any](v any) (T, error) {
	var zero T

	// Fast path: already the target type
	if result, ok := v.(T); ok {
		return result, nil
	}

	// Handle nil input
	if v == nil {
		return zero, &ConversionError{
			From:   "nil",
			To:     fmt.Sprintf("%T", zero),
			Value:  nil,
			Reason: "cannot convert nil",
		}
	}

	// Dispatch on target type
	switch any(zero).(type) {
	case string:
		if result, ok := any(convertToString(v)).(T); ok {
			return result, nil
		}
	case int:
		if n, err := convertToInt(v); err == nil {
			if result, ok := any(n).(T); ok {
				return result, nil
			}
		} else {
			return zero, err
		}
	case int64:
		if n, err := convertToInt64(v); err == nil {
			if result, ok := any(n).(T); ok {
				return result, nil
			}
		} else {
			return zero, err
		}
	case float64:
		if n, err := convertToFloat64(v); err == nil {
			if result, ok := any(n).(T); ok {
				return result, nil
			}
		} else {
			return zero, err
		}
	case float32:
		if n, err := convertToFloat64(v); err == nil {
			if result, ok := any(float32(n)).(T); ok {
				return result, nil
			}
		} else {
			return zero, err
		}
	case bool:
		if b, err := convertToBool(v); err == nil {
			if result, ok := any(b).(T); ok {
				return result, nil
			}
		} else {
			return zero, err
		}
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
// Pointer dereferencing
// =============================================================================

// Val returns the value pointed to by p, or the zero value of T if p is nil.
// For creating pointers, use Go 1.26's new(v) builtin directly.
func Val[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// ValOr returns the value pointed to by p, or fallback if p is nil.
func ValOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// =============================================================================
// Internal conversion helpers
// =============================================================================

func convertToString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func convertToInt(v any) (int, error) {
	switch x := v.(type) {
	case int:
		return x, nil
	case int8:
		return int(x), nil
	case int16:
		return int(x), nil
	case int32:
		return int(x), nil
	case int64:
		return int(x), nil
	case uint:
		if x > math.MaxInt {
			return 0, &ConversionError{
				From:   "uint",
				To:     "int",
				Value:  v,
				Reason: "value exceeds math.MaxInt",
			}
		}
		return int(x), nil
	case uint8:
		return int(x), nil
	case uint16:
		return int(x), nil
	case uint32:
		return int(x), nil
	case uint64:
		if x > math.MaxInt {
			return 0, &ConversionError{
				From:   "uint64",
				To:     "int",
				Value:  v,
				Reason: "value exceeds math.MaxInt",
			}
		}
		return int(x), nil
	case float32:
		return int(x), nil
	case float64:
		return int(x), nil
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	case string:
		n, err := strconv.Atoi(x)
		if err != nil {
			return 0, &ConversionError{
				From:   "string",
				To:     "int",
				Value:  v,
				Reason: "invalid numeric string",
				Cause:  err,
			}
		}
		return n, nil
	}
	return 0, &ConversionError{
		From:   fmt.Sprintf("%T", v),
		To:     "int",
		Value:  v,
		Reason: "unsupported source type",
	}
}

func convertToInt64(v any) (int64, error) {
	n, err := convertToInt(v)
	if err != nil {
		// Rewrite the target in the error for accuracy
		if cerr, ok := err.(*ConversionError); ok {
			cerr.To = "int64"
		}
		return 0, err
	}
	return int64(n), nil
}

func convertToFloat64(v any) (float64, error) {
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
				From:   "string",
				To:     "float64",
				Value:  v,
				Reason: "invalid numeric string",
				Cause:  err,
			}
		}
		return n, nil
	}
	return 0, &ConversionError{
		From:   fmt.Sprintf("%T", v),
		To:     "float64",
		Value:  v,
		Reason: "unsupported source type",
	}
}

func convertToBool(v any) (bool, error) {
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
		switch x {
		case "true", "1", "yes", "on", "TRUE", "True", "YES", "Yes", "ON", "On":
			return true, nil
		case "false", "0", "no", "off", "FALSE", "False", "NO", "No", "OFF", "Off", "":
			return false, nil
		}
		return false, &ConversionError{
			From:   "string",
			To:     "bool",
			Value:  v,
			Reason: "unrecognized boolean string (expected true/false/1/0/yes/no/on/off)",
		}
	}
	return false, &ConversionError{
		From:   fmt.Sprintf("%T", v),
		To:     "bool",
		Value:  v,
		Reason: "unsupported source type",
	}
}
