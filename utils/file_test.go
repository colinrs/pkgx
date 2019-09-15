package utils

import "testing"

func TestWriteMsgToFile(t *testing.T) {
	type args struct {
		filename string
		msg      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestWriteMsgToFile",
			args: args{
				filename: "test_data/test_data.log",
				msg:      "test_data/test_data.log",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteMsgToFile(tt.args.filename, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("WriteMsgToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
