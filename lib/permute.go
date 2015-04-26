package lib

import (
	"fmt"
)

func Permute(origin []int) [][]int {

	if len(origin) == 1 {
		return [][]int{origin}
	}
	results := [][]int{}
	for i, e := range origin {
		c := make([]int, len(origin)-1, len(origin)-1)
		copy(c, origin[0:i])
		copy(c[i:], origin[i+1:])
		for _, x := range Permute(c) {
			p := append([]int{e}, x...)
			results = append(results, p)
		}
	}
	return results
}

type Permutation struct {
	index []int
	zero  int
}

func NewPermutation(permutation []int, zero int) *Permutation {

	index := make([]int, len(permutation), len(permutation))

	for i, e := range permutation {
		index[e-zero] = i
	}

	return &Permutation{
		index: index,
		zero:  zero,
	}
}

func (p *Permutation) Index(e int) int {
	return p.index[e-p.zero]
}

// Generate an identifier for a permutation of N digits numbers 0 to N-1
func Number(permutation []int) uint {

	counts := make([]int, len(permutation))
	var number func([]int) (uint, uint)
	number = func(p []int) (uint, uint) {
		if len(p) == 1 {
			return 0, 1
		} else {
			h := p[0]
			t := p[1:]
			for _, e := range t {
				if e > h {
					counts[e] += 1
				}
			}
			ts, tf := number(t)
			if h < counts[h] {
				panic(fmt.Errorf("illegal state: h < counts[h]"))
			}
			rs := uint(h-counts[h])*tf + ts
			return rs, tf * uint(len(p))
		}
	}
	n, _ := number(permutation)
	return n
}
