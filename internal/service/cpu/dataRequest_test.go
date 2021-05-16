package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataRequest(t *testing.T) {
	t.Run("test func DataRequest 'CPU'", func(t *testing.T) {
		result, err := DataRequest()
		require.NoError(t, err)
		for _, val := range result {
			require.GreaterOrEqual(t, val, 0.0)
		}
	})
}
