package loadaverage

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// DataRequest читает данные load average через команду sysctl.
func DataRequest() ([]float64, error) {
	const countData = 3 // Кол-во ожидаемых данных по средней загрузке (1 мин, 5 мин, 15 мин).
	var raw string
	val := []float64{0.0, 0.0, 0.0}

	top := exec.Command("sysctl", "-n", "vm.loadavg")
	b, err := top.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cannot execute command 'exec': 'sysctl -n vm.loadavg': %w", err)
	}
	raw = string(b[2:])

	fray := strings.Split(raw, " ")
	for i := 0; i < countData; i++ {
		val[i], _ = strconv.ParseFloat(fray[i], 64)
	}
	return val, nil
}
