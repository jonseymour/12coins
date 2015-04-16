package main

func decide8(scale Scale, weight1 Weight) (int, Weight) {
	weight2 := scale.Weigh([]int{0, 5, 6}, []int{4, 1, 9})

	if weight2 == equal {
		weight3 := scale.Weigh([]int{2, 7}, []int{8, 9})
		if weight3 == equal {
			return 3, weight1
		} else if weight3 == weight1 {
			return 2, weight1
		}
		return 7, weight1.invert()
	} else if weight2 == weight1 {
		// # 0, 4
		if scale.Weigh([]int{0}, []int{8}) == weight1 {
			return 0, weight1
		} else {
			return 4, weight1.invert()
		}
	}
	// 5, 6, 1
	weight3 := scale.Weigh([]int{1, 5}, []int{8, 9})
	if weight3 == equal {
		return 6, weight1.invert()
	} else if weight3 == weight1 {
		return 1, weight1
	}
	return 5, weight1.invert()
}

// decide returns the identity of the different coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {
	weight := scale.Weigh([]int{0, 1, 2, 3}, []int{4, 5, 6, 7})
	if weight == equal {
		switch scale.Weigh([]int{8, 9}, []int{10, 0}) {
		case equal:
			return 11, scale.Weigh([]int{11}, []int{0})
		case light:
			switch scale.Weigh([]int{9}, []int{8}) {
			case light:
				return 9, light
			case heavy:
				return 8, light
			case equal:
				return 10, heavy
			}
		case heavy:
			switch scale.Weigh([]int{9}, []int{8}) {
			case light:
				return 8, heavy
			case heavy:
				return 9, heavy
			case equal:
				return 10, light
			}
		}
	} else {
		return decide8(scale, weight)
	}
	return -1, equal
}
