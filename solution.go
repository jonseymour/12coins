package main

import (
	"fmt"
)

// Decide returns the identity of the counterfeit coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {
	switch scale.Weigh([]int{0, 1, 2}, []int{9, 10, 11}) {
	case light:
		switch scale.Weigh([]int{0, 2}, []int{1, 4}) {
		case light:
			switch scale.Weigh([]int{0}, []int{4}) {
			case light:
				return 0, light
			case equal:
				return 2, light
			case heavy:
				panic(fmt.Errorf("incomplete case"))
			}
		case heavy:
			panic(fmt.Errorf("incomplete case"))
		case equal:
			panic(fmt.Errorf("incomplete case"))
		}
	case heavy:
		switch scale.Weigh([]int{0, 2}, []int{1, 4}) {
		case light:
			switch scale.Weigh([]int{0}, []int{4}) {
			case light:
				return 0, light
			case equal:
				return 1, heavy
			case heavy:
				panic(fmt.Errorf("incomplete case"))
			}
		case equal:
			panic(fmt.Errorf("incomplete case"))
		case heavy:
			panic(fmt.Errorf("incomplete case"))
		}
	case equal:
		panic(fmt.Errorf("incomplete case"))
	}
	return -1, equal
}
