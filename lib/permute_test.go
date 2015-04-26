package lib

import (
	"testing"
)

func TestNumberPermute0(t *testing.T) {
	n := Number([]int{0, 1, 2})
	if n != 0 {
		t.Fatalf("assertion failed: was: %d expected: %d", n, 0)
	}
}

func TestNumberPermute1(t *testing.T) {
	n := Number([]int{0, 2, 1})
	if n != 1 {
		t.Fatalf("assertion failed: was: %d expected: %d", n, 1)
	}
}

func TestNumberPermute2(t *testing.T) {
	n := Number([]int{1, 0, 2})
	if n != 2 {
		t.Fatalf("assertion failed: was: %d expected: %d", n, 2)
	}
}

func TestNumberPermute3(t *testing.T) {
	n := Number([]int{1, 2, 0})
	if n != 3 {
		t.Fatalf("assertion failed: was: %d expected: %d", n, 3)
	}
}

func TestNumberPermute4(t *testing.T) {
	n := Number([]int{2, 0, 1})
	if n != 4 {
		t.Fatalf("assertion failed: was: %d expected: %d", n, 4)
	}
}

func TestNumberPermute5(t *testing.T) {
	n := Number([]int{2, 1, 0})
	if n != 5 {
		t.Fatalf("assertion failed: was: %d expected: %d", n, 5)
	}
}
