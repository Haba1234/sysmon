// +build linux ios windows

package loadaverage

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
)

func runCMD() (string, error) {
	if runtime.GOOS != "windows" {
		/*grep := exec.Command("grep", "average")
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
		return string(b), err*/
		var info syscall.Sysinfo_t
		err := syscall.Sysinfo(&info)
		if err != nil {
			return "", err
		}
		const shift = 16
		str := fmt.Sprint(float64(info.Loads[0])/float64(1<<shift), " ")
		str += fmt.Sprint(float64(info.Loads[1])/float64(1<<shift), " ")
		str += fmt.Sprint(float64(info.Loads[2]) / float64(1<<shift))

		return str, nil
	}
	return "", errors.New("command 'load average' not supported operating system")
}
