package lib

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
