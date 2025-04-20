package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// ByteCounter counts bytes written to it.
type ByteCounter struct {
	n int64
}

func (bc *ByteCounter) Write(p []byte) (int, error) {
	n := len(p)
	bc.n += int64(n)
	return n, nil
}

// count the words in each line and the total number of lines
func (bc *ByteCounter) CountLines(r *bufio.Reader) (int, int, error) {
	lineCount := 0
	wordCount := 0
	for {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return 0, 0, err
		}
		if len(line) == 0 && err == io.EOF {
			break
		}
		lineCount++
		scanner := bufio.NewScanner(bytes.NewReader(line))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			wordCount++
		}
		if scanErr := scanner.Err(); scanErr != nil {
			return 0, 0, scanErr
		}
		if err == io.EOF {
			break
		}
	}
	return lineCount, wordCount, nil
}

func main() {
	bc := &ByteCounter{}
	n, err := bc.Write([]byte("Hello, World!"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Bytes written:", n)
	fmt.Println("Total bytes:", bc.n)
	// Example usage of CountLines
	input := "Hello, World!\nThis is a test.\nAnother line."
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))
	lines, words, err := bc.CountLines(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Lines: %d, Words: %d\n", lines, words)
}
