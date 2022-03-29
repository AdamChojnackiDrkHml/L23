package decoder

import (
	"fmt"
	"l23/pkg/node"
	"l23/pkg/reader"
	"l23/pkg/writer"
	"strconv"
	"strings"
)

type Decoder struct {
	reader       *reader.Reader
	writer       *writer.Writer
	huffTree     *node.Node
	buffer       []byte
	w            string
	numOfSymbols int64
}

func Decoder_createDecoder(reader *reader.Reader, writer *writer.Writer) *Decoder {
	decoder := &Decoder{
		reader: reader,
		writer: writer,
		buffer: make([]byte, 0)}

	return decoder
}

func (decoder *Decoder) decode() {

	bitsPointer := 0
	myByte := decoder.reader.Reader_readByte()
	bits := splitByteToBits(myByte)
	for decoder.reader.IsReading && decoder.numOfSymbols > 0 {

		currNode := decoder.huffTree

		//search for sign
		for currNode.IsInner {
			if bits[bitsPointer] == "1" {
				currNode = currNode.Right
				bitsPointer++
			} else if bits[bitsPointer] == "0" {
				currNode = currNode.Left
				bitsPointer++
			}
			if bitsPointer == 8 {
				myByte = decoder.reader.Reader_readByte()
				if !decoder.reader.IsReading {
					break
				}
				bits = splitByteToBits(myByte)
				bitsPointer = 0
			}
		}

		decoder.addToBuffer(currNode.Name)
		decoder.numOfSymbols--

	}

	decoder.w = string(decoder.buffer)
	decoder.writeCode()
	decoder.buffer = make([]byte, 0)
}

func (decoder *Decoder) addToBuffer(myByte byte) {
	decoder.buffer = append(decoder.buffer, myByte)

	if len(decoder.buffer) == 256 {
		decoder.w = string(decoder.buffer)
		decoder.writeCode()
		decoder.buffer = make([]byte, 0)
	}
}

func (decoder *Decoder) writeCode() {
	decoder.writer.Writer_writeToFile(decoder.w)
	decoder.w = ""
}

func (decoder *Decoder) getTree() {
	decoder.huffTree = node.Node_verySadAndCoupledFunctionToRecreateTree(decoder.reader)
}

func splitByteToBits(aByte byte) []string {
	return strings.Split(fmt.Sprintf("%08b", aByte), "")
}

func (decoder *Decoder) getNumOfSymbols() {

	s := decoder.reader.Reader_getFirstWord()
	if !decoder.reader.IsReading {
		return
	}
	decoder.numOfSymbols, _ = strconv.ParseInt(s, 10, 64)
}

func (decoder *Decoder) Decoder_run() {
	decoder.getNumOfSymbols()
	decoder.getTree()
	decoder.decode()
}
