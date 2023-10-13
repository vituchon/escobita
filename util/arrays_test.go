package util

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestGeneratePermutationsWorks(t *testing.T) {
	testRuns := []struct {
		title    string
		input    []int
		expected [][]int
	}{
		{
			title:    "Empty array",
			input:    []int{},
			expected: [][]int{},
		},
		{
			title:    "2 element array",
			input:    []int{1, 2},
			expected: [][]int{[]int{1, 2}, []int{2, 1}},
		},
		{
			title:    "3 element array",
			input:    []int{1, 2, 3},
			expected: [][]int{[]int{1, 2, 3}, []int{1, 3, 2}, []int{2, 1, 3}, []int{2, 3, 1}, []int{3, 1, 2}, []int{3, 2, 1}},
		},
		{
			title:    "4 element array",
			input:    []int{1, 2, 3, 4},
			expected: [][]int{[]int{1, 2, 3, 4}, []int{1, 2, 4, 3}, []int{1, 3, 2, 4}, []int{1, 3, 4, 2}, []int{1, 4, 2, 3}, []int{1, 4, 3, 2}, []int{2, 1, 3, 4}, []int{2, 1, 4, 3}, []int{2, 3, 1, 4}, []int{2, 3, 4, 1}, []int{2, 4, 1, 3}, []int{2, 4, 3, 1}, []int{3, 1, 2, 4}, []int{3, 1, 4, 2}, []int{3, 2, 1, 4}, []int{3, 2, 4, 1}, []int{3, 4, 1, 2}, []int{3, 4, 2, 1}, []int{4, 1, 2, 3}, []int{4, 1, 3, 2}, []int{4, 2, 1, 3}, []int{4, 2, 3, 1}, []int{4, 3, 1, 2}, []int{4, 3, 2, 1}},
		},
	}

	for _, testRun := range testRuns {
		t.Log("Running", testRun.title)
		computed := GeneratePermutations(testRun.input)
		if len(computed) != len(testRun.expected) {
			t.Error("Length doesn't match, computed was: ", len(computed), " and expected is: ", len(testRun.expected))
		}
		if !HasSameValuesDisregardingOrder(computed, testRun.expected) {
			t.Error("Computed was", computed, "and expected is", testRun.expected)
		}
	}
}

func sortIntSliceSlice(slice [][]int) {
	sort.Slice(slice, func(i, j int) bool {
		iSlice := slice[i]
		jSlice := slice[j]
		for k := 0; k < len(iSlice); k++ {
			if iSlice[k] < jSlice[k] {
				return true
			} else if iSlice[k] > jSlice[k] {
				return false
			}
			// if they are equals then continue until there is a order difference or the slice ends
		}
		return false
	})
}

func TestSortIntSliceSliceWork(t *testing.T) {
	testRuns := []struct {
		title    string
		input    [][]int
		expected [][]int
	}{
		{
			title:    "Empty",
			input:    [][]int{},
			expected: [][]int{},
		},
		{
			title:    "1 element (single, sorted)",
			input:    [][]int{[]int{1}},
			expected: [][]int{[]int{1}},
		},
		{
			title:    "1 element (double, sorted)",
			input:    [][]int{[]int{1, 2}},
			expected: [][]int{[]int{1, 2}},
		},
		{
			title:    "2 elements (single, sorted)",
			input:    [][]int{[]int{1}, []int{2}},
			expected: [][]int{[]int{1}, []int{2}},
		},
		{
			title:    "2 elements (single, not sorted)",
			input:    [][]int{[]int{2}, []int{1}},
			expected: [][]int{[]int{1}, []int{2}},
		},
		{
			title:    "2 elements (double, sorted)",
			input:    [][]int{[]int{1, 2}, []int{2, 1}},
			expected: [][]int{[]int{1, 2}, []int{2, 1}},
		},
		{
			title:    "2 elements (double, not sorted)",
			input:    [][]int{[]int{2, 1}, []int{1, 2}},
			expected: [][]int{[]int{1, 2}, []int{2, 1}},
		},
		{
			title:    "3 elements (triple, not sorted)",
			input:    [][]int{[]int{4, 1, 1}, []int{1, 4, 5}, []int{1, 3, 5}},
			expected: [][]int{[]int{1, 3, 5}, []int{1, 4, 5}, []int{4, 1, 1}},
		},
		{
			title:    "4 elements (single, not sorted)",
			input:    [][]int{[]int{2}, []int{5}, []int{1}, []int{3}},
			expected: [][]int{[]int{1}, []int{2}, []int{3}, []int{5}},
		},
	}
	for _, testRun := range testRuns {
		t.Log("Running", testRun.title)
		computed := DeepCopySlice(testRun.input)
		sortIntSliceSlice(computed)
		for i := 0; i < len(testRun.expected); i++ {
			if !reflect.DeepEqual(computed[i], testRun.expected[i]) {
				t.Error("Computed was", computed, "and expected is", testRun.expected)
				break
			}
		}
	}
}

