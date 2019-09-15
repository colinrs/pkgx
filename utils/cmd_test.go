package utils

import (
	"context"
	"testing"
	"time"
)

func TestRunCmdWithTimeOut(t *testing.T) {
	type args struct {
		cmd     string
		timeout int
	}
	tests := []struct {
		name     string
		args     args
		wantOut  string
		wantCode int
		wantErr  bool
	}{
		{
			name: "TestRunCmdWithTimeOut-py",
			args: args{
				cmd:     "python ./test_data/test.py",
				timeout: 30,
			},
		},
		{
			name: "TestRunCmdWithTimeOut-sh",
			args: args{
				cmd:     "sh ./test_data/test.sh",
				timeout: 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, gotCode, err := RunCmdWithTimeOut(tt.args.cmd, tt.args.timeout)
			t.Logf("RunCmdWithTimeOut() error = %v,gotOut = %v,gotCode = %v", err, gotOut, gotCode)
		})
	}
}

func TestRunCmdWithTimeOutContext(t *testing.T) {

	timeout := time.Duration(2) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	type args struct {
		ctx     context.Context
		command string
	}
	tests := []struct {
		name     string
		args     args
		wantOut  string
		wantCode int
		wantErr  bool
	}{
		{
			name: "TestRunCmdWithTimeOut-sh",
			args: args{
				ctx:     ctx,
				command: "sh ./test_data/test.sh 1",
			},
		},
		{
			name: "TestRunCmdWithTimeOut-sh",
			args: args{
				ctx:     ctx,
				command: "sh ./test_data/test.sh 10",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, gotCode, err := RunCmdWithTimeOutContext(tt.args.ctx, tt.args.command)
			t.Logf("RunCmdWithTimeOut() error = %v,gotOut = %v,gotCode = %v", err, gotOut, gotCode)
		})
	}
}
