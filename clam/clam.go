package clam

import (
	"bufio"
	"os/exec"
)

func RunAndListen(cmd *exec.Cmd, fn func(s string)) error {
	var err error
	r, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			fn(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	return err
}

func Run(cmd *exec.Cmd) (string, error) {
	var (
		out []byte
		err error
	)
	if out, err = cmd.Output(); err != nil {
		return "", err
	}
	return string(out), err
}
