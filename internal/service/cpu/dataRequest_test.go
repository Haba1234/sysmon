package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataRequest(t *testing.T) {
	t.Run("test func DataRequest", func(t *testing.T) {
		cpu := NewCPU(5)
		result, err := cpu.DataRequest()
		require.NoError(t, err)
		for _, val := range result {
			require.GreaterOrEqual(t, val, 0.0)
		}
	})
}
