package loadaverage

import (
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func (la *LoadAverage) DataRequest() ([]float64, error) {
	const countData = 3 // Кол-во ожидаемых данных по средней загрузке (1 мин, 5 мин, 15 мин).
	val := []float64{0.0, 0.0, 0.0}
	var raw string
	var fraw []string
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
		//raw = strings.Trim(string(b), "{ }")
		fraw = strings.Split(raw, " ")
		log.Println("MAC OS.AVG:", string(b))
		log.Println("MAC OS.AVG:", raw)
	default:
		return nil, errors.New("command 'load average' not supported operating system")
	}

	/*_, err = fmt.Sscanf(raw, "{ %f %f %f }",
		&val[0], &val[1], &val[2])
	if err != nil {
		return nil, err
	}*/
	val[0], _ = strconv.ParseFloat(fraw[1], 64)
	val[1], _ = strconv.ParseFloat(fraw[2], 64)
	val[2], _ = strconv.ParseFloat(fraw[3], 64)
	return val, nil
}

func readProcFile(str string) (string, error) {
	raw, err := ioutil.ReadFile(str)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
