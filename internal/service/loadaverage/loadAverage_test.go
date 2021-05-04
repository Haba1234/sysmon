package loadaverage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Haba1234/sysmon/internal/logger"
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

	tFile, err := ioutil.TempFile("/tmp", "test")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tFile.Name())

	logg, err := logger.New("INFO", tFile.Name())
	require.NoError(t, err)
	la := NewLoadAverage(logg, 5)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			la.addNewValue(tt.inData)
			out, err := la.AverageValue(1)
			require.NoError(t, err)
			require.Equal(t, tt.output, out, "test failed")
			_ = tt
		})
	}

	t.Run("M > N", func(t *testing.T) {
		_, err := la.AverageValue(10)
		require.Error(t, err)
	})
}
