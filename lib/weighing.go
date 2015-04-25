package lib

import (
	"fmt"
)

type Weighing interface {
	Left() CoinSet
	Right() CoinSet
	Both() CoinSet
	Pan(pan int) CoinSet
	Pans() []CoinSet
}

type weighing struct {
	left  CoinSet
	right CoinSet
	both  CoinSet
}

func NewWeighing(left CoinSet, right CoinSet) Weighing {
	return &weighing{
		left:  left,
		right: right,
		both:  left.Union(right),
	}
}

func (w *weighing) Left() CoinSet {
	return w.left
}

func (w *weighing) Right() CoinSet {
	return w.right
}

func (w *weighing) Both() CoinSet {
	return w.both
}

func (w *weighing) Pan(pan int) CoinSet {
	switch pan {
	case 0:
		return w.Left()
	case 1:
		return w.Right()
	default:
		panic(fmt.Errorf("pan %d ?!", pan))
	}
}

func (w *weighing) Pans() []CoinSet {
	return []CoinSet{w.left, w.right}
}
