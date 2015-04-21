package main

import (
	"fmt"
	"github.com/jonseymour/12coins/lib"
	"os"
)

const (
	light = lib.Light
	heavy = lib.Heavy
)

var (
	A = [2][]int{
		[]int{3, 5, 1, 7}, []int{6, 8, 2, 4},
	}
	B = [2][]int{
		[]int{6, 11, 8, 1}, []int{9, 2, 7, 10},
	}
	C = [2][]int{
		[]int{3, 12, 8, 2}, []int{6, 9, 11, 5},
	}

	IDENTITY = [12]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	WEIGHTS  = [12]lib.Weight{
		light, heavy, light, heavy,
		light, heavy, light, heavy,
		heavy, heavy, light, light,
	}
)

func main() {
	ref := &lib.Solver{
		Weighings: [3][2][]int{A, B, C},
		Coins:     IDENTITY,
		Weights:   WEIGHTS,
		ZeroCoin:  1,
	}
	for _, m := range []bool{false, true} {
		for _, p := range lib.Permute([]int{0, 1, 2}) {
			clone := ref.Clone()
			clone.Mirror = m
			clone.Permutation = p
			clone.Weighings = [3][2][]int{clone.Weighings[p[0]], clone.Weighings[p[1]], clone.Weighings[p[2]]}
			clone, _ = clone.Reverse()
			clone = clone.Relabel()
			if errors := lib.TestAll(clone.Decide); len(errors) != 0 {
				panic(fmt.Errorf("errors: %v", errors))
			}
			fmt.Fprintf(os.Stdout, "%s\n", clone)
		}
	}
}
