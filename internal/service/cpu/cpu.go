package cpu

import (
	"sync"

	"github.com/Haba1234/sysmon/internal/sysmon"
)

const (
	userMode = 0 // Процент загрузки CPU для пользовательских процессов.
	sysMode  = 1 // Процент загрузки CPU для системных процессов.
	idle     = 2 // Процент не используемой загрузки CPU.
	total    = 3
)

func NewCPU(bufSize int, f sysmon.DataRequestFunc) *sysmon.DataStore {
	st := make([][]float64, total)
	st[userMode] = make([]float64, bufSize)
	st[sysMode] = make([]float64, bufSize)
	st[idle] = make([]float64, bufSize)

	return &sysmon.DataStore{
		Mu:          &sync.Mutex{},
		BufSize:     bufSize,
		Stats:       st,
		Total:       total,
		FuncDataReq: f,
	}
}
