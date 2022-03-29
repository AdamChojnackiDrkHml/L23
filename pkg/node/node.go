package node

import "l23/pkg/reader"

type Node struct {
	Left    *Node
	Right   *Node
	Probs   float64
	Name    byte
	IsInner bool
}

func Node_createNode(name int, probs float64) *Node {
	return &Node{
		Left:    nil,
		Right:   nil,
		Probs:   probs,
		Name:    byte(name),
		IsInner: false}
}

func Node_createInnerNode(probs float64) *Node {
	return &Node{
		Left:    nil,
		Right:   nil,
		Probs:   probs,
		Name:    0,
		IsInner: true}
}

func Node_recreateNode(name int) *Node {
	return &Node{
		Left:    nil,
		Right:   nil,
		Probs:   0.0,
		Name:    byte(name),
		IsInner: false}
}

func Node_joinNodes(n1, n2 *Node) *Node {
	newRoot := Node_createInnerNode(n1.Probs + n2.Probs)

	newRoot.Left = n1
	newRoot.Right = n2

	return newRoot
}

func postOrderTraverse(node *Node) string {

	if node == nil {
		return ""
	}

	dictionary := ""

	dictionary += postOrderTraverse(node.Left)
	dictionary += postOrderTraverse(node.Right)

	if node.IsInner {
		dictionary += "0"
	} else {
		dictionary += "1" + string(rune(node.Name))
	}

	return dictionary
}

func Node_toString(node *Node) string {

	return postOrderTraverse(node) + "0"

}

func Node_createCodes(node *Node) map[byte]string {
	ma := make(map[byte]string)

	s := ""

	stepDown(node, s, &ma)

	return ma
}

func stepDown(node *Node, s string, ma *map[byte]string) {

	if !node.IsInner {
		(*ma)[node.Name] = s
		return
	}

	if node.Left != nil {
		stepDown(node.Left, s+"0", ma)
	}

	if node.Right != nil {
		stepDown(node.Right, s+"1", ma)
	}
}

func Node_verySadAndCoupledFunctionToRecreateTree(reader *reader.Reader) *Node {
	stack := make([]*Node, 0)
	for {
		myByte := reader.Reader_readByte()
		if !reader.IsReading {
			return nil
		}
		if myByte == byte('1') {
			myByte = reader.Reader_readByte()
			stack = append(stack, Node_recreateNode(int(myByte)))
			continue
		}

		if myByte == byte('0') {
			if len(stack) == 1 {
				break
			}
			n := len(stack) - 1
			right := stack[n]
			left := stack[n-1]
			newInnerNode := Node_joinNodes(left, right)
			stack = stack[:n-1]
			stack = append(stack, newInnerNode)
		}

	}

	return stack[0]
}
