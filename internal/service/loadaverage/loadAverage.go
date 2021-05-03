package loadaverage

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Haba1234/sysmon/internal/logger"
)

// LoadAverage сбор статистики по средней загрузке.
type LoadAverage struct {
	mu      *sync.Mutex
	logg    *logger.Logger
	bufSize int         // bufSize - максимальная глубина данных.
	index   int         // index - текущая позиция для новых значений.
	stats   [][]float64 // stats - слайс для хранения статистики.
}

const (
	oneMin     = 0
	fiveMin    = 1
	fifteenMin = 2
	sizeLA     = 3
)

func NewLoadAverage(logg *logger.Logger, bufSize int) *LoadAverage {
	st := make([][]float64, sizeLA)
	st[oneMin] = make([]float64, bufSize)
	st[fiveMin] = make([]float64, bufSize)
	st[fifteenMin] = make([]float64, bufSize)

	la := &LoadAverage{
		mu:      &sync.Mutex{},
		logg:    logg,
		bufSize: bufSize,
		index:   0,
		stats:   st,
	}
	return la
}

// Start запуск сервиса.
func (la *LoadAverage) Start(ctx context.Context) error {
	la.logg.Info("service 'load average' starting...")

	go func() {
		out, err := runCMD()
		if err != nil {
			la.logg.Error(fmt.Sprintf("runCMD() failed with %v", err))
		}

		if err := la.addNewValue(out); err != nil {
			la.logg.Error(fmt.Sprintf("addNewValue() failed add new value: %v", err))
		}

		ticker := time.NewTicker(time.Second) // Сбор данных раз в секунду.
		for {
			select {
			case <-ctx.Done():
				la.logg.Info("'load average' stopped")
				ticker.Stop()
				return
			case <-ticker.C:
				out, err := runCMD()
				if err != nil {
					la.logg.Error(fmt.Sprintf("runCMD() failed with %v", err))
				}

				if err := la.addNewValue(out); err != nil {
					la.logg.Error(fmt.Sprintf("addNewValue() failed add new value: %v", err))
				}
			}
		}
	}()

	return nil
}

// Stop останов сервиса.
func (la *LoadAverage) Stop(ctx context.Context) error {
	la.logg.Info("service 'load average'  stopping...")
	return nil
}

var re = regexp.MustCompile(`: [0-9,. ]+`)

// Функция парсит среднюю загрузку за минуту, 5 минут, 15 минут и сохраняет
// в соответствующие слайсы.
func (la *LoadAverage) addNewValue(out string) error {
	result := re.FindString(out)
	if len(result) == 0 {
		return errors.New("addNewValue() wrong string")
	}

	result = strings.Trim(result, ": \n")
	result = strings.ReplaceAll(result, ", ", " ")
	arr := strings.SplitN(result, " ", 4)

	la.mu.Lock()
	defer la.mu.Unlock()

	for i, s := range arr {
		val, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		if err != nil {
			return err
		}
		la.stats[i][la.index] = val
	}
	la.shiftIndex()
	return nil
}

// Функция вычисляет новое значение индекса.
func (la *LoadAverage) shiftIndex() {
	la.index++
	if la.index >= la.bufSize {
		la.index = 0
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

// AverageValue подсчет среднего значения за последние m секунд.
func (la *LoadAverage) AverageValue(m int) ([]string, error) {
	if m <= 0 || m > la.bufSize {
		return nil, errors.New("parameter N (period) has an invalid value")
	}
	str := make([]string, sizeLA)
	avg := []float64{0.0, 0.0, 0.0}

	la.mu.Lock()
	defer la.mu.Unlock()

	for y := 0; y < sizeLA; y++ {
		for i := 0; i < m; i++ {
			avg[y] += la.readValue(y, i)
		}
		avg[y] /= float64(m)
		str[y] = strconv.FormatFloat(avg[y], 'f', 2, 64)
	}
	return str, nil
}
