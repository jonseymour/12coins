package main

// Decide returns the identity of the counterfeit coin and what the relative
// weight of that coin is with respect to any other coin.

var (
	coins   = []int{0, 4, 5, 1, 7, 2, 6, 3, 11, 10, 9, 8}
	weights = []Weight{
		light, heavy, heavy, light,
		heavy, light, heavy, light,
		light, heavy, light, heavy,
	}
)

// The simplest possible solution. 3 comparisons + 2 table lookup, deals with every case.
func decide(scale Scale) (int, Weight) {
	a := scale.Weigh([]int{0, 1, 2, 3}, []int{4, 5, 6, 7})
	b := scale.Weigh([]int{0, 6, 9, 11}, []int{2, 3, 4, 10})
	c := scale.Weigh([]int{2, 4, 7, 11}, []int{3, 5, 8, 9})
	i := a*9 + b*3 + c
	o := i

	if i > 12 {
		o = 26 - i
	}
	f := coins[o-1]
	w := weights[o-1]

	if i > 12 {
		w = heavy - w
	}

	return f, w
}
