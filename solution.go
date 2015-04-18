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
	a := scale.Weigh([]int{0, 1, 2, 3}, []int{4, 5, 6, 7})   // 8,9,10,11 // a loses 3 from b, 3 from c
	b := scale.Weigh([]int{0, 6, 9, 11}, []int{2, 3, 4, 10}) // 1,5,7,8   // b loses 3 from a, 3 from c
	c := scale.Weigh([]int{2, 4, 7, 11}, []int{3, 5, 8, 9})  // 0,1,6,10  // c loses 3 from a, 3 from b

	/*
		Documenting the permutation

		c[0] = b[4] = a[2]   // a,b,c
		c[1] = b[6] = a[4]   // a,b,c
		c[2] = b[10] = a[7]  // a,c
		c[3] = b[3] = a[11]  // b,c
		c[4] = b[5] = a[3]   // a,b,c
		c[5] = b[9] = a[5]   // a,c
		c[6] = b[11] = a[8]  // c
		c[7] = b[2] = a[9]   // b,c
		c[8] = b[0] = a[0]   // a,b
		c[9] = b[8] = a[1]   // a
		c[10] = b[1] = a[6]  // a,b
		c[11] = b[7] = a[10] // b
	*/

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
