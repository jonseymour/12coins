package main

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func decide(scale Scale) (int, Weight) {

	a := scale.Weigh([]int{9, 4, 11, 5}, []int{7, 6, 10, 8})
	b := scale.Weigh([]int{6, 1, 3, 11}, []int{4, 10, 5, 2})
	c := scale.Weigh([]int{9, 3, 7, 10}, []int{6, 0, 1, 4})

	i := (a-1)*9 + (b-1)*3 + c - 1
	o := abs(int(i)) - 1

	f := o

	// Frank Cole's solution - http://www.iwriteiam.nl/Ha12coins.html
	wabc := []Weight{a, b, c}
	t := []int{-3, 2, -2, 3, 1, 1, 2, 3, -1, 1, -1, 2}[o]
	w := wabc[abs(t)-1]
	if t < 0 {
		w = heavy - w
	}

	return f, w
}
