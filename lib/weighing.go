package lib

type Weighing interface {
	Left() CoinSet
	Right() CoinSet
	Both() CoinSet
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
