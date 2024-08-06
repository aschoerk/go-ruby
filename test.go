package main

import (
	"fmt"

	"aschoerk.de/go-ruby/ruby"
)

func main() {

	fmt.Printf("Bytes: %t\n", ruby.CompareGenerally(byte(1), byte(2))) // true
	fmt.Printf("Dimensions of byte(1): %d\n", ruby.CountDimensions(byte(1)))

	fmt.Printf("Strings: %t\n", ruby.CompareGenerally("a", "b")) // true
	fmt.Printf("Dimensions of \"a\": %d\n", ruby.CountDimensions("a"))
	fmt.Printf("Integers: %t\n", ruby.CompareGenerally(1, 2)) // true
	fmt.Printf("Dimensions of 1: %d\n", ruby.CountDimensions(1))

	// 2D slice of integers
	a := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	b := [][]int{
		{1, 2, 3},
		{4, 5},
		{7, 8, 9},
	}
	fmt.Println("Comparing 2D slices:")
	fmt.Println(ruby.CompareGenerally(a, b)) // true
	fmt.Printf("Dimensions of a: %d\n", ruby.CountDimensions(a))

	// 3D slice of integer pointers
	c := [][][]*int{
		{{intPtr(1), intPtr(2)}, {intPtr(3), intPtr(4)}},
		{{intPtr(5), intPtr(6)}, {intPtr(7), nil}},
	}
	d := [][][]*int{
		{{intPtr(1), intPtr(2)}, {intPtr(3), intPtr(4)}},
		{{intPtr(5), intPtr(6)}, {intPtr(7), intPtr(8)}},
	}
	fmt.Println("\nComparing 3D slices with pointers:")
	fmt.Println(ruby.CompareGenerally(c, d)) // true (because nil < 8)
	fmt.Printf("Dimensions of c: %d\n", ruby.CountDimensions(c))

	// 2D slice of string pointers
	e := [][]*string{
		{strPtr("a"), strPtr("b"), nil},
		{strPtr("d"), strPtr("e"), strPtr("f")},
	}
	f := [][]*string{
		{strPtr("a"), strPtr("b"), strPtr("c")},
		{strPtr("d"), strPtr("e"), strPtr("f")},
	}
	fmt.Println("\nComparing 2D string slices with pointers:")
	fmt.Println(ruby.CompareGenerally(e, f)) // true (because nil < "c")
	fmt.Printf("Dimensions of e: %d\n", ruby.CountDimensions(e))
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
