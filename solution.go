package main

// #cgo CFLAGS: -I./C
// #include "adapter.h"
import "C"

var theScale Scale

//export scaleAdapterGo
func scaleAdapterGo(left C.COINS, right C.COINS) C.WEIGHRESULT {

	fromBits := func(bits C.COINS) []int {
		c := make([]int, 0, 12)
		i := 0
		bits >>= 1
		for bits > 0 {
			if bits&1 == 1 {
				c = append(c, i)
			}
			bits >>= 1
			i += 1
		}
		return c
	}

	a := fromBits(left)
	b := fromBits(right)

	switch theScale.Weigh(a, b) {
	case light:
		return C.WEIGHRESULT_LESS
	case heavy:
		return C.WEIGHRESULT_GREATER
	}

	return C.WEIGHRESULT_EQUAL
}

// Decide returns the identity of the counterfeit coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {
	theScale = scale
	result := C.decide()
	var weight Weight
	if result&1 == 1 {
		weight = light
	} else {
		weight = heavy
	}
	return int(result >> 1), weight
}
