package kingdom_algorithmia

import (
	"bufio"
	"fmt"
	"io"
)

type Quest01 struct{}

func (q Quest01) Solve(input io.Reader) (string, error) {
	scanner := bufio.NewScanner(input)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Your quest logic here
	result := fmt.Sprintf("Quest 1 processed %d lines", len(lines))
	return result, nil
}
