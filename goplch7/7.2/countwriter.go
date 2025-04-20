package main

import (
	"bytes"
	"fmt"
	"io"
)

// CountingWriter wraps an io.Writer and counts the number of bytes written to it.
type counter struct {
	w io.Writer
	n int64
}

// Write method to count the number of bytes written
func (c *counter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n += int64(n)
	return n, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	c := &counter{w: w}
	return c, &c.n
}

func main() {

	// Example usage
	w := &bytes.Buffer{}
	cw, n := CountingWriter(w)
	cw.Write([]byte("Hello, World!"))
	fmt.Println(*n)         // Output: 13
	fmt.Println(w.String()) // Output: Hello, World!

}
