package cpu

import (
	"strconv"
	"sync"
)

// CPU сбор статистики по средней загрузке CPU.
type CPU struct {
	mu      *sync.Mutex
	bufSize int         // bufSize - максимальная глубина данных.
	index   int         // index - текущая позиция для новых значений.
	stats   [][]float64 // stats - слайс для хранения статистики.
}

const (
	userMode = 0 // Процент загрузки CPU для пользовательских процессов.
	sysMode  = 1 // Процент загрузки CPU для системных процессов.
	idle     = 2 // Процент не используемой загрузки CPU.
	total    = 3
)

func NewCPU(bufSize int) *CPU {
	st := make([][]float64, total)
	st[userMode] = make([]float64, bufSize)
	st[sysMode] = make([]float64, bufSize)
	st[idle] = make([]float64, bufSize)

	return &CPU{
		mu:      &sync.Mutex{},
		bufSize: bufSize,
		index:   0,
		stats:   st,
	}
}

func (cp *CPU) readValue(indexBuf, indexPos int) float64 {
	var value float64
	laIndex := cp.index - 1 //nolint:ifshort
	if laIndex < 0 {
		laIndex = cp.bufSize - 1
	}
	index := (indexPos + laIndex) % cp.bufSize
	value = cp.stats[indexBuf][index]
	return value
}

// GetValue подсчет среднего значения за последние m секунд.
func (cp *CPU) GetValue(m int) []string {
	str := make([]string, total)
	avg := []float64{0.0, 0.0, 0.0}

	cp.mu.Lock()
	defer cp.mu.Unlock()

	for y := 0; y < total; y++ {
		for i := 0; i < m; i++ {
			avg[y] += cp.readValue(y, i)
		}
		avg[y] /= float64(m)
		str[y] = strconv.FormatFloat(avg[y], 'f', 1, 64)
	}
	return str
}

// ShiftIndex вычисляет новое значение индекса в кольцевом буфере.
func (cp *CPU) ShiftIndex() {
	cp.index++
	if cp.index >= cp.bufSize {
		cp.index = 0
	}
}

// WriteValue записывает новое значение в кольцевой буфер.
func (cp *CPU) WriteValue(out []float64) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	result := out
	for i, v := range result {
		cp.stats[i][cp.index] = v
	}
}
