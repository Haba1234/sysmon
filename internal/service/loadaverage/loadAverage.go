package loadaverage

import (
	"strconv"
	"sync"
)

// LoadAverage сбор статистики по средней загрузке.
type LoadAverage struct {
	mu      *sync.Mutex
	bufSize int         // bufSize - максимальная глубина данных.
	index   int         // index - текущая позиция для новых значений.
	stats   [][]float64 // stats - слайс для хранения статистики.
}

const (
	oneMin     = 0
	fiveMin    = 1
	fifteenMin = 2
	total      = 3
)

func NewLoadAverage(bufSize int) *LoadAverage {
	st := make([][]float64, total)
	st[oneMin] = make([]float64, bufSize)
	st[fiveMin] = make([]float64, bufSize)
	st[fifteenMin] = make([]float64, bufSize)

	return &LoadAverage{
		mu:      &sync.Mutex{},
		bufSize: bufSize,
		index:   0,
		stats:   st,
	}
}

func (la *LoadAverage) readValue(indexBuf, indexPos int) float64 {
	var value float64
	laIndex := la.index - 1 //nolint:ifshort
	if laIndex < 0 {
		laIndex = la.bufSize - 1
	}
	index := (indexPos + laIndex) % la.bufSize
	value = la.stats[indexBuf][index]
	return value
}

// GetValue подсчет среднего значения за последние m секунд.
func (la *LoadAverage) GetValue(m int) []string {
	str := make([]string, total)
	avg := []float64{0.0, 0.0, 0.0}

	la.mu.Lock()
	defer la.mu.Unlock()

	for y := 0; y < total; y++ {
		for i := 0; i < m; i++ {
			avg[y] += la.readValue(y, i)
		}
		avg[y] /= float64(m)
		str[y] = strconv.FormatFloat(avg[y], 'f', 2, 64)
	}
	return str
}

// ShiftIndex вычисляет новое значение индекса в кольцевом буфере.
func (la *LoadAverage) ShiftIndex() {
	la.index++
	if la.index >= la.bufSize {
		la.index = 0
	}
}

// WriteValue записывает новое значение в кольцевой буфер.
func (la *LoadAverage) WriteValue(out []float64) {
	la.mu.Lock()
	defer la.mu.Unlock()

	result := out
	for i, v := range result {
		la.stats[i][la.index] = v
	}
}
