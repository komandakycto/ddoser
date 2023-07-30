package processor_test

import (
	"context"
	"github.com/komandakycto/ddoser/app/parser/nginx"
	"github.com/komandakycto/ddoser/app/processor"
	"github.com/sirupsen/logrus"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIPAnalysis_ProcessDefault(t *testing.T) {
	t.Parallel()

	log := logrus.New()
	l := logrus.NewEntry(log)

	parser := nginx.NewParser(true, nginx.DefaultTimeLayout)

	// Assume that these are valid log lines
	groups := [][]string{
		{
			`192.168.0.1 - - [18/Jul/2023:13:44:17 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.1 - - [18/Jul/2023:13:44:17 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.1 - - [18/Jul/2023:13:44:18 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.1 - - [18/Jul/2023:13:44:18 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.1 - - [18/Jul/2023:13:44:19 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.1 - - [18/Jul/2023:13:44:19 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.3 - - [18/Jul/2023:13:44:20 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.3 - - [18/Jul/2023:13:44:20 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
		},
		{
			`192.168.0.5 - - [18/Jul/2023:13:44:15 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.5 - - [18/Jul/2023:13:44:15 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.5 - - [18/Jul/2023:13:44:15 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.7 - - [18/Jul/2023:13:44:16 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
		},
		{
			`192.168.0.17 - - [18/Jul/2023:13:44:21 +0000] "GET /login/form HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.17 - - [18/Jul/2023:13:44:21 +0000] "GET /login/form HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.17 - - [18/Jul/2023:13:44:21 +0000] "GET /login/form HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
			`192.168.0.17 - - [18/Jul/2023:13:44:21 +0000] "GET /login/form HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0"`,
		},
	}

	// Successful processing of groups
	t.Run("Successful processing. threshold: 5, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(5, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.7": true, "192.168.0.1": true}, result)
	})

	t.Run("Successful processing. threshold: 2, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(2, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.1": true, "192.168.0.3": true, "192.168.0.5": true, "192.168.0.7": true, "192.168.0.17": true}, result)
	})

	t.Run("Successful processing. threshold: 10, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(10, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.7": true}, result)
	})

	t.Run("Successful processing. threshold: 20, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(20, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{}, result)
	})

	t.Run("Successful processing. threshold: 3, windows: 10, url: login", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(3, 10, "login", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.17": true}, result)
	})

	// test context cancel
	t.Run("Context cancel", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(5, 10, "", parser, false, l)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{}, result)
	})
}

func TestIPAnalysis_ProcessJson(t *testing.T) {
	t.Parallel()

	log := logrus.New()
	l := logrus.NewEntry(log)

	parser := nginx.NewParser(false, time.RFC3339)

	// Assume that these are valid log lines
	groups := [][]string{
		{
			`{"time":"2023-07-18T13:44:17+00:00","ts":"1690184907.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:17+00:00","ts":"1690184907.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:18+00:00","ts":"1690184908.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:18+00:00","ts":"1690184908.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:19+00:00","ts":"1690184909.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:19+00:00","ts":"1690184909.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"192.168.0.3","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"192.168.0.3","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
		},
		{
			`{"time":"2023-07-18T13:44:15+00:00","ts":"1690184905.681","ip":"192.168.0.5","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:15+00:00","ts":"1690184905.681","ip":"192.168.0.5","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:15+00:00","ts":"1690184905.681","ip":"192.168.0.5","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:16+00:00","ts":"1690184906.681","ip":"192.168.0.7","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
		},
		{
			`{"time":"2023-07-18T13:44:21+00:00","ts":"1690184901.681","ip":"192.168.0.17","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/login/form","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:21+00:00","ts":"1690184901.681","ip":"192.168.0.17","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/login/form","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:21+00:00","ts":"1690184901.681","ip":"192.168.0.17","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/login/form","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:21+00:00","ts":"1690184901.681","ip":"192.168.0.17","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/login/form","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
		},
	}

	// Successful processing of groups
	t.Run("Successful processing. threshold: 5, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(5, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.7": true, "192.168.0.1": true}, result)
	})

	t.Run("Successful processing. threshold: 2, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(2, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.1": true, "192.168.0.3": true, "192.168.0.5": true, "192.168.0.7": true, "192.168.0.17": true}, result)
	})

	t.Run("Successful processing. threshold: 10, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(10, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.7": true}, result)
	})

	t.Run("Successful processing. threshold: 20, windows: 10", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(20, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{}, result)
	})

	t.Run("Successful processing. threshold: 3, windows: 10, url: login", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(3, 10, "login", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.17": true}, result)
	})

	// test context cancel
	t.Run("Context cancel", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(5, 10, "", parser, false, l)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{}, result)
	})
}

func TestIPAnalysis_ProcessIPv6(t *testing.T) {
	t.Parallel()

	log := logrus.New()
	l := logrus.NewEntry(log)

	parser := nginx.NewParser(false, time.RFC3339)

	// Assume that these are valid log lines
	groups := [][]string{
		{
			`{"time":"2023-07-18T13:44:17+00:00","ts":"1690184907.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:17+00:00","ts":"1690184907.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:18+00:00","ts":"1690184908.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:18+00:00","ts":"1690184908.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:19+00:00","ts":"1690184909.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:19+00:00","ts":"1690184909.681","ip":"192.168.0.1","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"192.168.0.3","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"192.168.0.3","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"192.168.0.3","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"2001:0db8:85a3:0000:0000:8a2e:0370:7334","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"2001:0db8:85a3:0000:0000:8a2e:0370:7334","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"2001:0db8:85a3:0000:0000:8a2e:0370:7334","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
			`{"time":"2023-07-18T13:44:20+00:00","ts":"1690184910.681","ip":"2001:0db8:0000:0000:0000:0000:1428:57ab","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/","args":"","status":"499","bytes_sent":"0","request_length":"81","referer":"-","user_agent":"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0","request_time":"0.001","upstream_response_time":"0.001","x_forwarded_for":"","hostname":"","request_id":""}`,
		},
	}

	// Successful processing of groups
	t.Run("Successful processing. threshold: 3, windows: 10, ipv4", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(3, 10, "", parser, true, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.3": true, "192.168.0.1": true}, result)
	})

	t.Run("Successful processing. threshold: 3, windows: 10, all api", func(t *testing.T) {
		t.Parallel()

		analysis := processor.NewIPAnalysis(2, 10, "", parser, false, l)
		ctx := context.Background()

		result, err := analysis.Process(ctx, groups)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"192.168.0.3": true, "192.168.0.1": true, "2001:0db8:85a3:0000:0000:8a2e:0370:7334": true}, result)
	})
}
