package lib

import (
	"fmt"
)

type CoinMask uint16

type CoinSet interface {
	AsCoins(zeroCoin int) []int
	Size() uint8
	Sort() CoinSet
	Union(other CoinSet) CoinSet
	Intersection(other CoinSet) CoinSet
	Complement(other CoinSet) CoinSet
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

func (s *orderedCoinSet) AsCoins(zeroCoin int) []int {
	c := make([]int, len(s.coins))
	diff := zeroCoin - s.zeroCoin
	for i, e := range s.coins {
		c[i] = e + diff
	}
	return c
}

func (s *coinSet) Size() uint8 {
	return s.size
}

func (s *coinSet) Sort() CoinSet {
	return s
}

func (s *orderedCoinSet) Sort() CoinSet {
	return &s.coinSet
}

func (s *coinSet) String() string {
	return fmt.Sprintf("%v", s.AsCoins(1))
}

func (s *orderedCoinSet) String() string {
	return fmt.Sprintf("%v", s.AsCoins(s.zeroCoin))
}

func (s *coinSet) Union(other CoinSet) CoinSet {
	var o *coinSet
	var ok bool
	if o, ok = other.(*coinSet); !ok {
		o = NewCoinSet(other.AsCoins(0), 0).(*coinSet)
	}
	return NewCoinSetFromMask(s.mask | o.mask)
}

func (s *coinSet) Intersection(other CoinSet) CoinSet {
	var o *coinSet
	var ok bool
	if o, ok = other.(*coinSet); !ok {
		o = NewCoinSet(other.AsCoins(0), 0).(*coinSet)
	}
	return NewCoinSetFromMask(s.mask & o.mask)
}

func (s *coinSet) Complement(other CoinSet) CoinSet {
	var o *coinSet
	var ok bool
	if o, ok = other.(*coinSet); !ok {
		o = NewCoinSet(other.AsCoins(0), 0).(*coinSet)
	}
	return NewCoinSetFromMask(s.mask &^ o.mask)
}

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
