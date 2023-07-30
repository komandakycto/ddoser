package reader_test

import (
	"github.com/komandakycto/ddoser/ddoser/reader"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogReader_ReadLastBytes(t *testing.T) {
	r := require.New(t)

	// Define the test file path.
	filepath := "./../../test_data/access.log"

	// Create a new LogReader instance.
	// Since the averageLogBytes is not known in this context, we set it to 100 for testing purposes.
	logReader := reader.NewLogReader(filepath, 10000, 100)

	// Read the last bytes from the log file.
	lines, err := logReader.ReadLastBytes()
	r.NoError(err, "Failed to read last bytes from log file")

	// Expected number of lines to be read.
	expectedNumLines := 63 // Let's assume we want to read the last 10 lines.

	// Verify the number of lines read from the log file.
	r.Equal(expectedNumLines, len(lines), "Number of lines read is not as expected")

	// Verify the content of the lines.
	expectedLastLine := `142.116.105.111 - - [18/Jul/2023:13:44:19 +0000] "GET / HTTP/2.0" 499 0 "-" "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:108.0) Gecko/20100101 Firefox/106.0"`
	r.Equal(expectedLastLine, lines[expectedNumLines-1], "Last line content is not as expected")
}
