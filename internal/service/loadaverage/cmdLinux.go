package loadaverage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
)

func runCMD() (string, error) {
	log.Println("GOOS: ", runtime.GOOS)
	if runtime.GOOS != "ios" {
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
		return string(b), err

		*/

		raw, err := ioutil.ReadFile("/proc/loadavg")
		if err != nil {
			return "", err
		}
		str := []string{"", "", ""}
		fmt.Sscanf(string(raw), "%s %s %s",
			&str[0], &str[1], &str[2])

		return str[0] + " " + str[1] + " " + str[1], nil
	}
	return "", errors.New("command 'load average' not supported operating system")
}