func HasSameValuesDisregardingOrder(left, right [][]int) bool {
	if len(left) != len(right) {
		return false
	}
	sortIntSliceSlice(left)
	sortIntSliceSlice(right)
	return hasSameValuesRegardingOrder(left, right)
}

func hasSameValuesRegardingOrder(left, right [][]int) bool {
	for i := 0; i < len(left); i++ {
		if !reflect.DeepEqual(left[i], right[i]) {
			return false
		}
	}
	return true
}

func TestHasSameValuesDisregardingOrder(t *testing.T) {
	testRuns := []struct {
		title    string
		left     [][]int
		right    [][]int
		expected bool
	}{
		{
			title:    "Empty",
			left:     [][]int{},
			right:    [][]int{},
			expected: true,
		},
		{
			title:    "1 element (single, equals)",
			left:     [][]int{[]int{1}},
			right:    [][]int{[]int{1}},
			expected: true,
		},
		{
			title:    "1 element (single, not equals: different values)",
			left:     [][]int{[]int{1}},
			right:    [][]int{[]int{2}},
			expected: false,
		},
		{
			title:    "1 element (single, not equals: different size)",
			left:     [][]int{[]int{1}},
			right:    [][]int{},
			expected: false,
		},
		{
			title:    "1 element (double, equals)",
			left:     [][]int{[]int{1, 1}},
			right:    [][]int{[]int{1, 1}},
			expected: true,
		},
		{
			title:    "1 element (double, not equals: different values)",
			left:     [][]int{[]int{1, 1}},
			right:    [][]int{[]int{1, 2}},
			expected: false,
		},
		{
			title:    "2 elements (single, equals, same order)",
			left:     [][]int{[]int{1}, []int{2}},
			right:    [][]int{[]int{1}, []int{2}},
			expected: true,
		},
		{
			title:    "2 elements (single, equals, different order)",
			left:     [][]int{[]int{1}, []int{2}},
			right:    [][]int{[]int{2}, []int{1}},
			expected: true,
		},
		{
			title:    "2 elements (single, not equals: different values)",
			left:     [][]int{[]int{1}, []int{2}},
			right:    [][]int{[]int{2}, []int{2}},
			expected: false,
		},
		{
			title:    "2 elements (single, not equals: different size)",
			left:     [][]int{[]int{1}},
			right:    [][]int{[]int{2}, []int{1}},
			expected: false,
		},
		{
			title:    "3 elements (triple, equals, same order)",
			left:     [][]int{[]int{1, 3, 5}, []int{1, 4, 5}, []int{4, 1, 1}},
			right:    [][]int{[]int{1, 3, 5}, []int{1, 4, 5}, []int{4, 1, 1}},
			expected: true,
		},
		{
			title:    "3 elements (triple, equals, different order)",
			left:     [][]int{[]int{4, 1, 1}, []int{1, 4, 5}, []int{1, 3, 5}},
			right:    [][]int{[]int{1, 3, 5}, []int{1, 4, 5}, []int{4, 1, 1}},
			expected: true,
		},
		{
			title:    "4 elements (single, equals, different order)",
			left:     [][]int{[]int{2}, []int{5}, []int{1}, []int{3}},
			right:    [][]int{[]int{1}, []int{2}, []int{3}, []int{5}},
			expected: true,
		},
	}
	for _, testRun := range testRuns {
		t.Log("Running", testRun.title)
		computed := HasSameValuesDisregardingOrder(testRun.left, testRun.right)
		if computed != testRun.expected {
			t.Error("Computed was", computed, "and expected is", testRun.expected)
		}
	}
}

var slices = [][]int{
	[]int{1, 2, 3},
	[]int{1, 2, 3, 4},
	[]int{1, 2, 3, 4, 5},
	[]int{1, 2, 3, 4, 5, 6},
}

func BenchmarkGeneratePermutations(b *testing.B) {
	for _, slice := range slices {
		b.Run(fmt.Sprintf("slice.length:%v", len(slice)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GeneratePermutations(slice)
			}
		})
	}
}
