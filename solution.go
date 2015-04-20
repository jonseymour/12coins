package main

import (
	"github.com/jonseymour/12coins/lib"
)

// Decide returns the identity of the counterfeit coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale lib.Scale) (int, lib.Weight) {
	return -1, lib.Equal // always wrong
}
