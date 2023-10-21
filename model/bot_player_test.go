package model

import (
	"testing"
)

/**
on my machine tests favors the native implementation
escobita$ go test -v ./model/ -bench=. -run=none -benchtime=10s
goos: linux
goarch: amd64
pkg: github.com/vituchon/escobita/model
cpu: Intel(R) Core(TM) i7-3630QM CPU @ 2.40GHz
BenchmarkDetermineIsGoldenSevenIsUsedNative
BenchmarkDetermineIsGoldenSevenIsUsedNative-8   	1000000000	         7.274 ns/op
BenchmarkDetermineIsGoldenSevenIsUsed
BenchmarkDetermineIsGoldenSevenIsUsed-8         	652290052	        18.66 ns/op
PASS


I will continue with the native "using for loop" implementation as gitlab CI jobs gets broken because of
```
Run go build -v ./...
go: downloading github.com/gorilla/handlers v1.5.1
go: downloading github.com/gorilla/mux v1.8.0
go: downloading github.com/gorilla/securecookie v1.1.1
go: downloading github.com/gorilla/sessions v1.2.1
go: downloading github.com/gorilla/websocket v1.5.0
go: downloading github.com/lib/pq v1.10.9
go: downloading github.com/felixge/httpsnoop v1.0.1
Error: model/bot_player.go:5:2: no required module provides package golang.org/x/exp/slices; to add it:
	go get golang.org/x/exp/slices
Error: Process completed with exit code 1.
```

And also... this benchmarks shows that native implementation is twice fasters than the slices function
*/
var cards []Card = []Card{Card{1, GOLD, 1}, Card{2, SWORD, 1}, Card{3, CLUB, 1}, Card{4, CUP, 1}}

func BenchmarkDetermineIsGoldenSevenIsUsedNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetermineIsGoldenSevenIsUsedNative(cards)
	}
}

func BenchmarkDetermineIsGoldenSevenIsUsed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetermineIsGoldenSevenIsUsed(cards)
	}
}
