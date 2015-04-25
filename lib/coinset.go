package lib

type CoinMask uint16

type CoinSet interface {
	AsCoins(zeroCoin int) []int
	Size() uint8
	Sort() CoinSet
	Union(other CoinSet) CoinSet
	Intersection(other CoinSet) CoinSet
	Add(coin int, zeroCoin int) CoinSet
}

type coinset struct {
	mask CoinMask
	size uint8
}

func (s *coinset) AsCoins(zeroCoin int) []int {
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

func (s *coinset) Size() uint8 {
	return s.size
}

func (s *coinset) Sort() CoinSet {
	return s
}

func (s *coinset) Union(other CoinSet) CoinSet {
	var o *coinset
	var ok bool
	if o, ok = other.(*coinset); !ok {
		o = NewCoinSet(other.AsCoins(0), 0).(*coinset)
	}
	return NewCoinSetFromMask(s.mask | o.mask)
}

func (s *coinset) Intersection(other CoinSet) CoinSet {
	var o *coinset
	var ok bool
	if o, ok = other.(*coinset); !ok {
		o = NewCoinSet(other.AsCoins(0), 0).(*coinset)
	}
	return NewCoinSetFromMask(s.mask & o.mask)
}

func (s *coinset) Add(coin int, zeroCoin int) CoinSet {
	mask := CoinMask(1 << uint(coin-zeroCoin))
	if (s.mask & mask) == 0 {
		return NewCoinSetFromMask(s.mask | mask)
	} else {
		return s
	}
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
	return &coinset{
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
	return &coinset{
		mask: mask,
		size: count,
	}
}
