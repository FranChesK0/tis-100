package emu

import (
	"errors"
	"strings"

	"github.com/FranChesK0/tis-100/internal/constants"
	"github.com/FranChesK0/tis-100/internal/types"
)

type Program struct {
	Nodes       []*Node
	NodeList    *NodeList
	ActiveNodes *NodeList
	Outputs     []*Output
}

func NewProgram() *Program {
	nodes := make([]*Node, 0, constants.NodesNumber)
	var n *Node
	for i := range constants.NodesNumber {
		n = NewNode()
		n.Index = uint8(i)
		nodes = append(nodes, n)
	}
	p := &Program{
		Nodes:   nodes,
		Outputs: make([]*Output, 0),
	}

	for i := range p.Nodes {
		if i != 8 && i != 9 && i != 10 && i != 11 {
			p.Nodes[i].Ports[DOWN] = p.Nodes[i+4]
		}
		if i != 0 && i != 1 && i != 2 && i != 3 {
			p.Nodes[i].Ports[UP] = p.Nodes[i-4]
		}
		if i != 3 && i != 7 && i != 11 {
			p.Nodes[i].Ports[RIGHT] = p.Nodes[i+1]
		}
		if i != 0 && i != 4 && i != 8 {
			p.Nodes[i].Ports[LEFT] = p.Nodes[i-1]
		}
	}

	return p
}

func (p *Program) Tick() (bool, error) {
	allBlocked := true
	var err error
	for list := p.ActiveNodes; list != nil; list = list.Next {
		if err = list.Node.Tick(); err != nil {
			return false, err
		}
		allBlocked = allBlocked && list.Node.Blocked
	}
	return allBlocked, nil
}

func (p *Program) LoadStreams(streams []types.Stream) error {
	for _, stream := range streams {
		switch stream.Type {
		case constants.INPUT:
			n := p.createInputNode(stream)
			p.ActiveNodes = Append(p.ActiveNodes, n)
		case constants.OUTPUT:
			n := p.createOutputNode(stream)
			p.ActiveNodes = Append(p.ActiveNodes, n)
		default:
			return errors.New("unknown stream type")
		}
	}
	return nil
}

func (p *Program) LoadCode(code types.ProgramCode) error {
	if len(code.NodesCode) != constants.NodesNumber {
		return errors.New("wrong nodes number")
	}

	allInput := make([]InputCode, 0)
	for range constants.NodesNumber {
		allInput = append(allInput, NewInputCode())
	}

	for i, nc := range code.NodesCode {
		for _, line := range nc {
			formatted := strings.ToUpper(strings.TrimSpace(line))
			allInput[i].AddLine(formatted)
		}
	}

	for _, n := range p.Nodes {
		if err := n.ParseCode(&allInput[n.Index]); err != nil {
			return err
		}
		if len(n.Instructions) > 0 {
			p.ActiveNodes = Append(p.ActiveNodes, n)
		}
	}

	return nil
}

func (p *Program) createNode() *Node {
	n := NewNode()
	p.NodeList = Append(p.NodeList, n)
	return n
}

func (p *Program) createInputNode(stream types.Stream) *Node {
	inputNode := p.createNode()
	inputNode.Index = stream.Position
	belowNode := p.Nodes[stream.Position]

	inputNode.Ports[DOWN] = belowNode
	belowNode.Ports[UP] = inputNode

	for _, value := range stream.Values {
		ins := inputNode.CreateInstruction(MOV)
		ins.SrcType = NUMBER
		ins.Src.Number = value
		ins.DestType = ADDRESS
		ins.Dest.Direction = DOWN
	}

	ins := inputNode.CreateInstruction(JRO)
	ins.SrcType = NUMBER
	ins.Src.Number = 0

	return inputNode
}

func (p *Program) createOutputNode(stream types.Stream) *Node {
	outputNode := p.createNode()
	outputNode.Index = stream.Position + 8
	aboveNode := p.Nodes[stream.Position+8]

	outputNode.Ports[UP] = aboveNode
	aboveNode.Ports[DOWN] = outputNode

	ins := outputNode.CreateInstruction(MOV)
	ins.SrcType = ADDRESS
	ins.Src.Direction = UP
	ins.DestType = ADDRESS
	ins.Dest.Direction = ACC
	outputNode.CreateInstruction(OUT)

	p.Outputs = append(p.Outputs, NewOutput(stream.Position))
	outputNode.Output = p.Outputs[len(p.Outputs)-1]

	return outputNode
}
