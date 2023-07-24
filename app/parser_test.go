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

func Test_parseJson(t *testing.T) {

	logTime, _ := time.Parse(time.RFC3339, "2023-07-24T10:30:25+00:00")

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
			name: "Test parseJson",
			args: args{
				logLine: `{"time":"2023-07-24T10:30:25+00:00","ts":"1690194625.050","ip":"49.229.253.154","x_real_ip":"","method":"GET","scheme":"https","domain":"deplab.g2afse.com","uri":"/click?pid=1770&offer_id=68&ref_id=bb56e5a7ln","args":"pid=1770&offer_id=68&ref_id=bb56e5a7ln","status":"302","bytes_sent":"319","request_length":"420","referer":"","user_agent":"Mozilla/5.0 (Linux; Android 13; V2109 Build/TP1A.220624.014) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.5735.196 Mobile Safari/537.36","request_time":"0.008","upstream_response_time": "0.009","x_forwarded_for":"","hostname":"nginx-26","request_id":"fd33a0cff7d48fc8f2ccaac10cad3fba"}`,
			},
			want: &LogEntry{
				IPAddress:    "49.229.253.154",
				RequestedURL: "/click?pid=1770&offer_id=68&ref_id=bb56e5a7ln",
				Timestamp:    logTime,
				UserAgent:    "Mozilla/5.0 (Linux; Android 13; V2109 Build/TP1A.220624.014) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.5735.196 Mobile Safari/537.36",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJson(tt.args.logLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}
