package main

// decide3 given two coins from the left side of a previous weighing, one coin from the
// right side of that weighing, the result of that weighing and 2 coins known to be authentic,
// decide which of the 3 possibly counterfeit coins is counterfeit using a single weighing
//
// We put one coin from each side of the previous weighing on the left side and
// put two authentic coins on the right side.
//
// If the new weighing is balanced, the counterfeit must be the coin that wasn't weighed
// with bias of the original weighing.
//
// If the new weighing isn't balanced then if it is biased in the same direction as the
// original measurement, the counterfeit is the member of the original pair that was
// weighed twice, otherwise it is the coin from the right hand side of the previous weighing.
//
func decide3(scale Scale, weight1 Weight, two []int, one int, authentic []int) (int, Weight) {
	weight2 := scale.Weigh([]int{two[0], one}, authentic)
	if weight2 == equal {
		return two[1], weight1
	} else if weight2 == weight1 {
		return two[0], weight1
	}
	return one, weight2
}

func decide8(scale Scale, weight1 Weight) (int, Weight) {
	weight2 := scale.Weigh([]int{0, 5, 6}, []int{4, 1, 9})
	if weight2 == equal {
		return decide3(scale, weight1, []int{2, 3}, 7, []int{8, 9})
	} else if weight2 == weight1 {
		return decide3(scale, weight1, []int{0, 1}, 4, []int{8, 9})
	}
	return decide3(scale, weight2, []int{5, 6}, 1, []int{8, 9})
}

// decide returns the identity of the different coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {
	weight := scale.Weigh([]int{0, 1, 2, 3}, []int{4, 5, 6, 7})
	if weight == equal {
		weight2 := scale.Weigh([]int{8, 9}, []int{10, 0})
		if weight2 == equal {
			return 11, scale.Weigh([]int{11}, []int{0})
		} else {
			return decide3(scale, weight2, []int{8, 9}, 10, []int{0, 1})
		}
	} else {
		return decide8(scale, weight)
	}
	return -1, equal
}
