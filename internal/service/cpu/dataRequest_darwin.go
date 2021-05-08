// +build darwin ios

package cpu

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// DataRequest с помощью exec читает данные по использованию CPU через команду top,
// парсит и выбирает только необходимые параметры.
func (cp *CPU) DataRequest() ([]float64, error) {
	const countData = 3 // Кол-во ожидаемых данных по проценту использования CPU (user, sys, idle).
	val := []float64{0.0, 0.0, 0.0}
	var raw string
	var err error

	grep := exec.Command("grep", "CPU")
	top := exec.Command("top", "-bn1")
	pipe, _ := top.StdoutPipe()
	defer pipe.Close()
	grep.Stdin = pipe
	err = top.Start()
	if err != nil {
		return nil, err
	}
	b, err := grep.Output()
	if err != nil {
		return nil, err
	}
	log.Println(b)
	raw = strings.ReplaceAll(string(b), ", ", " ")
	raw = strings.ReplaceAll(raw, "ni,", "ni ")
	raw = strings.ReplaceAll(raw, ",", ".")
	raw = strings.TrimPrefix(raw, "%Cpu(s): ")

	n, err := fmt.Sscanf(raw, "%f us %f sy %f ni %f id",
		&val[0], &val[1], &val[2], &val[3])
	if err != nil {
		return nil, err
	}
	if n < countData {
		return nil, errors.New("data 'load average' not fully read")
	}
	return []float64{val[0], val[1], val[3]}, nil
}
