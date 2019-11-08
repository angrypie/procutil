package procutil

import (
	"os"
	"syscall"
	"time"
)

//DefaultTerminateDuaration is timout duration before sending SIGKILL.
const DefaultTerminateDuaration = time.Second * 2

//Terminate stops process with SIGTERM or with SIGKILL while timeout exceeded.
func Terminate(proc *os.Process, timeout ...time.Duration) error {
	duration := DefaultTerminateDuaration
	if len(timeout) > 0 {
		duration = timeout[0]
	}
	done := make(chan error, 1)
	go func() {
		done <- proc.Signal(syscall.SIGTERM)
	}()
	select {
	case <-time.After(duration):
		return proc.Kill()
	case err := <-done:
		//Kill if sigterm returns error or process just still alive
		if err != nil || proc.Signal(syscall.Signal(0)) == nil {
			return proc.Kill()
		}
	}

	return nil
}
