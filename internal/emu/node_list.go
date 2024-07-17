package emu

type NodeList struct {
	Node *Node
	Next *NodeList
}

func Append(list *NodeList, n *Node) *NodeList {
	tail := &NodeList{
		Node: n,
		Next: nil,
	}
	if list == nil {
		return tail
	}

	head := list
	for head.Next != nil {
		head = head.Next
	}
	head.Next = tail
	return list
}
