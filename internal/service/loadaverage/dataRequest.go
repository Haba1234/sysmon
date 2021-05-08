package loadaverage

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strconv"
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
		raw = string(b[2:])
	default:
		return nil, errors.New("command 'load average' not supported operating system")
	}

	fray := strings.Split(raw, " ")
	for i := 0; i < countData; i++ {
		val[i], _ = strconv.ParseFloat(fray[i], 64)
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
