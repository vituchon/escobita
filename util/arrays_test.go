package util

import (
	"fmt"
	"testing"
)

/*
user@computer:<escobita root path>$ go test -v ./util/ -bench=. -run=none -benchtime=10s
goos: linux
goarch: amd64
pkg: github.com/vituchon/escobita/util
cpu: Intel(R) Core(TM) i7-3630QM CPU @ 2.40GHz
BenchmarkGeneratePermutations_v1a
BenchmarkGeneratePermutations_v1a-8   	  100096	    118079 ns/op
BenchmarkGeneratePermutations_v1b
BenchmarkGeneratePermutations_v1b-8   	  104904	    123525 ns/op
BenchmarkGeneratePermutations_v2
BenchmarkGeneratePermutations_v2-8    	   94162	    128963 ns/op
BenchmarkGeneratePermutations_v3
BenchmarkGeneratePermutations_v3-8    	  101770	    113785 ns/op <-- FASTER
PASS
ok  	github.com/vituchon/escobita/util	53.371s

go test -v ./util/ -bench=. -run=none -benchtime=10s
goos: linux
goarch: amd64
pkg: github.com/vituchon/escobita/util
cpu: Intel(R) Core(TM) i7-3630QM CPU @ 2.40GHz
BenchmarkGeneratePermutations_v1a
BenchmarkGeneratePermutations_v1a/slice.length:3
BenchmarkGeneratePermutations_v1a/slice.length:3-8         	 3393565	      3685 ns/op
BenchmarkGeneratePermutations_v1a/slice.length:4
BenchmarkGeneratePermutations_v1a/slice.length:4-8         	  576556	     19936 ns/op
BenchmarkGeneratePermutations_v1a/slice.length:5
BenchmarkGeneratePermutations_v1a/slice.length:5-8         	   98236	    120050 ns/op
BenchmarkGeneratePermutations_v1a/slice.length:6
BenchmarkGeneratePermutations_v1a/slice.length:6-8         	   14029	    847624 ns/op
BenchmarkGeneratePermutations_v1b
BenchmarkGeneratePermutations_v1b/slice.length:3
BenchmarkGeneratePermutations_v1b/slice.length:3-8         	 3205861	      3714 ns/op
BenchmarkGeneratePermutations_v1b/slice.length:4
BenchmarkGeneratePermutations_v1b/slice.length:4-8         	  555519	     20484 ns/op
BenchmarkGeneratePermutations_v1b/slice.length:5
BenchmarkGeneratePermutations_v1b/slice.length:5-8         	  101011	    123917 ns/op
BenchmarkGeneratePermutations_v1b/slice.length:6
BenchmarkGeneratePermutations_v1b/slice.length:6-8         	   13896	    866823 ns/op
BenchmarkGeneratePermutations_v2
BenchmarkGeneratePermutations_v2/slice.length:3
BenchmarkGeneratePermutations_v2/slice.length:3-8          	 3149324	      3742 ns/op
BenchmarkGeneratePermutations_v2/slice.length:4
BenchmarkGeneratePermutations_v2/slice.length:4-8          	  529614	     20521 ns/op
BenchmarkGeneratePermutations_v2/slice.length:5
BenchmarkGeneratePermutations_v2/slice.length:5-8          	   90957	    127230 ns/op
BenchmarkGeneratePermutations_v2/slice.length:6
BenchmarkGeneratePermutations_v2/slice.length:6-8          	   13502	    872656 ns/op
BenchmarkGeneratePermutations_v3
BenchmarkGeneratePermutations_v3/slice.length:3
BenchmarkGeneratePermutations_v3/slice.length:3-8          	 3432536	      3468 ns/op
BenchmarkGeneratePermutations_v3/slice.length:4
BenchmarkGeneratePermutations_v3/slice.length:4-8          	  632379	     18908 ns/op
BenchmarkGeneratePermutations_v3/slice.length:5
BenchmarkGeneratePermutations_v3/slice.length:5-8          	  102729	    115519 ns/op
BenchmarkGeneratePermutations_v3/slice.length:6
BenchmarkGeneratePermutations_v3/slice.length:6-8          	   14456	    798634 ns/op

*/
var slice []int = []int{1, 2, 3, 4, 5}

var slices = [][]int{
	[]int{1, 2, 3},
	[]int{1, 2, 3, 4},
	[]int{1, 2, 3, 4, 5},
	[]int{1, 2, 3, 4, 5, 6},
}

func BenchmarkGeneratePermutations_v1a(b *testing.B) {
	for _, slice := range slices {
		b.Run(fmt.Sprintf("slice.length:%v", len(slice)), func(b *testing.B) { // taking advice from  https://blog.logrocket.com/benchmarking-golang-improve-function-performance#Benchmarking with various inputs
			for i := 0; i < b.N; i++ {
				GeneratePermutations_v1a(slice)
			}
		})
	}
}

func BenchmarkGeneratePermutations_v1b(b *testing.B) {
	for _, slice := range slices {
		b.Run(fmt.Sprintf("slice.length:%v", len(slice)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GeneratePermutations_v1b(slice)
			}
		})
	}
}

func BenchmarkGeneratePermutations_v2(b *testing.B) {
	for _, slice := range slices {
		b.Run(fmt.Sprintf("slice.length:%v", len(slice)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GeneratePermutations_v2(slice)
			}
		})
	}
}

func BenchmarkGeneratePermutations_v3(b *testing.B) {
	for _, slice := range slices {
		b.Run(fmt.Sprintf("slice.length:%v", len(slice)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GeneratePermutations_v3(slice)
			}
		})
	}
}
