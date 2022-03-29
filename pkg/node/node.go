package node

import "strconv"

type Node struct {
	left    *Node
	right   *Node
	Probs   float64
	name    byte
	isInner bool
}

func Node_createNode(name int, probs float64) *Node {
	return &Node{
		left:    nil,
		right:   nil,
		Probs:   probs,
		name:    byte(name),
		isInner: false}
}

func Node_createInnerNode(probs float64) *Node {
	return &Node{
		left:    nil,
		right:   nil,
		Probs:   probs,
		name:    0,
		isInner: true}
}

func Node_joinNodes(n1, n2 *Node) *Node {
	newRoot := Node_createInnerNode(n1.Probs + n2.Probs)

	newRoot.left = n1
	newRoot.right = n2

	return newRoot
}

func postOrderTraverse(node *Node) string {

	if node == nil {
		return ""
	}

	dictionary := ""

	dictionary += postOrderTraverse(node.left)
	dictionary += postOrderTraverse(node.right)

	if node.isInner {
		dictionary += "0"
	} else {
		dictionary += "1" + strconv.Itoa(int(node.name))
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

	if !node.isInner {
		(*ma)[node.name] = s
		return
	}

	if node.left != nil {
		stepDown(node.left, s+"0", ma)
	}

	if node.right != nil {
		stepDown(node.right, s+"1", ma)
	}
}
