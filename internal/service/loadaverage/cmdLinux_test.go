package loadaverage

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdLinux(t *testing.T) {
	t.Run("test func runCMD", func(t *testing.T) {
		result, err := runCMD()

		if runtime.GOOS != "darwin" {
			require.NoError(t, err)
			for _, val := range result {
				require.GreaterOrEqual(t, val, 0.0)
			}
			return
		}
		require.NotNil(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := readProcFile("test")
		require.NotNil(t, err)
	})
}
