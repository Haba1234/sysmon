package cpu

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// DataRequest с помощью exec читает данные по использованию CPU через команду top,
// парсит и выбирает только необходимые параметры.
func DataRequest() ([]float64, error) {
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
		return nil, fmt.Errorf("cannot execute command 'exec': %w", err)
	}
	b, err := grep.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot execute command 'exec': %w", err)
	}

	raw = strings.ReplaceAll(string(b), ", ", " ")
	raw = strings.ReplaceAll(raw, "%", "")
	raw = strings.TrimPrefix(raw, "CPU usage: ")
	n, err := fmt.Sscanf(raw, "%f user %f sys %f idle",
		&val[0], &val[1], &val[2])
	if err != nil {
		return nil, fmt.Errorf("cannot parse 'CPU average': %w", err)
	}
	if n < countData {
		return nil, errors.New("data 'CPU average' not fully read")
	}
	return []float64{val[0], val[1], val[2]}, nil
}
