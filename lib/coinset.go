package lib

import (
	"encoding/json"
	"fmt"
)

const (
	ZERO_BASED = 0
	ONE_BASED  = 1
)

type CoinMask uint16

type CoinSet interface {
	AsCoins(zeroCoin int) []int
	Size() uint8
	Sort() CoinSet
	Union(other CoinSet) CoinSet
	Intersection(other CoinSet) CoinSet
	Complement(other CoinSet) CoinSet
	ExactlyOne(zeroCoin int) int
}

type hasCoinSet interface {
	asCoinSet() *coinSet
}

type coinSet struct {
	mask CoinMask
	size uint8
}

type orderedCoinSet struct {
	coinSet
	coins    []int
	zeroCoin int
}

func (s *coinSet) asCoinSet() *coinSet {
	return s
}

func (s *orderedCoinSet) asCoinSet() *coinSet {
	return &s.coinSet
}

// Answer the coins of the set as an array of integers. The
// coins are numbered w.r.t. to the specified 0 coin.
func (s *coinSet) AsCoins(zeroCoin int) []int {
	if s == nil {
		return []int{}
	}
	result := make([]int, s.Size())
	bits := s.mask
	mask := CoinMask(1)
	coin := 0
	i := uint8(0)
	for i < s.size {
		if bits&mask != 0 {
			result[i] = zeroCoin + coin
			i = i + 1
		}
		coin += 1
		mask <<= 1
	}
	return result
}

// Answer the coins of the ordered set in the original order.
func (s *orderedCoinSet) AsCoins(zeroCoin int) []int {
	c := make([]int, len(s.coins))
	diff := zeroCoin - s.zeroCoin
	for i, e := range s.coins {
		c[i] = e + diff
	}
	return c
}

// Answer the size of the set.
func (s *coinSet) Size() uint8 {
	return s.size
}

// Answer a sorted version of the current set.
func (s *coinSet) Sort() CoinSet {
	return s
}

// Answer a sorted version of the current set.
func (s *orderedCoinSet) Sort() CoinSet {
	return &s.coinSet
}

// Answer a description of the set.
func (s *coinSet) String() string {
	tmp, _ := json.Marshal(s.AsCoins(1))
	return string(tmp)
}

// Answer a description of the set.
func (s *orderedCoinSet) String() string {
	tmp, _ := json.Marshal(s.AsCoins(s.zeroCoin))
	return string(tmp)
}

// Return a new set which is the union of the receiver
// and the specified set.
func (s *coinSet) Union(other CoinSet) CoinSet {
	var o hasCoinSet
	var ok bool
	if o, ok = other.(hasCoinSet); !ok {
		o = NewCoinSet(other.AsCoins(ZERO_BASED), ZERO_BASED).(hasCoinSet)
	}
	return NewCoinSetFromMask(s.mask | o.asCoinSet().mask)
}

// Return a new set which is the intersection of the receiver
// and the specified set.
func (s *coinSet) Intersection(other CoinSet) CoinSet {
	var o hasCoinSet
	var ok bool
	if o, ok = other.(hasCoinSet); !ok {
		o = NewCoinSet(other.AsCoins(ZERO_BASED), ZERO_BASED).(hasCoinSet)
	}
	return NewCoinSetFromMask(s.mask & o.asCoinSet().mask)
}

// Return a new set which is a complement of the coins
// specified set w.r.t to the receiver's set.
func (s *coinSet) Complement(other CoinSet) CoinSet {
	var o hasCoinSet
	var ok bool
	if o, ok = other.(hasCoinSet); !ok {
		o = NewCoinSet(other.AsCoins(ZERO_BASED), ZERO_BASED).(hasCoinSet)
	}
	return NewCoinSetFromMask(s.mask &^ o.asCoinSet().mask)
}

func (s *coinSet) ExactlyOne(zeroCoin int) int {
	if s.Size() != 1 {
		panic(fmt.Errorf("illegal state: expected a set of exactly one: was %d", s.Size()))
	}
	return s.AsCoins(zeroCoin)[0]
}

// Return a new unordered set from the specified coins, assuming
// the zero coin is the specified coin.
func NewCoinSet(coins []int, zeroCoin int) CoinSet {
	mask := CoinMask(0)
	count := uint8(0)
	for _, e := range coins {
		nextmask := mask | (1 << uint(e-zeroCoin))
		if nextmask != mask {
			mask = nextmask
			count += 1
		}
	}
	return &coinSet{
		mask: mask,
		size: count,
	}
}

// Return a new unordered coin set from the specified coin mask.
func NewCoinSetFromMask(mask CoinMask) CoinSet {
	newMask := CoinMask(0)
	bit := CoinMask(1)
	count := uint8(0)
	for newMask != mask {
		if bit&mask != 0 {
			newMask |= bit
			count += 1
		}
		bit <<= 1
	}
	return &coinSet{
		mask: mask,
		size: count,
	}
}

func NewOrderedCoinSet(coins []int, zeroCoin int) CoinSet {
	result := orderedCoinSet{
		coins:    coins,
		zeroCoin: zeroCoin,
	}
	result.coinSet = *(NewCoinSet(coins, zeroCoin).(*coinSet))
	result.size = uint8(len(result.coins)) // co-erce size to be consistent with coins
	return &result
}
