package main

import (
	"fmt"
	"os"
)

const (
	coins      = 12
	nextBits   = 8
	coinBits   = 8
	weightBits = 2
	outcomes   = 3
)

var verbose = false
var debug = ""

var (
	weightShift = uint(coinBits)
	lightShift  = uint(64) - nextBits
	equalShift  = lightShift - nextBits
	heavyShift  = equalShift - nextBits
	leftShift   = heavyShift - coins
	rightShift  = leftShift - coins
	nextMask    = uint64(1<<nextBits - 1)
	coinsMask   = uint64(1<<coins - 1)
	coinMask    = uint64(1<<coinBits - 1)
	weightMask  = uint64(1<<weightBits - 1)
)

var table = []uint64{
	0x01190400F0F00000, //00
	0x020B070611120000, //01
	0x0300120113000000, //02
	0x0000000000000000, //03
	0x090D050611120000, //04
	0x1100060113000000, //05
	0x0000000000000200, //06
	0x0816140223000000, //07
	0x0000000000000001, //08
	0x13150A0223000000, //09
	0x0000000000000201, //0A
	0x0C0F180843000000, //0B
	0x0000000000000002, //0C
	0x17100E0843000000, //0D
	0x0000000000000202, //0E
	0x0000000000000003, //0F
	0x0000000000000203, //10
	0x0000000000000004, //11
	0x0000000000000204, //12
	0x0000000000000005, //13
	0x0000000000000205, //14
	0x0000000000000006, //15
	0x0000000000000206, //16
	0x0000000000000007, //17
	0x0000000000000207, //18
	0x1A221C3004010000, //19
	0x1B1E215000030000, //1A
	0x0000000000000008, //1B
	0x201F1D5000030000, //1C
	0x0000000000000208, //1D
	0x0000000000000009, //1E
	0x0000000000000209, //1F
	0x000000000000000A, //20
	0x000000000000020A, //21
	0x2300248000010000, //22
	0x000000000000000B, //23
	0x000000000000020B, //24
}

var index = make(map[uint64][]int)

func nextShift(w Weight) uint {
	return uint(heavy-w)*nextBits + heavyShift
}

func init() {
	verbose = len(debug) > 0
	f := func(c uint64) []int {
		b := c
		a := []int{}
		i := 0
		for b != 0 {
			if b&1 == 1 {
				a = append(a, i)
			}
			i += 1
			b >>= 1
		}
		return a
	}

	for _, e := range table {
		c := e >> leftShift & coinsMask
		if c != 0 {
			if _, ok := index[c]; !ok {
				index[c] = f(c)
			}
		}
		d := e >> rightShift & coinsMask
		if d != 0 {
			if _, ok := index[d]; !ok {
				index[d] = f(d)
			}
		}
	}
}

// Decide returns the identity of the counterfeit coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {

	i := 0
	r := table[i]
	for r>>heavyShift != 0 {
		a := index[r>>leftShift&coinsMask]
		b := index[r>>rightShift&coinsMask]
		w := scale.Weigh(a, b)
		p := r >> nextShift(w) & nextMask
		if verbose {
			fmt.Fprintf(os.Stderr, "%d: %016x, %v, %v & %v -> do %d\n", i, r, a, b, w, p)
		}
		i = int(p)
		r = table[i]
	}
	coin := int(r & coinMask)
	weight := Weight(r >> weightShift & weightMask)
	if verbose {
		fmt.Fprintf(os.Stderr, "%d: %016x -> stop %d, %v\n", i, r, coin, weight)
	}
	return coin, weight
}
