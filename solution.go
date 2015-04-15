package main

// Decide returns the identity of the counterfeit coin and what the relative
// weight of that coin is with respect to any other coin.
func decide(scale Scale) (int, Weight) {
    var index int = -1
    var weight Weight = equal
    var result = scale.Weigh([]int{ 0, 1, 2, 3 }, []int{ 4, 5, 6, 7 })
    if result == equal {
        //All coins 0 .. 7 are good
        result = scale.Weigh([]int{ 8, 9, 10 }, []int{ 0, 1, 2 })
        if result == equal {
            //8, 9 and 10 are good too
            result = scale.Weigh([]int{ 11 }, []int{ 0 })
            index = 11
            weight = result
        } else if result == heavy {
            result = scale.Weigh([]int{ 8 }, []int{ 9 })
            if result == equal {
                index = 10
                weight = heavy
            } else if result == heavy {
                index = 8
                weight = heavy
            } else {
                index = 9
                weight = heavy
            }
        } else {
            result = scale.Weigh([]int{ 8 }, []int{ 9 })
            if result == equal {
                index = 10
                weight = light
            } else if result == heavy {
                index = 9
                weight = light
            } else {
                index = 8
                weight = light
            }
        }
    } else if result == heavy {
        //Either one of 0,1,2,3 is a heavy fake or one of 4,5,6,7 is a light fake
        result = scale.Weigh([]int{ 0, 4, 5 }, []int{ 1, 6, 7 })
        if result == equal {
            //One of 2, 3 is a heavy fake
            result = scale.Weigh([]int{ 2 }, []int{ 3 })
            weight = heavy
            if result == heavy {
                index = 2
            } else {
                index = 3
            }
        } else if result == heavy {
            //Either 0 is a heavy fake or one of 6, 7 is a light fake
            result = scale.Weigh([]int{ 6 }, []int{ 7 })
            if result == equal {
                index = 0
                weight = heavy
            } else if result == heavy {
                index = 7
                weight = light
            } else {
                index = 6
                weight = light
            }
        } else {
            //Either 1 is a heavy fake or one of 4, 5 is a light fake
            result = scale.Weigh([]int{ 4 }, []int{ 5 })
            if result == equal {
                index = 1
                weight = heavy
            } else if result == heavy {
                index = 5
                weight = light
            } else {
                index = 4
                weight = light
            }
        }
    } else {
        //Either one of 0,1,2,3 is a light fake or one of 4,5,6,7 is a heavy fake
        result = scale.Weigh([]int{ 0, 4, 5 }, []int{ 1, 6, 7 })
        if result == equal {
            //One of 2, 3 is a light fake
            result = scale.Weigh([]int{ 2 }, []int{ 3 })
            weight = light
            if result == heavy {
                index = 3
            } else {
                index = 2
            }
        } else if result == heavy {
            //Either one of 4, 5 is a heavy fake or 1 is a light fake
            result = scale.Weigh([]int{ 4 }, []int{ 5 })
            if result == equal {
                index = 1
                weight = light
            } else if result == heavy {
                index = 4
                weight = heavy
            } else {
                index = 5
                weight = heavy
            }
        } else {
            //Either one of 6, 7 is a heavy fake or 0 is a light fake
            result = scale.Weigh([]int{ 6 }, []int{ 7 })
            if result == equal {
                index = 0
                weight = light
            } else if result == heavy {
                index = 6
                weight = heavy
            } else {
                index = 7
                weight = heavy
            }
        }
    }
	return index, weight // always wrong
}
