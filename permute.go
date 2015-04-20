package main

import ()

func permute(origin []int) [][]int {

	if len(origin) == 1 {
		return [][]int{origin}
	}
	results := [][]int{}
	for i, e := range origin {
		c := make([]int, len(origin)-1, len(origin)-1)
		copy(c, origin[0:i])
		copy(c[i:], origin[i+1:])
		for _, x := range permute(c) {
			p := append([]int{e}, x...)
			results = append(results, p)
		}
	}
	return results
}
