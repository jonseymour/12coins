#12 coins

Given:

* 11 coins of identical weight
* 1 coin of a different relative weight
* a set of scales
* a maximum of 3 weighings

Decide which coin has the different weight and whether that weight is less than
or greater than the weight of any of the other 11 coins.

The decide() function is a go-lang solution to this problem. When compiled and
executed, the program tests that function against all 24 possible
configurations.

#The Solution

The solution documented in solution.go is a very pleasing solution.

Unlike the other attempts, this solution only requires 3 weighings and these
weighings are sufficient to discriminate the 24 different possibilities for any
given configuration of the coins.

	func decide(scale Scale) (int, Weight) {

		a := scale.Weigh([]int{0, 3, 5, 7}, []int{1, 2, 4, 6})
		b := scale.Weigh([]int{0, 6, 8, 10}, []int{1, 5, 7, 9})
		c := scale.Weigh([]int{1, 4, 5, 8}, []int{2, 7, 10, 11})

		i := a*9 + b*3 + c
		o := i

		if i > 12 {
			o = 26 - i
		}

		f := int(o - 1)

		w := Weight((func() int {
			switch f >> 2 {
			case 0:
				return f >> 1
			case 1:
				return 1
			default:
				return 0
			}
		}() ^ f&1) << 1)


		if i > 12 {
			w = heavy - w
		}

		return f, w
	}

The following Venn diagram, which shows the intersections between the sets of
coins involved in all 3 weighings, helps to provide a heuristic justification
for why this set of weighings is capable of discriminating the 24 cases - each
weighing involves overlapping set and subsets of coins and the 12 coins are
evenly distributed across all sets and all intersections between all sets.

<img src="venn.png"/>

Some observations:

- all weighings share 3 coins {1,5,7}
- each weighing shares a different pair of coins with each other weighing
- each weighing has a single coin that is unique to itself
- each weighing shares exactly 5 coins with one weighing and a different (but partially overlapping) set of 5 coins with the other weighing

#Explanation Of Completeness

The following argument explains why the configuration of the weighings has
enough information to distinguish the 24 possible configurations. It doesn't
explain why the coin is selected using manipulations of a sum derived from the 3
weighings.

If the B and C weighings are balanced, then the A weighing must be unbalanced
because of the coin unique to A - namely 3.

If the A and B weighings are unbalanced and the C weighing is balanced, then the
cause must be a coin that is common to A and B and not shared by C, namely 0 or
6. If the A and B weighings have the same bias, then the counterfeit coin must
be 0, otherwise it is 6.

If the A, B and C weighings are unbalanced, then the cause must be a coin that
is common to A, B and C - namely 1,5 or 7. If the A and C weighings have the
same bias, then the counterfeit is 5. If B and C have the same bias, then the
counterfeit is 7. If A and B have the same bias, then the counterfeit is 1.

Symmetry arguments allow derivation of other the possible solutions - 2, 4, 8,
9, 10, 11.

#Explanation Of Indexing Behaviour Of Sum

The indexing behaviour of the sum which answers the identity of the coin appears
somewhat magical, and indeed, if I had happened upon a distribution of weights
that had this property it would have been amazing. In reality, in an earlier
iteration I discovered a set of weighings that had the ability to discriminate
the 24 configurations. I then used the sum to calculate an index into two arrays
of length 27 and stored the identity of the coin and the weight of the
counterfeit coin at the indexed element. This array served as a mapping function
from the sum to the identity of the coin.

When I did this, I observed that elements at 0, 13 and 26 were not assigned to
and that coins[13+(n+1)] = coins[13-(n+1)] and weights[13+(n+1)] = 2 -
weights[13-(n+1)] for all n in [0,11]. This realisation allowed the mapping
function to be realised using 2 arrays of 12 elements each and a test for the
magnitude of the sum.

Having done this, I then identified a permutation that allowed me to relabel the
coins such that the content of the i'th element was i, thereby allowing me to
replace the coin mapping function with the identity function and so eliminate
the need for this array in the solution.

