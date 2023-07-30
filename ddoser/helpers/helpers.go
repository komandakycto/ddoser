package helpers

import "net"

// SplitSliceIntoGroups splits the input slice into n groups with k elements in each group.
func SplitSliceIntoGroups(input []string, linesInGroup int) [][]string {
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

// IsIPv4 checks if the input string is a valid IPv4 address.
func IsIPv4(address string) bool {
	ip := net.ParseIP(address)
	return ip != nil && ip.To4() != nil
}
