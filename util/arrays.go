package util

import (
	"reflect"
	"sort"
)

// Generates an array containing all possible permutations regarding element's order
// chatgpt states: Keep in mind that the number of permutations grows rapidly with the size of the input array, so this function may not be practical for very large input arrays due to the combinatorial explosion of possibilities.
// so the length of the input array should not be larger than 10, as 10! = 3628800
func GeneratePermutations[T any](array []T) [][]T {
	if len(array) == 1 {
		return [][]T{array}
	} else {
		combinations := [][]T{}
		for index, value := range array {
			others := make([]T, len(array)-1, len(array)-1)
			j := 0
			for i := 0; i < index; i++ {
				others[j] = array[i]
				j++
			}
			for i := index + 1; i < len(array); i++ {
				others[j] = array[i]
				j++
			}
			// performs somethink like this => others := array[0:index].concat(array[index+1,len(array)]), but working over sliced slices is not a good choice!
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

func HasSameValuesDisregardingOrder[T any](left, right []T, comparisionFunc func(left, right T) int) bool {
	if len(left) != len(right) {
		return false
	}
	SortSlice(left, comparisionFunc)
	SortSlice(right, comparisionFunc)
	return HasSameValuesRegardingOrder(left, right)
}

func SortSlice[T any](slice []T, comparisionFunc func(left, right T) int) {
	sort.Slice(slice, func(i, j int) bool {
		order := comparisionFunc(slice[i], slice[j])
		if order <= -1 {
			return true
		}
		return false
	})
}

func HasSameValuesRegardingOrder[T any](left, right []T) bool {
	for i := 0; i < len(left); i++ {
		if !reflect.DeepEqual(left[i], right[i]) {
			return false
		}
	}
	return true
}

// taken from https://stackoverflow.com/a/37563128/903998, thanks :D!
func Filter[T any](ss []T, predicateFunc func(T) bool) (ret []T) {
	for _, s := range ss {
		if predicateFunc(s) {
			ret = append(ret, s)
		}
	}
	return
}

func Find[T any](ss []T, predicateFunc func(T) bool) *T {
	for _, s := range ss {
		if predicateFunc(s) {
			return &s
		}
	}
	return nil
}
