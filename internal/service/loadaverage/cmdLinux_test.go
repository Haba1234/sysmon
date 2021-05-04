package loadaverage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdLinux(t *testing.T) {
	t.Run("test func runCMD", func(t *testing.T) {
		result, err := runCMD()
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
