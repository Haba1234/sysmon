package loadaverage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadAverage(t *testing.T) { //nolint:tparallel
	tests := []struct {
		name        string
		inData      []float64
		expectedErr error
		output      []string
	}{
		{
			name:        "Test 1",
			inData:      []float64{1.00, 0.58, 0.44},
			expectedErr: nil,
			output:      []string{"1.00", "0.58", "0.44"},
		},
		{
			name:        "Test 2",
			inData:      []float64{1.00, 0.58, 0.44},
			expectedErr: nil,
			output:      []string{"1.00", "0.58", "0.44"},
		},
		{
			name:        "Test 3",
			inData:      []float64{1.00, 0.58, 0.44},
			expectedErr: nil,
			output:      []string{"1.00", "0.58", "0.44"},
		},
		{
			name:        "Test 4",
			inData:      []float64{1.00, 0.58, 0.44},
			expectedErr: nil,
			output:      []string{"1.00", "0.58", "0.44"},
		},
		{
			name:        "Test 5",
			inData:      []float64{1.00, 0.58, 0.44},
			expectedErr: nil,
			output:      []string{"1.00", "0.58", "0.44"},
		},
		{
			name:        "Test 6",
			inData:      []float64{1.00, 0.58, 0.44},
			expectedErr: nil,
			output:      []string{"1.00", "0.58", "0.44"},
		},
	}

	la := NewLoadAverage(5)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			la.mu.Lock()
			la.WriteValue(tt.inData)
			la.ShiftIndex()
			la.mu.Unlock()
			out := la.GetValue(1)
			require.Equal(t, tt.output, out, "test failed")
			_ = tt
		})
	}
}
