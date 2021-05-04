// +build linux darwin ios

package loadaverage

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

func runCMD() (string, error) {
	if runtime.GOOS != "windows" {
		grep := exec.Command("grep", "average")
		top := exec.Command("top", "-bn1")
		pipe, _ := top.StdoutPipe()
		defer pipe.Close()

		grep.Stdin = pipe
		err := top.Start()
		if err != nil {
			return "", err
		}
		b, err := grep.CombinedOutput()
		fmt.Println(string(b))
		return string(b), err
	}
	return "", errors.New("command 'load average' not supported operating system")
}
