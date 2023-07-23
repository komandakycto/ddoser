package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseLogLine(t *testing.T) {
	logTime, _ := time.Parse(TimeLayout, "18/Jul/2023:13:44:18 +0000")

	type args struct {
		logLine string
	}
	tests := []struct {
		name    string
		args    args
		want    *LogEntry
		wantErr bool
	}{
		{
			name: "Test parseLogLine",
			args: args{
				logLine: `192.236.195.162 - - [18/Jul/2023:13:44:18 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			},
			want: &LogEntry{
				IPAddress:    "192.236.195.162",
				RequestedURL: "/",
				Timestamp:    logTime,
				UserAgent:    "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogLine(tt.args.logLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLogLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLogLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}
