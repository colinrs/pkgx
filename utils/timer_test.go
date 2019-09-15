package utils

import (
	"testing"
	"time"
)

func TestParseTimeStr2Time(t *testing.T) {
	type args struct {
		timeStr string
	}
	tests := []struct {
		name    string
		args    args
		wantT   time.Time
		wantErr bool
	}{
		{
			name: "TestParseTimeStr2Time1",
			args: args{
				timeStr: "201909090000",
			},
		},
		{
			name: "TestParseTimeStr2Time2",
			args: args{
				timeStr: "1568305106",
			},
		},
		{
			name: "TestParseTimeStr2Time2",
			args: args{
				timeStr: "xxxx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := ParseTimeStr2Time(tt.args.timeStr)
			t.Logf("ParseTimeStr2Time() gotT %v", gotT.Format(FIVEMIN1))
			t.Logf("ParseTimeStr2Time() gotT %v", gotT.Unix())
			t.Logf("ParseTimeStr2Time() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
