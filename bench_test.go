package to_test

import (
	"testing"
	"time"

	"github.com/bold-minds/to"
)

func BenchmarkStr_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Str("hello")
	}
}

func BenchmarkStr_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Str(42)
	}
}

func BenchmarkInt_FromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Int("42")
	}
}

func BenchmarkInt_FromFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Int(42.0)
	}
}

func BenchmarkInt_FromInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Int(42)
	}
}

func BenchmarkBool_FromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Bool("true")
	}
}

func BenchmarkBool_FromInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Bool(1)
	}
}

func BenchmarkF64_FromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.F64("3.14")
	}
}

func BenchmarkF64_FromInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.F64(42)
	}
}

func BenchmarkIntOr_Success(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.IntOr("42", 0)
	}
}

func BenchmarkIntOr_Fallback(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.IntOr("abc", 99)
	}
}

func BenchmarkType_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = to.Type[int]("42")
	}
}

func BenchmarkPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = to.Ptr(42)
	}
}

func BenchmarkVal(b *testing.B) {
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = to.Val(&s)
	}
}

func BenchmarkValOr_NonNil(b *testing.B) {
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = to.ValOr(&s, "default")
	}
}

func BenchmarkValOr_Nil(b *testing.B) {
	var p *string
	for i := 0; i < b.N; i++ {
		_ = to.ValOr(p, "default")
	}
}

// =============================================================================
// Reflective slow path — convertNamed
//
// Type[int], Type[bool], etc. go through the fast concrete-type switch.
// Named numeric targets like time.Duration and user-defined int/uint types
// fall through to the reflection-based convertNamed helper, which is the
// reason for most of the allocation in this package. Benchmark it
// explicitly so regressions don't hide behind the fast-path numbers.
// =============================================================================

type benchNamedInt int64
type benchNamedUint uint32

func BenchmarkType_NamedInt_FromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = to.Type[benchNamedInt]("42")
	}
}

func BenchmarkType_NamedInt_FromInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = to.Type[benchNamedInt](int(42))
	}
}

func BenchmarkType_NamedUint_FromInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = to.Type[benchNamedUint](int(42))
	}
}

func BenchmarkType_Duration_FromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = to.Type[time.Duration]("5s")
	}
}

func BenchmarkType_Duration_FromInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = to.Type[time.Duration](int64(5_000_000_000))
	}
}
