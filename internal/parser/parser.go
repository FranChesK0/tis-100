package parser

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

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

// TODO: use goroutines to call all fetch functions
// TODO: refactor FetchPuzzle function
func FetchPuzzle(fileName string) (*Puzzle, error) {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(fileName); err != nil {
		return &Puzzle{}, fmt.Errorf("unable to load lua script %s: %w", fileName, err)
	}

	title, err := fetchTitle(L)
	if err != nil {
		return &Puzzle{}, err
	}
	description, err := fetchDescription(L)
	if err != nil {
		return &Puzzle{}, err
	}
	streams, err := fetchStreams(L)
	if err != nil {
		return &Puzzle{}, err
	}
	layout, err := fetchLayout(L)
	if err != nil {
		return &Puzzle{}, err
	}

	return &Puzzle{
		Title:       title,
		Description: description,
		Streams:     streams,
		Layout:      layout,
	}, nil
}

func runLuaFunction(L *lua.LState, functionName string) (lua.LValue, error)

func fetchTitle(L *lua.LState) (string, error)

func fetchDescription(L *lua.LState) ([]string, error)

func fetchStreams(L *lua.LState) ([]Stream, error)

func fetchLayout(L *lua.LState) ([]TileType, error)
