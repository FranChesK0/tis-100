package parser

import (
	"errors"
	"fmt"

	"github.com/yuin/gopher-lua"

	"github.com/FranChesK0/tis-100/internal/constants"
	"github.com/FranChesK0/tis-100/internal/types"
)

type Stream struct {
	Type     types.StreamType
	Name     string
	Position uint8
	Values   []int16
}

type Puzzle struct {
	Title       string
	Description []string
	Streams     []Stream
	Layout      []types.NodeType
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

func runLuaFunction(L *lua.LState, functionName string) (lua.LValue, error) {
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal(functionName),
		NRet:    1,
		Protect: true,
	}); err != nil {
		return nil, fmt.Errorf("error while calling %s function: %w", functionName, err)
	}
	value := L.Get(-1)
	L.Pop(1)
	return value, nil
}

func fetchTitle(L *lua.LState) (string, error) {
	runResult, err := runLuaFunction(L, "GetTitle")
	if err != nil {
		return "", err
	}
	if title, ok := runResult.(lua.LString); ok { // check whether the value is a string
		return title.String(), nil
	}
	return "", errors.New("cannot process the result of the GetTitle function")
}

func fetchDescription(L *lua.LState) ([]string, error) {
	runResult, err := runLuaFunction(L, "GetDescription")
	if err != nil {
		return nil, err
	}
	if descTable, ok := runResult.(*lua.LTable); ok { // check whether the value is a lua table
		desc := make([]string, 0, descTable.Len())
		descTable.ForEach(func(_, value lua.LValue) {
			if descLine, ok := value.(lua.LString); ok { // check whether each value is a string
				desc = append(desc, descLine.String())
			}
		})
		if len(desc) == descTable.Len() { // return array if all values were fetched
			return desc, nil
		}
	}
	return nil, errors.New("cannot process the result of the GetDescription function")
}

func fetchStreams(L *lua.LState) ([]Stream, error) {
	runResult, err := runLuaFunction(L, "GetStreams")
	if err != nil {
		return nil, err
	}
	if streamsTable, ok := runResult.(*lua.LTable); ok { // check whether the value is a lua table
		streams := make([]Stream, 0, streamsTable.Len())
		streamsTable.ForEach(func(_, value lua.LValue) {
			if streamTable, ok := value.(*lua.LTable); ok && streamTable.Len() == 4 {
				typeValue, typeOk := streamTable.RawGetInt(1).(lua.LNumber)          // check whether stream type is a number
				nameValue, nameOk := streamTable.RawGetInt(2).(lua.LString)          // check whether stream name is a string
				posValue, posOk := streamTable.RawGetInt(3).(lua.LNumber)            // check whether stream position is a number
				valuesTable, valuesTableOk := streamTable.RawGetInt(4).(*lua.LTable) // check whether stream values is a lua table

				streamValues := make([]int16, 0, valuesTable.Len())
				if valuesTable.Len() <= constants.MaxStreamValuesLength { // check stream values length
					valuesTable.ForEach(func(_, value lua.LValue) {
						val, ok := value.(lua.LNumber) // check wether each stream value is a number
						valuesTableOk = valuesTableOk && ok
						if ok && constants.MinACC <= val &&
							val <= constants.MaxACC {
							streamValues = append(streamValues, int16(val))
						}
					})
				}

				// if everything ok with streams table format return it
				if typeOk && nameOk && posOk && valuesTableOk &&
					len(streamValues) == valuesTable.Len() &&
					0 <= typeValue && typeValue < constants.StreamTypesNumber &&
					0 <= posValue && posValue < constants.IOPositionsNumber {
					streams = append(streams, Stream{
						Type:     types.StreamType(typeValue),
						Name:     nameValue.String(),
						Position: uint8(posValue),
						Values:   streamValues,
					})
				}
			}
		})
		if len(streams) == streamsTable.Len() {
			return streams, nil
		}
	}
	return nil, errors.New("cannot process the result of the GetStreams function")
}

func fetchLayout(L *lua.LState) ([]types.NodeType, error) {
	runResult, err := runLuaFunction(L, "GetLayout")
	if err != nil {
		return nil, err
	}
	if layoutTable, ok := runResult.(*lua.LTable); ok &&
		layoutTable.Len() == constants.NodesNumber {
		layout := make([]types.NodeType, 0, layoutTable.Len())
		layoutTable.ForEach(func(_, value lua.LValue) {
			if tileType, ok := value.(lua.LNumber); ok {
				if 0 <= tileType && tileType < constants.NodeTypesNumber {
					layout = append(layout, types.NodeType(tileType))
				}
			}
		})
		if len(layout) == layoutTable.Len() { // return array if all values were fetched
			return layout, nil
		}
	}
	return nil, errors.New("cannot process the result of the GetLayout function")
}
