package loadaverage

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// DataRequest читает данные load average из файла.
func DataRequest() ([]float64, error) {
	const countData = 3 // Кол-во ожидаемых данных по средней загрузке (1 мин, 5 мин, 15 мин).
	var raw string
	val := []float64{0.0, 0.0, 0.0}

	b, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, fmt.Errorf("cannot read 'loadavg' file: %w", err)
	}
	raw = string(b)

	fray := strings.Split(raw, " ")
	for i := 0; i < countData; i++ {
		val[i], _ = strconv.ParseFloat(fray[i], 64)
	}
	return val, nil
}
