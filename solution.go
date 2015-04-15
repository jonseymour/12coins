package main

import (
	"fmt"
)

// decide returns the identity of the different coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {
	switch scale.Weigh([]int{0, 1, 2, 3}, []int{4, 5, 6, 7}) {
	case light:
		switch scale.Weigh([]int{0, 5, 6, 8}, []int{4, 1, 9, 10}) {
		case equal:
			// 2, 3, 7
			switch scale.Weigh([]int{2, 7}, []int{8, 9}) {
			case equal:
				return 3, light
			case light:
				return 2, light
			case heavy:
				return 7, heavy
			}
		case light:
			// # 0, 4
			switch scale.Weigh([]int{0}, []int{8}) {
			case equal:
				return 4, heavy
			case light:
				return 0, light
			case heavy:
				panic(fmt.Errorf("cannot happen"))
			}
		case heavy:
			// 5, 6, 1
			switch scale.Weigh([]int{1, 5}, []int{8, 9}) {
			case equal:
				return 6, heavy
			case light:
				return 1, light
			case heavy:
				return 5, heavy
			}
		}
	case heavy:
		switch scale.Weigh([]int{0, 5, 6, 8}, []int{4, 1, 9, 10}) {
		case equal:
			// 2, 3, 7
			switch scale.Weigh([]int{2, 7}, []int{8, 9}) {
			case equal:
				return 3, heavy
			case light:
				return 7, light
			case heavy:
				return 2, heavy
			}
		case light:
			// 5, 6, 1
			switch scale.Weigh([]int{1, 5}, []int{8, 9}) {
			case equal:
				return 6, light
			case light:
				return 5, light
			case heavy:
				return 1, heavy
			}
		case heavy:
			switch scale.Weigh([]int{0}, []int{8}) {
			case equal:
				return 4, light
			case light:
				panic(fmt.Errorf("cannot happen"))
			case heavy:
				return 0, heavy
			}
		}
	case equal:
		switch scale.Weigh([]int{8, 9}, []int{10, 0}) {
		case equal:
			switch scale.Weigh([]int{10}, []int{11}) {
			case light:
				return 11, heavy
			case heavy:
				return 11, light
			}
		case light:
			switch scale.Weigh([]int{9, 1}, []int{8, 0}) {
			case light:
				return 9, light
			case heavy:
				return 8, light
			case equal:
				return 10, heavy
			}
		case heavy:
			switch scale.Weigh([]int{9, 1}, []int{8, 0}) {
			case light:
				return 8, heavy
			case heavy:
				return 9, heavy
			case equal:
				return 10, light
			}
		}
	}
	return -1, equal
}
