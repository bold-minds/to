package to_test

import (
	"errors"
	"testing"

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
