package loadaverage

import (
	"log"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdLinux(t *testing.T) {
	t.Run("test func runCMD", func(t *testing.T) {
		str, err := runCMD()
		log.Println("str:", str)
		log.Println("err:", err)
		if runtime.GOOS != "darwin" {
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(str), 0)
			return
		}
		require.NotNil(t, err)
	})
}
