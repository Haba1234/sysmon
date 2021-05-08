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
		raw = string(b[2:])
	default:
		return nil, errors.New("command 'load average' not supported operating system")
	}
	log.Println("MAC OS.AVG:", raw)
	fray := strings.Split(raw, " ")
	/*_, err = fmt.Sscanf(raw, "{ %f %f %f }",
		&val[0], &val[1], &val[2])
	if err != nil {
		return nil, err
	}*/
	log.Println(fray)
	for i := 0; i < countData; i++ {
		val[i], _ = strconv.ParseFloat(fray[i], 64)
	}
	log.Println(val)
	return val, nil
}

func readProcFile(str string) (string, error) {
	raw, err := ioutil.ReadFile(str)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
