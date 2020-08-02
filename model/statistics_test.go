package model

import (
	"testing"
)

func TestSeventiesCalculation(t *testing.T) {
	testRuns := []struct {
		title         string
		cards         []Card
		expectedValue int
	}{
		{
			title:         "1 + 1 + 1 + 1 = 8",
			cards:         []Card{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, SWORD, 1}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 8,
		},
	}
	for _, testRun := range testRuns {
		t.Logf("Running unit test: %s ", testRun.title)
		computedValue := calculateSeventiesScore(testRun.cards...)
		if computedValue != testRun.expectedValue {
			t.Errorf("aggregated card values differs! Expected is %d and computed value is %d", testRun.expectedValue, computedValue)
		}
	}
}
