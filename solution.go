package main

// A solution to the 12 coins problem.
//
// A notable feature of this solution is that 3 weighings are sufficient to discriminate
// the 24 possible solutions although this is not easily proven.
//
// Another notable feature of this solution is that the multiplication 9*a + 3*b + c - 1 is the
// identity of the coin if the sum is less than 12 and is 25 - that sum if the sum is greater than
// 13. Again, why this is so is not easily proven, except by exhaustive enumeration. Note also that
// the grouping of the coins was chosen so that the sum would have this property.
//
// Suffice to say that the Venn diagrams of the coins involved in each weighing are highly
// symmetrical. There are 3 coins common to all weighings. Each weighing shares a pair
// of coins with one weighing but not the other. In one weighing that shares a pair,
// the pair that is shared will be split if the other weighing grouped it and grouped if the other weighing split it.
// Each weighing has a coin unique to it.
//
// 1,5,7 are shared by all weighings.
//
// A splits 5,7 against 1
// B groups 1,5,7 together
// C splits 1,5 against 7
//
// If A,B,C are unbalanced, then the weighings which are identical contain the counterfeit coin on the same side.
//
// 0,6 are shared by the A and B. A splits 0 and 6, B groups 0 and 6.
// 2,4 are shared by the A and C. A groups 2 and 4, C splits 2 and 4.
// 8,10 are shared by the B and C. B groups 8 and 10, C splits 8 and 10.
//
// If only C is balanced, then the counterfeit is 0 if A and B agree or 6 otherwise.
// If only B is balanced, then the counterfeit is 2 if A and C agree or 4 otherwise.
// If only A is balanced, then the counterfeit is 8 if B and C agree or 10 otherwise.
//
// 3 is unique to A
// 9 is unique to B
// 11 is unique to C
//
// If only A is unbalanced the counterfeit must be 3
// If only B is unbalanced the counterfeit must be 9
// If only C is unbalanced the counterfeit must be 11
//

func decide(scale Scale) (int, Weight) {

	a := scale.Weigh([]int{0, 3, 5, 7}, []int{1, 2, 4, 6})
	b := scale.Weigh([]int{0, 6, 8, 10}, []int{1, 5, 7, 9})
	c := scale.Weigh([]int{1, 4, 5, 8}, []int{2, 7, 10, 11})

	i := a*9 + b*3 + c
	o := i

	if i > 12 {
		o = 26 - i
	}

	f := int(o - 1)

	w := Weight((func() int {
		switch f >> 2 {
		case 0:
			return f >> 1
		case 1:
			return 1
		default:
			return 0
		}
	}() ^ f&1) << 1)

	if i > 12 {
		w = heavy - w
	}

	return f, w
}
