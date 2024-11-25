package util

import (
	"os"
	"testing"
)

func TestFsIntegerSequenceWorks(t *testing.T) {
	tests := []struct {
		name         string
		initialValue int
		increment    int
		expected     int
	}{
		{
			name:         "Archivo no existe, inicializa en 0",
			initialValue: 0,
			increment:    1,
			expected:     1,
		},
		{
			name:         "Secuencia inicia en 10",
			initialValue: 10,
			increment:    1,
			expected:     11,
		},
		{
			name:         "Secuencia después de un incremento",
			initialValue: 5,
			increment:    2,
			expected:     7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpFile := "test_sequence.txt"
			defer os.Remove(tmpFile)

			var seq IntegerSequence = NewFsIntegerSequence(tmpFile, test.initialValue, test.increment)
			generated, err := seq.GetNext()
			if err != nil {
				t.Errorf("unexpected error: '%v'", err)
			}
			if generated != test.expected {
				t.Errorf("GetNext(). generated = %v, expected %v", generated, test.expected)
			}

			// checking a second call, shall generate the inmediate following according to increment
			generated, err = seq.GetNext()
			if err != nil {
				t.Errorf("unexpected error: '%v'", err)
			}
			if generated != test.expected+test.increment {
				t.Errorf("GetNext() + increment. generated = %v, expected %v", generated, test.expected+test.increment)
			}
		})
	}
}
