package main

import "golang.org/x/tour/reader"

type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.
func (mr MyReader) Read(bytes []byte) (int, error) {
	for index, _ := range bytes {
		bytes[index] = 65
	}
	return len(bytes), nil
}

func main() {
	reader.Validate(MyReader{})
}
