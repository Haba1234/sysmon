package loadaverage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataRequest(t *testing.T) {
	t.Run("test func DataRequest 'AVG'", func(t *testing.T) {
		la := NewLoadAverage(5)
		result, err := la.DataRequest()
		require.NoError(t, err)
		for _, val := range result {
			require.GreaterOrEqual(t, val, 0.0)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := readProcFile("test")
		require.NotNil(t, err)
	})
}
