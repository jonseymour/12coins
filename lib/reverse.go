package lib

import (
	"fmt"
)

//
// If the receiver is a valid solution to the 12 coins problem,
// return a clone of the receiver in which the Coins and Weights
// slice and the Flip pointer have been populated with values
// required to make Decide(Scale) return the correct values for
// all inputs.
//
// Otherwise, return a pointer to the receiver.
//
// The .flags value of the returned pointer is INVALID if the solution
// is invalid or REVERSED if the receiver is a valid solution to the
// the problem.
//
// If err is non-nil, then .Valid of the result will point to a false
// value and .Failures of the result will list the tests that
// caused the receiver to be marked invalid.
//
func (s *Solution) Reverse() (*Solution, error) {
	clone := s.Clone()

	clone.reset()
	clone.markInvalid()

	clone.Coins = make([]int, 27, 27)
	clone.Weights = make([]Weight, 27, 27)
	if clone.flags&RECURSE == 0 {
		clone.encoding.Flip = nil
	}

	for i, _ := range clone.Coins {
		clone.Coins[i] = clone.GetZeroCoin()
		clone.Weights[i] = Equal
	}

	failures := make(map[int]bool)

	fail := func(coin int, weight Weight) {
		k := coin*2 + (int(weight) / 2)
		if !failures[k] {
			failures[k] = true
			s.Failures = append(s.Failures, Failure{
				Coin:   coin,
				Weight: weight,
			})
		}
	}

	for _, w := range []Weight{Light, Heavy} {
		for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {
			o := NewOracle(i, w, 1)
			ri, _, rx := clone.decide(o)
			if ri != i {
				if clone.Weights[rx] != Equal {
					fail(clone.Coins[rx], clone.Weights[rx])
					fail(i, w)
					continue
				}
				clone.Coins[rx] = i
			}
			clone.Weights[rx] = w
		}
	}

	if len(s.Failures) != 0 {
		s.flags = INVALID
		return s, fmt.Errorf("not a valid solution because of %d failures", len(s.Failures))
	}

	// exploit symmetry where it exists

	if clone.Weights[0] == Equal {
		clone.Coins = clone.Coins[1:13]
		clone.Weights = clone.Weights[1:13]
	} else {

		//
		// A curious truth is that within the first 13 positions, of the
		// array, the unassigned slots can only ever occur at positions
		// 0,2,6 or 8.
		//
		// If it occurs an 0, then LLL is not a valid weighing. If it occurs
		// at 2, then LLH is not a valid weighing. If it occurs at 6, then LHL
		// is not a valid weighing. If it occurs at 8 then LHH is not a valid
		// weighing.
		//
		// The desired outcome is that the empty slots occur at
		// 0 which allows the sum derived from the other bits to
		// index the counterfeit coin directly.
		//
		// This can be arranged by identifying the weighing that is causing
		// the empty slot to happen at something other than 0 and flipping
		// the contribution of that weighing to the sum.
		//

		if clone.Weights[8] == Equal {
			clone.encoding.Flip = pi(0)
		} else if clone.Weights[6] == Equal {
			clone.encoding.Flip = pi(1)
		} else if clone.Weights[2] == Equal {
			clone.encoding.Flip = pi(2)
		} else {
			panic(fmt.Errorf("unexpected case: %v!", clone.Weights))
		}
		if clone.flags&RECURSE != 0 {
			panic(fmt.Errorf("infinite recursion detected: %v", clone.Weights))
		}
		defer func() {
			clone.flags = clone.flags &^ RECURSE
		}()
		clone.flags |= RECURSE
		return clone.Reverse()
	}
	clone.flags |= REVERSED
	return clone, nil
}
