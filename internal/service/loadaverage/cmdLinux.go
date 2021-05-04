package loadaverage

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

func runCMD() (string, error) {
	/*var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return "", err
	}
	log.Println("load1: ", float64(info.Loads[0])/float64(1<<16))
	log.Println("load2: ", float64(info.Loads[1])/float64(1<<16))
	log.Println("load3: ", float64(info.Loads[2])/float64(1<<16))*/

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
