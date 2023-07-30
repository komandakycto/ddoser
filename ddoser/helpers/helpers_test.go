package helpers_test

import (
	"github.com/komandakycto/ddoser/ddoser/helpers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitSliceIntoGroups(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	// Test case with a number of elements divisible by linesInGroup.
	input1 := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	linesInGroup1 := 4
	expectedOutput1 := [][]string{
		{"a", "b", "c", "d"},
		{"e", "f", "g", "h"},
		{"i", "j", "k", "l"},
	}
	r.Equal(expectedOutput1, helpers.SplitSliceIntoGroups(input1, linesInGroup1), "Grouping failed")

	// Test case with a number of elements not divisible by linesInGroup.
	input2 := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	linesInGroup2 := 5
	expectedOutput2 := [][]string{
		{"a", "b", "c", "d", "e"},
		{"f", "g", "h", "i", "j"},
		{"k", "l"},
	}
	r.Equal(expectedOutput2, helpers.SplitSliceIntoGroups(input2, linesInGroup2), "Grouping failed")
}

func TestIsIPv4(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	// Test case with a valid IPv4 address.
	validIP := "192.168.0.1"
	r.True(helpers.IsIPv4(validIP), "Expected true for valid IPv4 address")

	// Test case with an invalid IPv4 address.
	invalidIP := "192.168.0.300"
	r.False(helpers.IsIPv4(invalidIP), "Expected false for invalid IPv4 address")
}
