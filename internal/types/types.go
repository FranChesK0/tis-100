package types

type (
	StreamType uint8
	NodeType   uint8
)

type Stream struct {
	Type     StreamType
	Name     string
	Position uint8
	Values   []int16
}

type Puzzle struct {
	Title       string
	Description []string
	Streams     []Stream
	Layout      []NodeType
}

type ProgramCode struct {
	Title     string
	NodesCode [][]string
}
