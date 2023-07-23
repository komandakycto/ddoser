package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

// readEndBytes reads the last n lines from the log file.
func readEndBytes(logPath string, bytesToRead int64, log *logrus.Entry) ([]string, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.WithError(err).Errorf("Error closing file %s", logPath)
		}
	}(file)

	// Get the file size.
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Create a buffer to hold the read bytes.
	buffer := make([]byte, bytesToRead)

	// Read the last n lines from the file.
	lines := make([]string, 0, 1000)

	// Calculate the offset from the end of the file.
	offset := fileSize - bytesToRead
	if offset < 0 {
		offset = 0
	}

	// Calculate the number of bytes to read from the current offset.
	readed, err := file.ReadAt(buffer, offset)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	log.Infof("Readed %d bytes from file %s", readed, logPath)

	// start iterate over bytes from offset to the end of the file
	// and find the new line characters
	newLineByte := int64(0) // first byte in buffer
	for i := int64(0); i < bytesToRead; i++ {
		// build string until new line character
		str := string(buffer[newLineByte:i])
		if buffer[i] == '\n' || i == bytesToRead-1 {
			if i == bytesToRead-1 {
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

// splitSliceIntoGroups splits the input slice into n groups with k elements in each group.
func splitSliceIntoGroups(input []string, linesInGroup int) [][]string {
	numElements := len(input)
	numGroups := (numElements + linesInGroup - 1) / linesInGroup // Ceiling division to handle leftover elements.

	groups := make([][]string, numGroups)
	for i := 0; i < numGroups; i++ {
		start := i * linesInGroup
		end := (i + 1) * linesInGroup
		if end > numElements {
			end = numElements
		}
		groups[i] = input[start:end]
	}

	return groups
}
