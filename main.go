package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

const format = `FORMAT: go run main.go [pattern] [file]`

func readFile(filePath string, callback func(int, []byte)) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		callback(lineNumber, scanner.Bytes())
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func getIntervals(positions [][]int) []int {
	var intervals []int

	for _, position := range positions {
		initialInterval := position[0]
		lastInterval := position[1]

		for initialInterval < lastInterval {
			intervals = append(intervals, initialInterval)
			initialInterval += 1
		}
	}

	return intervals
}

func intervalContainsPositions(interval []int, position int) bool {
	for _, index := range interval {
		if position == index {
			return true
		}
	}
	return false
}

func applyColour(line []byte, intervals []int) string {
	var modifiedLine []string

	for char, index := range string(line) {
		if intervalContainsPositions(intervals, char) {
			modifiedLine = append(modifiedLine, color.RedString(string(index)))
			continue
		}
		modifiedLine = append(modifiedLine, string(index))
	}
	return strings.Join(modifiedLine, "")
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Missing arguments, both pattern string and target file are required!")
		fmt.Println(format)
		os.Exit(0)
	}

	if len(args) > 2 {
		fmt.Println("Only two arguments are required!")
		fmt.Println(format)
		os.Exit(0)
	}

	pattern := args[0]
	filePath := args[1]

	if _, err := os.Stat(filePath); err != nil {
		fmt.Printf("%s file does not exist!", filePath)
		os.Exit(0)
	}

	expression := regexp.MustCompile(pattern)

	err := readFile(filePath, func(lineNumber int, line []byte) {
		positions := expression.FindAllIndex(line, -1)
		occurrences := len(positions)
		if occurrences > 0 {
			intervals := getIntervals(positions)
			highlightedLine := applyColour(line, intervals)
			fmt.Printf("%d: %s\n", lineNumber, highlightedLine)
		}
	})

	if err != nil {
		log.Fatal(err)
		fmt.Println("Error in reading file!")
		os.Exit(0)
	}
}
