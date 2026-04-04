package to_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/bold-minds/to"
)

// =============================================================================
// Str / StrOr
// =============================================================================

func TestStr(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{"string passthrough", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"nil", nil, ""},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := to.Str(tt.in); got != tt.want {
				t.Errorf("Str(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestStrOr(t *testing.T) {
	if got := to.StrOr("hello", "default"); got != "hello" {
		t.Errorf("got %q, want hello", got)
	}
	if got := to.StrOr(nil, "default"); got != "default" {
		t.Errorf("got %q, want default", got)
	}
	if got := to.StrOr(42, "default"); got != "42" {
		t.Errorf("got %q, want 42", got)
	}
}

// =============================================================================
// Int / IntOr
// =============================================================================

func TestInt(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want int
	}{
		{"int passthrough", 42, 42},
		{"int64", int64(42), 42},
		{"int32", int32(42), 42},
		{"float64", 42.0, 42},
		{"float64 truncation", 42.9, 42},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"string numeric", "42", 42},
		{"string negative", "-42", -42},
		{"string invalid", "abc", 0},
		{"nil", nil, 0},
		{"uint", uint(42), 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := to.Int(tt.in); got != tt.want {
				t.Errorf("Int(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestIntOr(t *testing.T) {
	if got := to.IntOr("42", 99); got != 42 {
		t.Errorf("got %d, want 42", got)
	}
	if got := to.IntOr("abc", 99); got != 99 {
		t.Errorf("got %d, want 99", got)
	}
	if got := to.IntOr(nil, 99); got != 99 {
		t.Errorf("got %d, want 99", got)
	}
}

// =============================================================================
// Bool / BoolOr
// =============================================================================

func TestBool(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"bool true", true, true},
		{"bool false", false, false},
		{"int 1", 1, true},
		{"int 0", 0, false},
		{"int 42", 42, true},
		{"string true", "true", true},
		{"string false", "false", false},
		{"string yes", "yes", true},
		{"string no", "no", false},
		{"string 1", "1", true},
		{"string 0", "0", false},
		{"string on", "on", true},
		{"string off", "off", false},
		{"string invalid", "maybe", false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := to.Bool(tt.in); got != tt.want {
				t.Errorf("Bool(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestBoolOr(t *testing.T) {
	if got := to.BoolOr("true", false); got != true {
		t.Errorf("got %v, want true", got)
	}
	if got := to.BoolOr("maybe", true); got != true {
		t.Errorf("got %v, want true (fallback)", got)
	}
	if got := to.BoolOr(nil, true); got != true {
		t.Errorf("got %v, want true (fallback)", got)
	}
}

// =============================================================================
// F64 / F64Or
// =============================================================================

func TestF64(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want float64
	}{
		{"float64 passthrough", 3.14, 3.14},
		{"float32", float32(3.14), float64(float32(3.14))},
		{"int", 42, 42.0},
		{"int64", int64(42), 42.0},
		{"string numeric", "3.14", 3.14},
		{"string int", "42", 42.0},
		{"string invalid", "abc", 0.0},
		{"bool true", true, 1.0},
		{"bool false", false, 0.0},
		{"nil", nil, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := to.F64(tt.in); got != tt.want {
				t.Errorf("F64(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestF64Or(t *testing.T) {
	if got := to.F64Or("3.14", 9.99); got != 3.14 {
		t.Errorf("got %v, want 3.14", got)
	}
	if got := to.F64Or("abc", 9.99); got != 9.99 {
		t.Errorf("got %v, want 9.99", got)
	}
}

// =============================================================================
// Type[T] / TypeOr[T]
// =============================================================================

func TestType_SuccessPaths(t *testing.T) {
	// int
	if v, err := to.Type[int]("42"); err != nil || v != 42 {
		t.Errorf("Type[int](\"42\") = (%d, %v)", v, err)
	}
	// int64
	if v, err := to.Type[int64]("42"); err != nil || v != 42 {
		t.Errorf("Type[int64](\"42\") = (%d, %v)", v, err)
	}
	// float64
	if v, err := to.Type[float64]("3.14"); err != nil || v != 3.14 {
		t.Errorf("Type[float64](\"3.14\") = (%v, %v)", v, err)
	}
	// float32
	if v, err := to.Type[float32]("3.14"); err != nil || v != float32(3.14) {
		t.Errorf("Type[float32](\"3.14\") = (%v, %v)", v, err)
	}
	// bool
	if v, err := to.Type[bool]("true"); err != nil || v != true {
		t.Errorf("Type[bool](\"true\") = (%v, %v)", v, err)
	}
	// string (passthrough for non-string goes through convertToString)
	if v, err := to.Type[string](42); err != nil || v != "42" {
		t.Errorf("Type[string](42) = (%q, %v)", v, err)
	}
}

func TestType_FailurePaths(t *testing.T) {
	// Invalid numeric string
	_, err := to.Type[int]("abc")
	if err == nil {
		t.Fatal("expected error")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.From != "string" || cerr.To != "int" {
		t.Errorf("got From=%q To=%q, want From=string To=int", cerr.From, cerr.To)
	}
	if cerr.Cause == nil {
		t.Error("expected Cause to be set for strconv errors")
	}
}

func TestType_NilInput(t *testing.T) {
	_, err := to.Type[int](nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.From != "nil" {
		t.Errorf("From=%q, want nil", cerr.From)
	}
}

func TestType_UnsupportedBoolString(t *testing.T) {
	_, err := to.Type[bool]("maybe")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTypeOr(t *testing.T) {
	if got := to.TypeOr[int]("42", 0); got != 42 {
		t.Errorf("got %d, want 42", got)
	}
	if got := to.TypeOr[int]("abc", 99); got != 99 {
		t.Errorf("got %d, want 99", got)
	}
}

// =============================================================================
// ConversionError
// =============================================================================

func TestConversionError_MessageWithoutCause(t *testing.T) {
	err := &to.ConversionError{
		From:   "string",
		To:     "int",
		Value:  "abc",
		Reason: "invalid numeric string",
	}
	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestConversionError_MessageWithCause(t *testing.T) {
	cause := errors.New("underlying")
	err := &to.ConversionError{
		From:   "string",
		To:     "int",
		Value:  "abc",
		Reason: "invalid",
		Cause:  cause,
	}
	if !errors.Is(err, cause) {
		t.Error("expected errors.Is to find the cause")
	}
}

// =============================================================================
// Val / ValOr
// =============================================================================

func TestVal(t *testing.T) {
	s := "hello"
	if got := to.Val(&s); got != "hello" {
		t.Errorf("got %q, want hello", got)
	}
	var nilPtr *string
	if got := to.Val(nilPtr); got != "" {
		t.Errorf("got %q, want empty string", got)
	}
	var nilInt *int
	if got := to.Val(nilInt); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestValOr(t *testing.T) {
	s := "hello"
	if got := to.ValOr(&s, "default"); got != "hello" {
		t.Errorf("got %q, want hello", got)
	}
	var nilPtr *string
	if got := to.ValOr(nilPtr, "default"); got != "default" {
		t.Errorf("got %q, want default", got)
	}
}

// =============================================================================
// Numeric type coverage — exercises every branch in convertToInt/Int64/Float64/Bool
// =============================================================================

func TestType_Int_AllNumericSources(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want int
	}{
		{"int", int(42), 42},
		{"int8", int8(42), 42},
		{"int16", int16(42), 42},
		{"int32", int32(42), 42},
		{"int64", int64(42), 42},
		{"uint", uint(42), 42},
		{"uint8", uint8(42), 42},
		{"uint16", uint16(42), 42},
		{"uint32", uint32(42), 42},
		{"uint64", uint64(42), 42},
		{"float32", float32(42.5), 42},
		{"float64", float64(42.5), 42},
		{"bool true", true, 1},
		{"bool false", false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := to.Type[int](tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestType_Int_UnsupportedSourceError(t *testing.T) {
	_, err := to.Type[int]([]int{1, 2, 3})
	if err == nil {
		t.Fatal("expected error for slice input")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.To != "int" {
		t.Errorf("To=%q, want int", cerr.To)
	}
}

func TestType_Int64_AllNumericSources(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want int64
	}{
		{"int", int(42), 42},
		{"int32", int32(42), 42},
		{"int64 passthrough", int64(42), 42},
		{"float64", float64(42.5), 42},
		{"string", "42", 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := to.Type[int64](tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestType_Int64_ErrorRewriting(t *testing.T) {
	// When a string fails int conversion on the int64 path, the error
	// should report "int64" as the target, not "int".
	_, err := to.Type[int64]("abc")
	if err == nil {
		t.Fatal("expected error")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.To != "int64" {
		t.Errorf("To=%q, want int64", cerr.To)
	}
}

func TestType_Int64_UnsupportedSourceError(t *testing.T) {
	_, err := to.Type[int64]([]int{1})
	if err == nil {
		t.Fatal("expected error")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.To != "int64" {
		t.Errorf("To=%q, want int64", cerr.To)
	}
}

func TestType_Float64_AllNumericSources(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want float64
	}{
		{"float64 passthrough", float64(3.14), 3.14},
		{"float32", float32(3.0), 3.0},
		{"int", int(42), 42.0},
		{"int8", int8(42), 42.0},
		{"int16", int16(42), 42.0},
		{"int32", int32(42), 42.0},
		{"int64", int64(42), 42.0},
		{"uint", uint(42), 42.0},
		{"uint8", uint8(42), 42.0},
		{"uint16", uint16(42), 42.0},
		{"uint32", uint32(42), 42.0},
		{"uint64", uint64(42), 42.0},
		{"bool true", true, 1.0},
		{"bool false", false, 0.0},
		{"string", "3.14", 3.14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := to.Type[float64](tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_Float64_UnsupportedSource(t *testing.T) {
	_, err := to.Type[float64]([]int{1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestType_Float64_BadString(t *testing.T) {
	_, err := to.Type[float64]("abc")
	if err == nil {
		t.Fatal("expected error")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.Cause == nil {
		t.Error("expected Cause to be set")
	}
}

func TestType_Float32(t *testing.T) {
	got, err := to.Type[float32]("3.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3.5 {
		t.Errorf("got %v, want 3.5", got)
	}
}

func TestType_Float32_BadString(t *testing.T) {
	_, err := to.Type[float32]("abc")
	if err == nil {
		t.Fatal("expected error for float32 bad string")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.To != "float32" {
		t.Errorf("To=%q, want float32", cerr.To)
	}
}

func TestType_Float32_OverflowFromFloat64(t *testing.T) {
	// A float64 value that exceeds the float32 range should be rejected
	// rather than silently becoming ±Inf.
	_, err := to.Type[float32](1e40)
	if err == nil {
		t.Fatal("expected overflow error for 1e40 → float32")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) || cerr.To != "float32" {
		t.Errorf("want *ConversionError with To=float32, got %T %+v", err, err)
	}

	_, err = to.Type[float32](-1e40)
	if err == nil {
		t.Fatal("expected overflow error for -1e40 → float32")
	}
}

func TestType_Float32_OverflowFromString(t *testing.T) {
	_, err := to.Type[float32]("1e40")
	if err == nil {
		t.Fatal("expected overflow error for \"1e40\" → float32")
	}
}

func TestType_Int_NaNInfRejected(t *testing.T) {
	// NaN, ±Inf, and out-of-range floats must error out rather than
	// producing INT64_MIN (the default Go amd64 behavior for int(NaN)).
	cases := map[string]float64{
		"NaN":     math.NaN(),
		"+Inf":    math.Inf(1),
		"-Inf":    math.Inf(-1),
		"too big": 1e20,
		"too sml": -1e20,
	}
	for name, v := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := to.Type[int](v)
			if err == nil {
				t.Fatalf("expected error for Type[int](%v)", v)
			}
			var cerr *to.ConversionError
			if !errors.As(err, &cerr) {
				t.Fatalf("expected *ConversionError, got %T", err)
			}
		})
	}
}

func TestInt_NaNInfReturnsZero(t *testing.T) {
	// The outcome-named shortcut swallows the error but should return 0
	// (a predictable sentinel) rather than INT64_MIN.
	if got := to.Int(math.NaN()); got != 0 {
		t.Errorf("to.Int(NaN) = %d, want 0", got)
	}
	if got := to.Int(math.Inf(1)); got != 0 {
		t.Errorf("to.Int(+Inf) = %d, want 0", got)
	}
	if got := to.Int(math.Inf(-1)); got != 0 {
		t.Errorf("to.Int(-Inf) = %d, want 0", got)
	}
	if got := to.Int(1e20); got != 0 {
		t.Errorf("to.Int(1e20) = %d, want 0", got)
	}
}

func TestType_Duration_FromString(t *testing.T) {
	// time.Duration has underlying type int64 but is a named type, so the
	// concrete type switch in Type[T] does NOT match it. The reflect-based
	// fallback must handle it, and strings like "5s" must parse via
	// time.ParseDuration.
	got, err := to.Type[time.Duration]("5s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 5*time.Second {
		t.Errorf("got %v, want 5s", got)
	}

	got, err = to.Type[time.Duration]("1h30m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 90*time.Minute {
		t.Errorf("got %v, want 1h30m", got)
	}
}

func TestType_Duration_FromInt64(t *testing.T) {
	// Numeric sources should populate the underlying int64 directly (same
	// semantics as time.Duration(5_000_000_000)).
	got, err := to.Type[time.Duration](int64(5_000_000_000))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 5*time.Second {
		t.Errorf("got %v, want 5s", got)
	}
}

func TestType_Duration_BadString(t *testing.T) {
	_, err := to.Type[time.Duration]("not a duration")
	if err == nil {
		t.Fatal("expected error for bad duration string")
	}
}

// NamedInt is a user-defined named integer type used to verify that Type[T]
// supports arbitrary named numeric types, not just time.Duration.
type NamedInt int64

func TestType_NamedInt(t *testing.T) {
	got, err := to.Type[NamedInt]("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Errorf("got %v, want 42", got)
	}

	got, err = to.Type[NamedInt](int32(42))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Errorf("got %v, want 42", got)
	}
}

// NamedInt8 verifies that the reflect-based path correctly range-checks
// narrower underlying kinds.
type NamedInt8 int8

func TestType_NamedInt8_Overflow(t *testing.T) {
	_, err := to.Type[NamedInt8](int(300))
	if err == nil {
		t.Fatal("expected overflow error: 300 does not fit in int8")
	}
}

// NamedUint verifies the uint branch of the reflect fallback rejects
// negative inputs.
type NamedUint uint32

func TestType_NamedUint_NegativeRejected(t *testing.T) {
	_, err := to.Type[NamedUint](-1)
	if err == nil {
		t.Fatal("expected error for negative → unsigned named type")
	}
}

func TestType_Any_Nil(t *testing.T) {
	// When T is interface{} (any), Type[any](nil) should return (nil, nil)
	// via the fast path — a nil any trivially satisfies any.
	got, err := to.Type[any](nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}

	// And a non-nil any should round-trip unchanged.
	got, err = to.Type[any](42)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Errorf("got %v, want 42", got)
	}
}

func TestInt_NegativeFloatTruncation(t *testing.T) {
	// Complement to the existing positive-float truncation test: -42.9
	// must truncate toward zero (→ -42), matching Go's float→int rules.
	if got := to.Int(-42.9); got != -42 {
		t.Errorf("to.Int(-42.9) = %d, want -42", got)
	}
}

func TestBool_MixedCaseString(t *testing.T) {
	// With the switch to strconv.ParseBool + strings.EqualFold, unusual
	// casings like tRuE and yEs should be accepted.
	cases := map[string]bool{
		"tRuE": true,
		"FaLsE": false,
		"yEs":  true,
		"nO":   false,
		"oN":   true,
		"oFf":  false,
	}
	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := to.Type[bool](in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestStr_AllNumericKinds(t *testing.T) {
	// Covers every branch of the strconv fast path in Str.
	cases := []struct {
		in   any
		want string
	}{
		{int8(-8), "-8"},
		{int16(-16), "-16"},
		{int32(-32), "-32"},
		{int64(-64), "-64"},
		{uint(1), "1"},
		{uint8(8), "8"},
		{uint16(16), "16"},
		{uint32(32), "32"},
		{uint64(64), "64"},
		{float32(1.5), "1.5"},
		{float64(2.5), "2.5"},
	}
	for _, tc := range cases {
		got := to.Str(tc.in)
		if got != tc.want {
			t.Errorf("Str(%v %T) = %q, want %q", tc.in, tc.in, got, tc.want)
		}
	}
	// Unknown type falls back to fmt.Sprintf("%v", …).
	type thing struct{ X int }
	if got := to.Str(thing{7}); got != "{7}" {
		t.Errorf("Str(thing{7}) = %q, want {7}", got)
	}
}

// NamedFloat exercises the reflect Float branch of convertNamed.
type NamedFloat float32

func TestType_NamedFloat(t *testing.T) {
	got, err := to.Type[NamedFloat](3.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3.5 {
		t.Errorf("got %v, want 3.5", got)
	}
}

func TestType_NamedFloat_Overflow(t *testing.T) {
	// float32 range overflow via the reflect path.
	_, err := to.Type[NamedFloat](1e40)
	if err == nil {
		t.Fatal("expected overflow error")
	}
}

func TestType_NamedFloat_BadSource(t *testing.T) {
	_, err := to.Type[NamedFloat]([]int{1})
	if err == nil {
		t.Fatal("expected error for slice source")
	}
}

func TestType_NamedUint_Overflow(t *testing.T) {
	// uint8 only holds 0-255.
	type u8 uint8
	_, err := to.Type[u8](int(300))
	if err == nil {
		t.Fatal("expected overflow error")
	}
}

func TestType_NamedUint_FromString(t *testing.T) {
	type Port uint16
	got, err := to.Type[Port]("8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 8080 {
		t.Errorf("got %v, want 8080", got)
	}
}

func TestType_NamedInt_BadSource(t *testing.T) {
	_, err := to.Type[NamedInt]([]int{1})
	if err == nil {
		t.Fatal("expected error for slice source")
	}
}

func TestType_NamedInt_NaN(t *testing.T) {
	_, err := to.Type[NamedInt](math.NaN())
	if err == nil {
		t.Fatal("expected NaN error")
	}
}

func TestType_UnsupportedTarget_Struct(t *testing.T) {
	// Struct targets fall through convertNamed's default branch to the
	// generic "unsupported conversion" error.
	type custom struct{ X int }
	_, err := to.Type[custom]("hello")
	if err == nil {
		t.Fatal("expected error")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
}

func TestToBool_UnsupportedSource(t *testing.T) {
	type thing struct{ X int }
	_, err := to.Type[bool](thing{1})
	if err == nil {
		t.Fatal("expected error for struct source")
	}
}

func TestType_Duration_NumericOverflowViaString(t *testing.T) {
	// Strings that aren't valid Go durations AND aren't integers should
	// fail cleanly.
	_, err := to.Type[time.Duration]("5x")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStrOr_NonNilStructUsesSprintf(t *testing.T) {
	// Documented behavior: StrOr returns fallback ONLY when v is nil.
	// Any non-nil value (including arbitrary structs) is formatted via
	// Str/fmt.Sprintf and the fallback is unreachable.
	type point struct{ X, Y int }
	got := to.StrOr(point{1, 2}, "fallback")
	if got == "fallback" {
		t.Errorf("fallback should be unreachable for non-nil v, got %q", got)
	}
}

func TestType_Bool_AllNumericSources(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"int 0", int(0), false},
		{"int 1", int(1), true},
		{"int8 0", int8(0), false},
		{"int16 0", int16(0), false},
		{"int32 0", int32(0), false},
		{"int64 0", int64(0), false},
		{"uint 0", uint(0), false},
		{"uint8 0", uint8(0), false},
		{"uint16 0", uint16(0), false},
		{"uint32 0", uint32(0), false},
		{"uint64 0", uint64(0), false},
		{"uint 1", uint(1), true},
		{"float32 0", float32(0), false},
		{"float32 1", float32(1), true},
		{"float64 0", float64(0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := to.Type[bool](tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_Bool_AllStringForms(t *testing.T) {
	trueForms := []string{"true", "1", "yes", "on", "TRUE", "True", "YES", "Yes", "ON", "On"}
	falseForms := []string{"false", "0", "no", "off", "FALSE", "False", "NO", "No", "OFF", "Off", ""}

	for _, s := range trueForms {
		t.Run("true:"+s, func(t *testing.T) {
			got, err := to.Type[bool](s)
			if err != nil || !got {
				t.Errorf("Type[bool](%q) = (%v, %v), want (true, nil)", s, got, err)
			}
		})
	}
	for _, s := range falseForms {
		t.Run("false:"+s, func(t *testing.T) {
			got, err := to.Type[bool](s)
			if err != nil || got {
				t.Errorf("Type[bool](%q) = (%v, %v), want (false, nil)", s, got, err)
			}
		})
	}
}

func TestType_Bool_UnsupportedSource(t *testing.T) {
	_, err := to.Type[bool]([]int{1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestType_Int_UintOverflow(t *testing.T) {
	// uint values exceeding math.MaxInt cannot fit in int safely
	huge := uint(^uint(0)) // max uint
	_, err := to.Type[int](huge)
	if err == nil {
		t.Fatal("expected overflow error for max uint")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.From != "uint" {
		t.Errorf("From=%q, want uint", cerr.From)
	}
}

func TestType_Int_Uint64Overflow(t *testing.T) {
	huge := uint64(^uint64(0)) // max uint64
	_, err := to.Type[int](huge)
	if err == nil {
		t.Fatal("expected overflow error for max uint64")
	}
	var cerr *to.ConversionError
	if !errors.As(err, &cerr) {
		t.Fatal("expected *ConversionError")
	}
	if cerr.From != "uint64" {
		t.Errorf("From=%q, want uint64", cerr.From)
	}
}

func TestType_UnsupportedTarget(t *testing.T) {
	// A target type with no conversion support (e.g., a struct) should return
	// an error via the fallthrough "unsupported conversion" path.
	type custom struct{ X int }
	_, err := to.Type[custom](42)
	if err == nil {
		t.Fatal("expected error for unsupported target type")
	}
}

func TestConversionError_ErrorMessageNoCause(t *testing.T) {
	err := &to.ConversionError{
		From:   "string",
		To:     "int",
		Value:  "abc",
		Reason: "invalid",
	}
	msg := err.Error()
	if msg == "" {
		t.Fatal("expected non-empty message")
	}
	// Should not contain ": <nil>" style cause formatting
	if msg != "cannot convert string(abc) to int: invalid" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestConversionError_ErrorMessageWithCause(t *testing.T) {
	err := &to.ConversionError{
		From:   "string",
		To:     "int",
		Value:  "abc",
		Reason: "invalid",
		Cause:  errors.New("underlying"),
	}
	msg := err.Error()
	if msg != "cannot convert string(abc) to int: invalid: underlying" {
		t.Errorf("unexpected message: %q", msg)
	}
}

// =============================================================================
// Integration tests — realistic usage patterns
// =============================================================================

func TestIntegration_EnvVarParsing(t *testing.T) {
	// Simulating os.Getenv which returns strings
	portStr := "8080"
	port := to.IntOr(portStr, 3000)
	if port != 8080 {
		t.Errorf("got %d, want 8080", port)
	}

	missing := ""
	port = to.IntOr(missing, 3000)
	if port != 3000 {
		t.Errorf("got %d, want 3000 (fallback)", port)
	}
}

func TestIntegration_ConfigMap(t *testing.T) {
	// Simulating a config parsed from JSON/YAML where values are any
	cfg := map[string]any{
		"timeout": 30.0,    // JSON numbers unmarshal as float64
		"host":    "localhost",
		"debug":   true,
	}
	timeout := to.IntOr(cfg["timeout"], 60)
	if timeout != 30 {
		t.Errorf("got %d, want 30", timeout)
	}
	host := to.StrOr(cfg["host"], "")
	if host != "localhost" {
		t.Errorf("got %q, want localhost", host)
	}
	debug := to.BoolOr(cfg["debug"], false)
	if !debug {
		t.Errorf("got %v, want true", debug)
	}
}
