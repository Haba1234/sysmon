package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Haba1234/sysmon/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	t.Run("test service", func(t *testing.T) {
		tFile, err := ioutil.TempFile("/tmp", "test")
		if err != nil {
			panic(err)
		}
		defer os.Remove(tFile.Name())

		logg, err := logger.New("INFO", tFile.Name())
		require.NoError(t, err)
		settings := Collection{
			LoadAverageEnabled: true,
			CPUEnabled:         true,
			BufSize:            10,
		}

		col := NewCollector(logg, settings)
		col.getDataStatistics()
		col.ReadStats(1)
		require.Equal(t, 3, len(col.StatsData.La))
		require.Equal(t, 3, len(col.StatsData.CPU))
	})
}
