package main

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_splitSliceIntoGroups(t *testing.T) {
	type args struct {
		input        []string
		linesInGroup int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Test splitSliceIntoGroups",
			args: args{
				input: []string{
					"line1",
					"line2",
					"line3",
					"line4",
					"line5",
					"line6",
					"line7",
					"line8",
					"line9",
					"line10",
				},
				linesInGroup: 2,
			},
			want: [][]string{
				{"line1", "line2"},
				{"line3", "line4"},
				{"line5", "line6"},
				{"line7", "line8"},
				{"line9", "line10"},
			},
		},
		{
			name: "Test splitSliceIntoGroups 3",
			args: args{
				input: []string{
					"line1",
					"line2",
					"line3",
				},
				linesInGroup: 2,
			},
			want: [][]string{
				{"line1", "line2"},
				{"line3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitSliceIntoGroups(tt.args.input, tt.args.linesInGroup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitSliceIntoGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readLastNLines(t *testing.T) {
	type args struct {
		logPath           string
		numberLinesToRead int64
		log               *logrus.Entry
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Test readEndBytes",
			args: args{
				logPath:           "./../data/test.txt",
				numberLinesToRead: 100,
				log:               logrus.NewEntry(logrus.New()),
			},
			want: []string{
				`line 20`,
				`line 21`,
				`line 22`,
				`line 23`,
				`line 24`,
				`line 25`,
				`line 26`,
				`line 27`,
				`line 28`,
				`line 29`,
				`line 30`,
				`line 31`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readEndBytes(tt.args.logPath, tt.args.numberLinesToRead, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("readEndBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readEndBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
