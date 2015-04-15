package main

// Decide returns the identity of the counterfeit coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {
	return -1, equal // always wrong
}
