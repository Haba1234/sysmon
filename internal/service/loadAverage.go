package service

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Haba1234/sysmon/internal/logger"
)

type LoadAverage struct {
	mu         *sync.Mutex
	logg       *logger.Logger
	bufSize    int64
	oneMin     []float64
	fiveMin    []float64
	fifteenMin []float64
}

func NewLoadAverage(logg *logger.Logger, bufSize int64) *LoadAverage {
	la := &LoadAverage{
		mu:         &sync.Mutex{},
		logg:       logg,
		bufSize:    bufSize,
		oneMin:     make([]float64, bufSize),
		fiveMin:    make([]float64, bufSize),
		fifteenMin: make([]float64, bufSize),
	}

	return la
}

func (la *LoadAverage) Start(ctx context.Context) error {
	la.logg.Info("service 'load average' starting...")

	var args = []string{ //nolint:gofumpt
		"-c",
		"top -bn1 | fgrep 'average'", // | tail -2
	}

	go func() {
		if err := la.addNewValue(args); err != nil {
			la.logg.Error(fmt.Sprintf("addNewValue() failed add new value: %s", err))
		}

		ticker := time.NewTicker(time.Second) // Сбор данных раз в секунду.
	loop:
		for {
			select {
			case <-ctx.Done():
				la.logg.Info("'load average' stopped")
				ticker.Stop()
				break loop
			case <-ticker.C:
				if err := la.addNewValue(args); err != nil {
					la.logg.Error(fmt.Sprintf("addNewValue() failed add new value: %s", err))
				}
			}
		}
	}()

	return nil
}

func (la *LoadAverage) Stop(ctx context.Context) error {
	la.logg.Info("service 'load average'  stopping...")
	return nil
}

func runCMD(path string, args []string) (string, error) {
	cmd := exec.Command(path, args...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

var re = regexp.MustCompile(`: [0-9,. ]+`)

// Функция парсит среднюю загрузку за минуту, 5 минут, 15 минут и сохраняет
// в соответствующие слайсы.
func (la *LoadAverage) addNewValue(args []string) error {
	out, err := runCMD("bash", args)
	if err != nil {
		la.logg.Error(fmt.Sprintf("runCMD() failed with %s", err))
	}
	out = re.FindString(out)
	if len(out) == 0 {
		return errors.New("addNewValue() wrong string")
	}

	la.deleteLastValue()
	out = strings.Trim(out, ": \n")
	out = strings.ReplaceAll(out, ", ", " ")
	arr := strings.SplitN(out, " ", 4)
	la.mu.Lock()
	defer la.mu.Unlock()
	for i, s := range arr {
		val, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		if err != nil {
			return err
		}
		switch i {
		case 0:
			la.oneMin[0] = val
		case 1:
			la.fiveMin[0] = val
		case 2:
			la.fifteenMin[0] = val
		}
	}
	return nil
}

// Функция сдвигает значения в слайсах на одно значение назад,
// освобождая место под новое значение в конце слайса.
func (la *LoadAverage) deleteLastValue() {
	la.mu.Lock()
	defer la.mu.Unlock()
	for i := la.bufSize - 1; i > 0; i-- {
		la.oneMin[i] = la.oneMin[i-1]
		la.fiveMin[i] = la.fiveMin[i-1]
		la.fifteenMin[i] = la.fifteenMin[i-1]
	}
}

// Подсчет среднего значения за последние m секунд.
func (la *LoadAverage) AverageValue(m int64) ([]string, error) {
	if m <= 0 || m > la.bufSize {
		return nil, errors.New("parameter N (period) has an invalid value")
	}
	str := make([]string, 3)
	avg0 := 0.0
	avg1 := 0.0
	avg2 := 0.0

	for i := 0; i < int(m)-1; i++ {
		avg0 += la.oneMin[i]
		avg1 += la.fiveMin[i]
		avg2 += la.fifteenMin[i]
	}
	avg0 /= float64(m)
	avg1 /= float64(m)
	avg2 /= float64(m)

	str[0] = strconv.FormatFloat(avg0, 'f', 2, 64)
	str[1] = strconv.FormatFloat(avg1, 'f', 2, 64)
	str[2] = strconv.FormatFloat(avg2, 'f', 2, 64)

	return str, nil
}
