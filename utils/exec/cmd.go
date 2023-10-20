package exec

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"syscall"
	"time"

	"github.com/wonderivan/logger"
)

// CMDResult ...
type CMDResult struct {
	output string
	err    error
}

// RunCmdWithTimeOut ...
func RunCmdWithTimeOut(command string, timeout int) (out string, code int, err error) {
	bufOut := new(bytes.Buffer)
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", command)
	done := make(chan error)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = bufOut
	cmd.Stderr = bufOut
	// 子进程也杀死
	err = cmd.Start()
	if err != nil {
		return
	}
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		logger.Error("cmd %s run %d sec，process:%d kill now", timeout, cmd.Process.Pid)
		go func() {
			<-done // 读出上面的goroutine数据，避免阻塞导致无法退出
		}()
		if err = cmd.Process.Kill(); err != nil {
			logger.Error("cmd %s  process:%d kill err:%s", command, cmd.Process.Pid, err)
		}
		code = cmd.ProcessState.ExitCode()
		out = bufOut.String()
		return
	case err = <-done:
		out = bufOut.String()
		code = cmd.ProcessState.ExitCode()
		return
	}
}

// RunCmdWithTimeOutContext ...手动触发退出
func RunCmdWithTimeOutContext(ctx context.Context, command string) (out string, code int, err error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	resultChan := make(chan CMDResult)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- CMDResult{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {

			err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			if err != nil {
				logger.Error("cmd %s  process:%d kill err:%s", command, cmd.Process.Pid, err)
			}
			out = "timeout killed"
			code = cmd.ProcessState.ExitCode()
			err = errors.New("timeout killed")
		}
		return
	case result := <-resultChan:
		return result.output, cmd.ProcessState.ExitCode(), result.err
	}
}

// timeout := time.Duration(taskReq.Timeout) * time.Second
// ctx, cancel := context.WithTimeout(context.Background(), timeout)
// defer cancel()
