package to_test

import (
	"testing"

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
