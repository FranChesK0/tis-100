package emu

type NodeList struct {
	Node *Node
	Next *NodeList
}

func Append(list *NodeList, n *Node) *NodeList
