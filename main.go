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

func readFile(filePath string, callback func([]byte)) ([]string, error) {
	path, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer path.Close()

	scanner := bufio.NewScanner(path)

	for scanner.Scan() {
		callback(scanner.Bytes())
	}

	return nil, nil
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
		fmt.Println("Missing arguements, pattern string and target file both are required!")
		fmt.Println(format)
		os.Exit(0)
	}

	if len(args) > 2 {
		fmt.Println("Only two arguements are required!")
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

	_, readLineErr := readFile(filePath, func(line []byte) {
		positions := expression.FindAllIndex(line, -1)
		occurences := len(positions)
		if occurences > 0 {
			intervals := getIntervals(positions)
			highLigtedLine := applyColour(line, intervals)
			fmt.Println(highLigtedLine)
		}
	})

	if readLineErr != nil {
		log.Fatal(readLineErr)
		fmt.Println("Error in reading file!")
		os.Exit(0)
	}
}
