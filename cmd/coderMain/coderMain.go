package main

import (
	"fmt"
	"l23/pkg/coder"
	"l23/pkg/reader"
	"l23/pkg/writer"
	"os"
)

func main() {
	fmt.Println(os.Getwd())
	reader := reader.Reader_createReader("data/input/test")
	writer := writer.Writer_createReader("data/output/test")

	coder := coder.Coder_createCoder(reader, writer)

	coder.Coder_run()

	writer.CloseFile()
}
