package lib

type CoinSet interface {
	AsCoins(zeroCoin int) []int
	Size() uint8
	Sort() CoinSet
}

type coinset struct {
	mask uint16
	size uint8
}

func (s *coinset) AsCoins(zeroCoin int) []int {
	result := make([]int, s.Size())
	bits := s.mask
	mask := uint16(1)
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

func NewCoinSet(coins []int, zeroCoin int) CoinSet {
	mask := uint16(0)
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
