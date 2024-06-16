package command

import (
	"os"
	"os/exec"
	"time"
)

const waitDelay = 1000 * time.Hour

func SetCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = waitDelay
}
