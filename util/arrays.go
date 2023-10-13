package util

// Generates an array containing all possible permutations regarding element's order
// chatgpt states: Keep in mind that the number of permutations grows rapidly with the size of the input array, so this function may not be practical for very large input arrays due to the combinatorial explosion of possibilities.
// so the length of the input array should not be larger than 10, as 10! = 3628800
func GeneratePermutations[T any](array []T) [][]T {
	if len(array) == 1 {
		return [][]T{array}
	} else {
		combinations := [][]T{}
		for index, value := range array {
			others := make([]T,len(array)-1,len(array)-1)
			j := 0
			for i := 0; i < index; i++ {
				others[j] = array[i]
				j++
			}
			for i := index + 1; i < len(array); i++ {
				others[j] = array[i]
				j++
			}
			subcombinations := GeneratePermutations(others)
			for _, subcombination := range subcombinations {
				combination := append([]T{value}, subcombination...)
				combinations = append(combinations, combination)
			}
		}
		return combinations
	}
}

func DeepCopySlice[T any](original []T) []T {
	// Create a new slice with the same length as the original.
	copied := make([]T, len(original))

	// Copy each element from the original slice to the new slice.
	for i, v := range original {
		copied[i] = v
	}

	return copied
}

func ShallowCopySlice[T any](original []T) []T {
	copied := make([]T, len(original))
	copy(copied, original)
	return copied
}


func Flatten[T any](lists [][]T) []T {
    var res []T
    for _, list := range lists {
        res = append(res, list...)
    }
    return res
}