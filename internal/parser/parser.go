package parser

import "github.com/yuin/gopher-lua"

// TODO: move const and type to constants
type (
	StreamType uint8
	TileType   uint8
)

const (
	INPUT StreamType = iota
	OUTPUT
)

const (
	COMPUTE TileType = iota
	DAMAGED
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
	Layout      []TileType
}

func FetchPuzzle(fileName string) (*Puzzle, error)

func runLuaFunction(L *lua.LState, functionName string) (lua.LValue, error)

func fetchTitle(L *lua.LState) (string, error)

func fetchDescription(L *lua.LState) ([]string, error)

func fetchStreams(L *lua.LState) ([]Stream, error)

func fetchLayout(L *lua.LState) ([]TileType, error)
