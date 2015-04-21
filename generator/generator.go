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
		[]int{2, 4, 0, 6}, []int{5, 7, 1, 3},
	}
	B = [2][]int{
		[]int{5, 10, 7, 0}, []int{8, 1, 6, 9},
	}
	C = [2][]int{
		[]int{2, 11, 7, 1}, []int{5, 8, 10, 4},
	}

	IDENTITY = [12]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
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
	}
	for _, p := range lib.Permute([]int{0, 1, 2}) {
		clone := ref.Clone()
		clone.Permutation = p
		clone.Weighings = [3][2][]int{clone.Weighings[p[0]], clone.Weighings[p[1]], clone.Weighings[p[2]]}
		for i, _ := range IDENTITY {
			o := lib.NewOracle(i, lib.Light)
			ri, rw := clone.Decide(o)
			if ri != i {
				clone.Coins[ri] = i
			}
			if rw != lib.Light {
				clone.Weights[ri] = lib.Heavy - clone.Weights[ri]
			}
		}
		clone.Relabel()
		fmt.Fprintf(os.Stdout, "%s\n", clone)
	}

}
