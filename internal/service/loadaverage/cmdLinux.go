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

func runCMD() ([]float64, error) {
	log.Println("GOOS: ", runtime.GOOS)
	const countData = 3 // Кол-во ожидаемых данных по средней загрузке (1 мин, 5 мин, 15 мин).

	if runtime.GOOS == "linux" {
		raw, err := readProcFile("/proc/loadavg")
		if err != nil {
			return nil, err
		}

		result := strings.ReplaceAll(raw, ",", ".")
		val := []float64{0.0, 0.0, 0.0}
		n, err := fmt.Sscanf(result, "%f %f %f",
			&val[0], &val[1], &val[2])
		if err != nil {
			return nil, err
		}
		if n < countData {
			return nil, errors.New("data 'load average' not fully read")
		}
		return val, nil
	}

	if runtime.GOOS == "darwin" {
		top := exec.Command("sysctl", "-n", "vm.loadavg")
		b, err := top.CombinedOutput()
		if err != nil {
			return nil, err
		}
		log.Println("MAC OS:", string(b))

		top = exec.Command("uptime")
		b, err = top.CombinedOutput()
		if err != nil {
			return nil, err
		}
		log.Println("MAC OS:", string(b))
		return nil, nil
	}
	return nil, errors.New("command 'load average' not supported operating system")
}

func readProcFile(str string) (string, error) {
	raw, err := ioutil.ReadFile(str)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
