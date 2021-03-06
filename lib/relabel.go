package lib

import (
	"sort"
)

// Relabel the coins of the weighing such that the Coins
// slice is numbered in strictly increasing numerical order.
func (s *Solution) Relabel() (*Solution, error) {

	var clone *Solution
	var err error

	if s.flags&REVERSED == 0 {
		if clone, err = s.Reverse(); err != nil {
			return s, err
		}
	} else {
		clone = s.Clone()
	}

	clone.reset()

	c := make([]int, len(clone.Coins), len(clone.Coins))
	for i, e := range clone.Coins {
		c[i] = e
	}
	z := clone.GetZeroCoin()
	p := NewPermutation(c, z)

	for i, w := range clone.Weighings {
		coinSet := [2]CoinSet{}
		for j, pan := range w.Pans() {
			coins := pan.AsCoins(z)
			for k, e := range coins {
				coins[k] = p.Index(e) + z
			}
			sort.Sort(sort.IntSlice(coins))
			coinSet[j] = NewCoinSet(coins, z)
		}
		clone.Weighings[i] = NewWeighing(coinSet[0], coinSet[1])
	}

	for i, _ := range clone.Coins {
		clone.Coins[i] = i + z
	}

	clone.flags |= (RELABLED | NORMALISED)

	return clone, nil
}
