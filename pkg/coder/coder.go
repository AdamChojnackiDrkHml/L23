package coder

import (
	"container/heap"
	"fmt"
	"l23/pkg/node"
	"l23/pkg/reader"
	"l23/pkg/writer"
	"strconv"
	"strings"
)

type Coder struct {
	reader         *reader.Reader
	writer         *writer.Writer
	huffTree       *node.Node
	codeMap        map[byte]string
	probs          []float64
	counterSymbols []int64
	currentPatch   []byte
	lastPatch      bool
	w              string
	buffer         []string
	bytesBuffer    []byte
}

func Coder_createCoder(reader *reader.Reader, writer *writer.Writer) *Coder {
	coder := &Coder{reader: reader,
		writer:         writer,
		probs:          make([]float64, 256),
		counterSymbols: make([]int64, 256),
		currentPatch:   make([]byte, 0),
		lastPatch:      false,
		buffer:         make([]string, 0),
		bytesBuffer:    make([]byte, 0)}

	for i := range coder.counterSymbols {
		coder.counterSymbols[i] = 0
	}

	return coder
}

type ProbsHeap []*node.Node

func (h ProbsHeap) Len() int {
	return len(h)
}

func (h ProbsHeap) Less(i, j int) bool {
	return h[i].Probs < h[j].Probs
}

func (h ProbsHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *ProbsHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*node.Node))
}

func (h *ProbsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	*h = old[0 : n-1]
	return node
}

func (coder *Coder) calcProbs() {

	allSymbolsCounter := coder.reader.Counter
	probsss := 0.0
	for i := 0; i < len(coder.counterSymbols); i++ {
		coder.probs[i] = float64(coder.counterSymbols[i]) / float64(allSymbolsCounter)
		probsss += coder.probs[i]
	}
	fmt.Println(probsss)
	coder.w = strconv.Itoa(allSymbolsCounter) + " "
	coder.writeCode()
	coder.w = ""
}

func (coder *Coder) getData() {
	coder.currentPatch = coder.reader.Reader_readDataPatch()
	coder.lastPatch = !coder.reader.IsReading
}

func (coder *Coder) updateCounter() {
	for _, n := range coder.currentPatch {
		coder.counterSymbols[n]++
	}
}

func (coder *Coder) writeCode() {
	coder.writer.Writer_writeToFile(coder.w)
	coder.w = ""
}

func (coder *Coder) code() {
	reader.Reader_resetFile(&coder.reader)
	coder.lastPatch = !coder.lastPatch
	for !coder.lastPatch {
		coder.getData()

		for _, n := range coder.currentPatch {
			coder.addToBuffer(strings.Split(coder.codeMap[n], ""))
		}
	}

	if len(coder.buffer) != 0 {
		padding := 8 - len(coder.buffer)

		for ; padding > 0; padding-- {
			coder.buffer = append(coder.buffer, "0")
		}
		coder.passToStringByteFromBufferBits()

	}

	coder.w = string(coder.bytesBuffer)
	coder.writeCode()

}

func (coder *Coder) addToBuffer(bits []string) {

	for len(bits) > 0 {
		howManyToAdd := 8 - len(coder.buffer)

		if howManyToAdd > len(bits) {
			howManyToAdd = len(bits)
		}
		coder.buffer = append(coder.buffer, bits[:howManyToAdd]...)
		if len(coder.buffer) == 8 {
			coder.passToStringByteFromBufferBits()
		}
		bits = bits[howManyToAdd:]
	}
}

func (coder *Coder) passToStringByteFromBufferBits() {
	var acc byte
	acc = 0

	for _, n := range coder.buffer {
		acc *= 2
		val, _ := strconv.Atoi(n)
		acc += byte(val)
	}

	coder.bytesBuffer = append(coder.bytesBuffer, acc)

	if len(coder.bytesBuffer) == 256 {
		coder.w = string(coder.bytesBuffer)
		coder.writeCode()
		coder.bytesBuffer = make([]byte, 0)
	}
	coder.buffer = make([]string, 0)
}

func (coder *Coder) buildTree() {
	probsHeap := make(ProbsHeap, 0)

	for i, n := range coder.probs {
		probsHeap = append(probsHeap, node.Node_createNode(i, n))
	}

	heap.Init(&probsHeap)

	for len(probsHeap) != 1 {
		n1 := heap.Pop(&probsHeap).(*node.Node)
		n2 := heap.Pop(&probsHeap).(*node.Node)
		if n1.Probs == 0.25 {
			fmt.Println("a")
		}
		if n2.Probs == 0.25 {
			fmt.Println("b")
		}
		newNode := node.Node_joinNodes(n1, n2)

		heap.Push(&probsHeap, newNode)
		heap.Init(&probsHeap)

	}
	coder.huffTree = probsHeap[0]
	coder.codeMap = node.Node_createCodes(coder.huffTree)
	coder.w = node.Node_toString(coder.huffTree)

	coder.writeCode()
	fmt.Println("DONE")
}

func (coder *Coder) Coder_run() {
	for !coder.lastPatch {
		coder.getData()
		coder.updateCounter()
	}
	coder.calcProbs()
	coder.buildTree()
	coder.code()
	coder.writeCode()
}
