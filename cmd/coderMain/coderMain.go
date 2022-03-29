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
	var readerA *reader.Reader
	var writerA *writer.Writer
	fmt.Println(os.Args)
	if len(os.Args) != 3 {
		readerA = reader.Reader_createReader("../data/input/test")
		writerA = writer.Writer_createReader("../data/output/test")
	} else {
		readerA = reader.Reader_createReader(os.Args[1])
		writerA = writer.Writer_createReader(os.Args[2])
	}

	coder := coder.Coder_createCoder(readerA, writerA)

	coder.Coder_run()

	writerA.CloseFile()
}
