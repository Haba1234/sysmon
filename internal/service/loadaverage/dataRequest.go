package loadaverage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func (la *LoadAverage) DataRequest() ([]float64, error) {
	const countData = 3 // Кол-во ожидаемых данных по средней загрузке (1 мин, 5 мин, 15 мин).
	val := []float64{0.0, 0.0, 0.0}
	var raw string
	var err error

	switch runtime.GOOS {
	case "linux":
		raw, err = readProcFile("/proc/loadavg")
		if err != nil {
			return nil, err
		}
	case "darwin":
		top := exec.Command("sysctl", "-n", "vm.loadavg")
		b, err := top.CombinedOutput()
		if err != nil {
			return nil, err
		}
		raw = strings.Trim(string(b), "{ }")
		log.Println("MAC OS.AVG:", string(b))
		log.Println("MAC OS.AVG:", raw)
	default:
		return nil, errors.New("command 'load average' not supported operating system")
	}

	n, err := fmt.Sscanf(raw, "%f %f %f",
		&val[0], &val[1], &val[2])
	if err != nil {
		return nil, err
	}
	if n < countData {
		return nil, errors.New("data 'load average' not fully read")
	}
	return val, nil
}

func readProcFile(str string) (string, error) {
	raw, err := ioutil.ReadFile(str)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
