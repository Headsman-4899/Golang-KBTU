package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rot *rot13Reader) Read(bytes []byte) (n int, err error) {
	n, err = rot.r.Read(bytes)
	for i := 0; i < len(bytes); i++ {
		if (bytes[i] >= 'A' && bytes[i] < 'N') || (bytes[i] >= 'a' && bytes[i] < 'n') {
			bytes[i] += 13
		} else if (bytes[i] > 'M' && bytes[i] <= 'Z') || (bytes[i] > 'm' && bytes[i] <= 'z') {
			bytes[i] -= 13
		}
	}
	return
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
