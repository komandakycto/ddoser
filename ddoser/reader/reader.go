package reader

import (
	"fmt"
	"os"
)

type LogReader struct {
	filepath        string
	bytes2read      int64
	averageLogBytes int64
}

// NewLogReader creates a new LogReader instance.
func NewLogReader(filepath string, bytes2read, averageLogBytes int64) *LogReader {
	return &LogReader{
		filepath:        filepath,
		bytes2read:      bytes2read,
		averageLogBytes: averageLogBytes,
	}
}

// ReadLastBytes reads the last n bytes from the end of the log file.
func (lr *LogReader) ReadLastBytes() ([]string, error) {
	file, err := os.Open(lr.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Get the file size.
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get the file size: %v", err)
	}
	fileSize := fileInfo.Size()

	// Create a buffer to hold the read bytes.
	buffer := make([]byte, lr.bytes2read)

	// will parse bytes to string and split by new line character.
	lines := make([]string, 0, lr.bytes2read/lr.averageLogBytes)

	// Calculate the offset from the end of the file.
	offset := fileSize - lr.bytes2read
	if offset < 0 {
		offset = 0
	}

	// Calculate the number of bytes to read from the current offset.
	_, err = file.ReadAt(buffer, offset)
	if err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("failed to read from the file: %v", err)
	}

	// start iterate over bytes from offset to the end of the file
	// and find the new line characters
	newLineByte := int64(0) // first byte in buffer
	for i := int64(0); i < lr.bytes2read; i++ {
		// build string until new line character
		str := string(buffer[newLineByte:i])
		if buffer[i] == '\n' || i == lr.bytes2read-1 {
			if i == lr.bytes2read-1 {
				// add last byte
				str = string(buffer[newLineByte : i+1])
			}
			lines = append(lines, str)
			newLineByte = i + 1 // next byte after new line character
		}
	}

	// drop first line because it may be incomplete
	if offset != 0 && len(lines) > 0 {
		lines = lines[1:]
	}

	return lines, nil
}
