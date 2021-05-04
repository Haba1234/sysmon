package loadaverage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdLinux(t *testing.T) {
	t.Run("test func runCMD", func(t *testing.T) {
		str, err := runCMD()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(str), 0)
	})
}
