package forest

import (
	"github.com/samber/lo"
	"math/rand/v2"
)

func PickIndex(rng *rand.Rand, weights []int) (index int) {
	sum := lo.Sum(weights)
	r := rng.IntN(sum)

	if len(weights) == 0 {
		index = -1
	}

	for ; index < len(weights); index++ {
		if weights[index] <= r {
			r -= weights[index]
		} else {
			break
		}
	}

	return index
}
