package ping

import (
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	type args struct {
		IPAddr string
		maxrtt time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantRtt float64
		wantErr bool
	}{
		{
			name: "TestPing",
			args: args{
				IPAddr: "127.0.0.1",
				maxrtt: 30 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRtt, err := Ping(tt.args.IPAddr, tt.args.maxrtt)
			t.Logf("Ping() = %v, want %v err:%v", gotRtt, tt.wantRtt, err)
		})
	}
}