Observe that for each set of coins A, B, C one weighing acts as the weight bit
and the other weighings act as an indexing function for the coin. For example,
for the coins 0-7 and the corresponding weights
[light,heavy,heavy,light,heavy,light,heavy,light], the A weighings will all be
biased in the light direction (e.g. produce a trit of 0). The other weighings
enumerate the 8 trits between 01 and 22. Observe that the B and C weighings
never register 00 (light, light) for any of these coins because any coin that
causes weighing A to produce 0 cannot simultaneously cause both the B and C
weighings to produce 0. Also observe that when the B and C weighings register 11
(equal, equal), then the A weighing directly reveales the weight of the
counterfeit coin that is unique to the A weighing.

#Explanation Of Weight Deriviation Function

One thing that bugged me about the original solution was the somewhat arbitrary nature
of the weight derivation function.

		w := Weight((func() int {
			switch f >> 2 {
			case 0:
				return f >> 1
			case 1:
				return 1
			default:
				return 0
			}
		}() ^ f&1) << 1)

		if i>12 {
			w = heavy - w
		}

This function had been reversed engineered from the mapping of 9a*3b*c-1 via a
table that gave the correct weight for each configuration and certainly can't be
justified merely by inspecting the code.

I did notice that the function appeared to have 3 branches controlled by the
higher order bits of sum - in effect, it was a function that was composed of 3
other functions although the significance of this was lost on me at the time.

I also realised that the choice of linear combination I chose was somewhat
arbitrary: I chose 9a + 3b + c, but I could have just as easily selected
9b + 3a + c or any of the other 4 permutations of a,b,c that are possible.

So, I shuffled the sum. This broke the identity function and changed the weight
derivation function so I reintroduced arrays to implement those functions. I
kept shuffling until I stumbled across a weight deriviation function that had a
simpler structure that the original one. I then applied the re-labeling
trick to the resulting solution and ended up with this:

	func decide(scale Scale) (int, Weight) {

		a := scale.Weigh([]int{2, 4, 0, 6}, []int{5, 7, 1, 3})
		b := scale.Weigh([]int{5, 10, 7, 0}, []int{8, 1, 6, 9})
		c := scale.Weigh([]int{2, 11, 7, 1}, []int{5, 8, 10, 4})

		i := a*9 + b*3 + c
		o := i

		if i > 12 {
			o = 26 - i
		}

		f := int(o - 1)

		w := Weight((func() int {
			if f&8 == 0 {
				return f
			} else {
				return 1 ^ f>>1
			}
		})() & 1 << 1)

		if i > 12 {
			w = heavy - w
		}

		return f, w
	}

This is a much simpler function as the weight varies according to the first bit
of the coin identity for the first 8 coins and according to the negation of the
2nd bit for the remaining 4.

Another interesting thing is that while the order of coins in each weighing has
changed, each weighing has the same coins as the original solution. My
conjecture, therefore, is that the indexing property is a function of which
coins are in each weighing and the weight derivation function is a function of
the permutation of the coins in each weighing.  I intend to explore this further
by generating the 4 remaining solutions of this type.

The fact that the coins 0 -> 7 appear in the first weighing of both solutions
is indubitably related to the fact that for these coins to be indexed in the first
8 positions by 9a+3b+3c-1, the first weighing trit (e.g. a) has to be 0 for these 8 coins.
This can only be true for the 4 configurations where one of the coins on the left hand side
is light and the 4 configurations where one of the coins on the right hand side is heavy.

Presumably this style of argument can be extended so show why the other weighings
must contain the coins they do.

#Other notes

Adding one to each coin identifier (so they are numbered 1->12 instead of 0-11),
yields this Venn diagram where the even coins are shared by an odd number of
weighings and the odd coins by an even number of weighings.

<img src="venn-1-based.png">

Converting these numbers to base 3 yeilds a diagram which shows that the coins
in the intersection between two sets share the same base 3 digit and this digit
is also shared by the coins unique to each set.

<img src="venn-base-3.png">

