package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	t.Run("FileExists", func(t *testing.T) {
		filePath := "test/testfile.txt"
		expectedLines := []string{"line 1", "line 2", "line 3"}

		lines, err := readFile(filePath, func(line []byte) {})
		assert.NoError(t, err)
		assert.Equal(t, expectedLines, lines)
	})

	t.Run("FileDoesNotExist", func(t *testing.T) {
		filePath := "nonexistentfile.txt"

		_, err := readFile(filePath, func(line []byte) {})
		assert.Error(t, err)
	})
}

func TestGetIntervals(t *testing.T) {
	positions := [][]int{{1, 4}, {6, 9}, {11, 14}}
	expectedIntervals := []int{1, 2, 3, 6, 7, 8, 11, 12, 13}

	intervals := getIntervals(positions)
	assert.Equal(t, expectedIntervals, intervals)
}

func TestIntervalContainsPositions(t *testing.T) {
	interval := []int{1, 2, 3, 6, 7, 8, 11, 12, 13}

	assert.True(t, intervalContainsPositions(interval, 2))
	assert.False(t, intervalContainsPositions(interval, 5))
}

func TestApplyColour(t *testing.T) {
	line := []byte("This is a test line")
	intervals := []int{0, 5, 8, 14}

	expectedOutput := "\x1b[31mT\x1b[0mhis \x1b[31mis a \x1b[0mtest \x1b[31ml\x1b[0mine"
	output := applyColour(line, intervals)

	assert.Equal(t, expectedOutput, output)
}

func TestMain(t *testing.T) {
	// Set up testing environment
	os.Args = []string{"cmd", "pattern", "test/testfile.txt"}

	// Capture standard output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	expectedOutput := "\x1b[31mline 1\x1b[0m\n\x1b[31mline 2\x1b[0m\n\x1b[31mline 3\x1b[0m\n"
	assert.Equal(t, expectedOutput, output)
}
