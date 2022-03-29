package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Reader struct {
	path               string
	file               *os.File
	PatchSize          int64
	ReadSymbolsCounter int
	IsReading          bool
	scanner            *bufio.Scanner
	Counter            int
}

func Print(a string) {
	fmt.Println(a)
}

func (reader *Reader) openFile() {
	file, err := os.Open(reader.path)

	if err != nil {
		panic(err)
	}

	reader.file = file

}

func Reader_createReader(path string) *Reader {
	reader := &Reader{path: path, PatchSize: 256, IsReading: true}

	reader.openFile()
	reader.Counter = 0
	reader.scanner = bufio.NewScanner(reader.file)
	return reader
}

func (reader *Reader) Reader_readDataPatch() []byte {
	currPatch := make([]byte, reader.PatchSize)
	control, err := reader.file.Read(currPatch)

	if err == io.EOF {
		reader.closeFile()
		reader.IsReading = false

	}

	reader.ReadSymbolsCounter = control
	reader.Counter += control
	fmt.Println(reader.Counter)
	return currPatch[:control]
}

func (reader *Reader) Reader_readLine() []string {
	reader.IsReading = reader.scanner.Scan()

	return strings.Split(reader.scanner.Text(), " ")
}

func (reader *Reader) closeFile() {
	reader.file.Close()
}

func Reader_resetFile(reader **Reader) {
	(*reader) = Reader_createReader((*reader).path)
}
