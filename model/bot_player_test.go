package model

import (
	"testing"
)

var cards []Card = []Card{Card{1, GOLD, 1}, Card{2, SWORD, 1}, Card{3, CLUB, 1}, Card{4, CUP, 1}}

func BenchmarkDetermineIsGoldenSevenIsUsed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetermineIsGoldenSevenIsUsed(cards)
	}
}
