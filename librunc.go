package librunc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/opencontainers/specs"
)

const defaultSpecName = "config.json"

var errorCodesMap map[int]string

type Container struct {
	directory string
}

// Takes in empty directory and json configuration
func New(dir string, spec specs.LinuxSpec) (c *Container, err error) {
	if _, err := exec.LookPath("runc"); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dir, 0655); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(&spec, "", "\t")
	if err != nil {
		return nil, err
	}

	configPath := path.Join(dir, defaultSpecName)
	if err := ioutil.WriteFile(configPath, data, 0666); err != nil {
		return nil, err
	}

	return &Container{dir}, nil
}

// Takes in directory. Tries to open configuration from that directory.
func NewFromDirectory(dir string) error {
	return fmt.Errorf("not-implemented") // todo
}

func (c *Container) Start(stdin io.Reader, stdout, stderr io.Writer) error {
	cmd, err := c.getCommand()
	if err != nil {
		return err
	}

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // detach from parent so library can go down
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// nothing interesting to do with this yet. possibly never
	// go func() {
	// 	cmd.Wait()
	// }()

	return nil
}

func (c *Container) Kill(signal string) error {
	cmd, err := c.getCommand()
	if err != nil {
		return err
	}

	cmd.Args = append(cmd.Args, "kill", signal)

	if err := cmd.Run(); err != nil {
		return mapError(err)
	}

	return nil
}

func (c *Container) getCommand() (*exec.Cmd, error) {
	runcPath, err := exec.LookPath("runc")
	if err != nil {
		return nil, err
	}

	return &exec.Cmd{
		Path: runcPath,
		Args: []string{"runc"},
		Dir:  c.directory,
	}, nil

}

func mapError(err error) error {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if procExit, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			if newerror, ok := errorCodesMap[procExit.ExitStatus()]; ok {
				return fmt.Errorf(newerror)
			}
		}
	}
	return err
}

// func (*Container) exec(cmd *exec.Cmd)

// func (*Container) Stats(interval int)
// func (*Container) NotifyOOM()
// func (*Container) Checkpoint()
// func (*Container) Restore()
