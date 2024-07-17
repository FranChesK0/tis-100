package emu

type Node struct {
	Index          uint8
	Blocked        bool
	CursorPosition uint8
	Instructions   []*Instruction
	ACC            int16
	BAK            int16
	OutputPort     *Node
	Last           *Node
	OutputValue    int16
	Ports          [4]*Node
}

type ReadResult struct {
	Blocked bool
	Value   int16
}

func NewNode() *Node

func (n *Node) CreateInstruction(op Operation) *Instruction

func (n *Node) ParseCode(ic *InputCode) error

func (n *Node) ParseLine(ic *InputCode, line string) error

func (n *Node) Read(locType LocationType, loc Location) (ReadResult, error)

func (n *Node) Write(dir LocationDirection, value int16) (bool, error)

func (n *Node) MoveCursor()

func (n *Node) Tick() error
