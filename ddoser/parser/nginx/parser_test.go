package nginx_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/komandakycto/ddoser/ddoser/entities"
	"github.com/komandakycto/ddoser/ddoser/parser/nginx"
)

func TestParseDefaultFormat(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	parser := nginx.NewParser(true, "")
	logLine := `142.116.105.111 - - [18/Jul/2023:13:44:19 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/106.0"`

	entry, err := parser.Parse(logLine)
	r.NoError(err, "Failed to parse log line")

	expectedTime, _ := time.Parse(nginx.DefaultTimeLayout, "18/Jul/2023:13:44:19 +0000")
	expectedEntry := &entities.LogEntry{
		IPAddress:    "142.116.105.111",
		RequestedURL: "/",
		Timestamp:    expectedTime,
		UserAgent:    "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/106.0",
	}

	r.Equal(expectedEntry, entry, "Parsed entry mismatch")
}

func TestParseJsonFormat(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	parser := nginx.NewParser(false, time.RFC3339)
	logLine := `{"time":"2023-07-24T07:48:27+00:00","ts":"1690184907.681","ip":"43.249.187.131","x_real_ip":"","method":"GET","scheme":"https","domain":"test.com","uri":"/v2/login","args":"","status":"200","bytes_sent":"117327","request_length":"81","referer":"https://google.com","user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36","request_time":"0.001","upstream_response_time": "0.001","x_forwarded_for":"","hostname":"nginx-14","request_id":"a1f523901934ba7d051625b105b52604"}`

	entry, err := parser.Parse(logLine)
	r.NoError(err, "Failed to parse log line")

	expectedTime, _ := time.Parse(time.RFC3339, "2023-07-24T07:48:27+00:00")
	expectedEntry := &entities.LogEntry{
		IPAddress:    "43.249.187.131",
		RequestedURL: "/v2/login",
		Timestamp:    expectedTime,
		UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}

	r.Equal(expectedEntry, entry, "Parsed entry mismatch")
}

func TestParseInvalidDefaultLine(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	parser := nginx.NewParser(true, "")

	invalidLine := "invalid log line"
	_, err := parser.Parse(invalidLine)
	r.Error(err, "Expected an error for invalid log line")
}

func TestParseInvalidJsonLine(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	parser := nginx.NewParser(false, time.RFC3339)

	invalidLine := `{"time":"invalid timestamp","ip":"43.249.187.131","uri":"/v2/login","user_agent":"Mozilla/5.0"}`
	_, err := parser.Parse(invalidLine)
	r.Error(err, "Expected an error for invalid json timestamp")
}

func TestParseInvalidTimestampLayout(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	parser := nginx.NewParser(false, "invalid_layout")

	jsonLogLine := `{"time":"2023-07-24T07:48:27+00:00","ip":"43.249.187.131","uri":"/v2/login","user_agent":"Mozilla/5.0"}`
	_, err := parser.Parse(jsonLogLine)
	r.Error(err, "Expected an error for invalid timestamp layout")
}
