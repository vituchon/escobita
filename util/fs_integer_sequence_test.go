package util

import (
	"os"
	"testing"
)

func TestFsIntegerSequenceWorks(t *testing.T) {
	tests := []struct {
		name         string
		initialValue int
		expected     int
	}{
		{
			name:         "Archivo no existe, inicializa en 0",
			initialValue: 0,
			expected:     1,
		},
		{
			name:         "Secuencia inicia en 10",
			initialValue: 10,
			expected:     11,
		},
		{
			name:         "Secuencia después de un incremento",
			initialValue: 5,
			expected:     6,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpFile := "test_sequence.txt"
			defer os.Remove(tmpFile)

			var seq IntegerSequence = NewIntegerSequence(tmpFile, test.initialValue, 1)
			generated, err := seq.GetNext()
			if err != nil {
				t.Errorf("unexpected error: '%v'", err)
			}
			if generated != test.expected {
				t.Errorf("GetNext(). generated = %v, expected %v", generated, test.expected)
			}
		})
	}
}
