package examples_test

import (
	"errors"
	"fmt"

	"github.com/bold-minds/to"
)

func ExampleInt() {
	fmt.Println(to.Int("42"))
	fmt.Println(to.Int(3.14))
	fmt.Println(to.Int("abc"))
	// Output:
	// 42
	// 3
	// 0
}

func ExampleIntOr() {
	fmt.Println(to.IntOr("42", 0))
	fmt.Println(to.IntOr("abc", 99))
	// Output:
	// 42
	// 99
}

func ExampleStr() {
	fmt.Println(to.Str("hello"))
	fmt.Println(to.Str(42))
	fmt.Println(to.Str(true))
	// Output:
	// hello
	// 42
	// true
}

func ExampleBool() {
	fmt.Println(to.Bool("true"))
	fmt.Println(to.Bool(1))
	fmt.Println(to.Bool("yes"))
	fmt.Println(to.Bool("maybe"))
	// Output:
	// true
	// true
	// true
	// false
}

func ExampleF64() {
	fmt.Println(to.F64("3.14"))
	fmt.Println(to.F64(42))
	// Output:
	// 3.14
	// 42
}

func ExampleType() {
	// Generic escape hatch for uncommon types
	v, err := to.Type[int64]("4096")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(v)
	// Output: 4096
}

func ExampleType_error() {
	// Structured error context
	_, err := to.Type[int]("not a number")

	var ce *to.ConversionError
	if errors.As(err, &ce) {
		fmt.Printf("from %s to %s: %s\n", ce.From, ce.To, ce.Reason)
	}
	// Output: from string to int: invalid numeric string
}

func ExampleVal() {
	name := "alice"
	fmt.Println(to.Val(&name))

	var missing *string
	fmt.Printf("%q\n", to.Val(missing))
	// Output:
	// alice
	// ""
}

func ExampleValOr() {
	var missing *string
	fmt.Println(to.ValOr(missing, "anonymous"))

	name := "alice"
	fmt.Println(to.ValOr(&name, "anonymous"))
	// Output:
	// anonymous
	// alice
}
