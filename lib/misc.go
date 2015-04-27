package lib

// Calculate the absolute value of the specified integer.
func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

// convert an integer value into a pointer to that value.
func pi(i int) *int {
	return &i
}

func pu(u uint) *uint {
	return &u
}

// convert a boolean value into a pointer to that value.
func pbool(b bool) *bool {
	return &b
}

// calculate n factorial
func fact(n int) int {
	if n < 2 {
		return 1
	} else {
		return n * fact(n-1)
	}
}
