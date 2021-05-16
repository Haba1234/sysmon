package loadaverage

import (
	"sync"

	"github.com/Haba1234/sysmon/internal/sysmon"
)

const (
	oneMin     = 0
	fiveMin    = 1
	fifteenMin = 2
	total      = 3
)

func NewLoadAverage(bufSize int, f sysmon.DataRequestFunc) *sysmon.DataStore {
	st := make([][]float64, total)
	st[oneMin] = make([]float64, bufSize)
	st[fiveMin] = make([]float64, bufSize)
	st[fifteenMin] = make([]float64, bufSize)

	return &sysmon.DataStore{
		Mu:          &sync.Mutex{},
		BufSize:     bufSize,
		Stats:       st,
		Total:       total,
		FuncDataReq: f,
	}
}
