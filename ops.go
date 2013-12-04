package ops

import (
	"github.com/etsuo/log"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(cmd string) error {
	return runCmd(cmd, false)
}

func RunCriticalCommand(cmd string) error {
	return runCmd(cmd, true)
}

func runCmd(cmd string, critical bool) error {
	parts := strings.Split(cmd, " ")

	c := exec.Command(parts[0])
	c.Args = parts

	log.Send.Infof("Executing OS Command: %s", cmd)

	c.Stdout = os.Stdout
	c.Stderr = os.Stdout

	err := c.Run()
	if err != nil {
		if critical {
			log.Send.Fatalf(err.Error())
		} else {
			log.Send.Warningf(err.Error())
		}
	}

	return err
}

// FsObjExists returns true or false based on
// whether or not a file system object (file / directory)
// exists
func FsObjExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
