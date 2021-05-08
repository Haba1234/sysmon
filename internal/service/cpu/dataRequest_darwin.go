// +build darwin

package cpu

import (
	"errors"
	"fmt"
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
	top := exec.Command("top", "-l 1")
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
	fmt.Println("MAC OS.CPU:", string(b))
	raw = strings.ReplaceAll(string(b), ", ", " ")
	fmt.Println("MAC OS.CPU1:", raw)
	raw = strings.ReplaceAll(raw, "%", "")
	fmt.Println("MAC OS.CPU2:", raw)
	raw = strings.TrimPrefix(raw, "CPU usage: ")
	fmt.Println("MAC OS.CPU3:", raw)
	n, err := fmt.Sscanf(raw, "%f user %f sys %f idle",
		&val[0], &val[1], &val[2])
	fmt.Println("MAC OS.CPU4:", val)
	if err != nil {
		return nil, err
	}
	if n < countData {
		return nil, errors.New("data 'load average' not fully read")
	}
	return []float64{val[0], val[1], val[3]}, nil
}
