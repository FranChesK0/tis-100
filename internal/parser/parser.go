package parser

import (
	"errors"
	"fmt"

	"github.com/yuin/gopher-lua"

	"github.com/FranChesK0/tis-100/internal/constants"
	"github.com/FranChesK0/tis-100/internal/types"
)

// TODO: use goroutines to call all fetch functions
// TODO: refactor FetchPuzzle function
func FetchPuzzle(fileName string) (*types.Puzzle, error) {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(fileName); err != nil {
		return &types.Puzzle{}, fmt.Errorf("unable to load lua script %s: %w", fileName, err)
	}

	title, err := fetchTitle(L)
	if err != nil {
		return &types.Puzzle{}, err
	}
	description, err := fetchDescription(L)
	if err != nil {
		return &types.Puzzle{}, err
	}
	streams, err := fetchStreams(L)
	if err != nil {
		return &types.Puzzle{}, err
	}
	layout, err := fetchLayout(L)
	if err != nil {
		return &types.Puzzle{}, err
	}

	return &types.Puzzle{
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
	title, ok := runResult.(lua.LString)
	if !ok {
		return "", errors.New("title is not a string")
	}
	return title.String(), nil
}

func fetchDescription(L *lua.LState) ([]string, error) {
	runResult, err := runLuaFunction(L, "GetDescription")
	if err != nil {
		return nil, err
	}
	descTable, ok := runResult.(*lua.LTable)
	if !ok {
		return nil, errors.New("description is not an array")
	}

	desc := make([]string, 0, descTable.Len())
	descTable.ForEach(func(_, value lua.LValue) {
		if err != nil {
			return
		}
		descLine, ok := value.(lua.LString)
		if !ok {
			err = errors.New("description line is not a string")
			return
		}
		desc = append(desc, descLine.String())
	})

	if err != nil {
		return nil, err
	}
	return desc, nil
}

func fetchStreams(L *lua.LState) ([]types.Stream, error) {
	runResult, err := runLuaFunction(L, "GetStreams")
	if err != nil {
		return nil, err
	}
	streamsTable, ok := runResult.(*lua.LTable)
	if !ok {
		return nil, errors.New("streams is not an array")
	}

	streams := make([]types.Stream, 0, streamsTable.Len())
	streamsTable.ForEach(func(_, value lua.LValue) {
		if err != nil {
			return
		}
		streamTable, ok := value.(*lua.LTable)
		if !ok {
			err = errors.New("stream is not an array")
			return
		}
		if streamTable.Len() != 4 {
			err = fmt.Errorf("wrong stream arguments number: expected 4, got %d", streamTable.Len())
			return
		}

		typeValue, ok := streamTable.RawGetInt(1).(lua.LNumber)
		if !ok || typeValue < 0 || typeValue >= constants.StreamTypesNumber {
			err = errors.New("first value of stream is not a StreamType value")
			return
		}
		nameValue, ok := streamTable.RawGetInt(2).(lua.LString)
		if !ok {
			err = errors.New("second value of stream is not a string")
			return
		}
		posValue, ok := streamTable.RawGetInt(3).(lua.LNumber)
		if !ok {
			err = errors.New("third value of stream is not a number")
			return
		}
		if posValue < 0 || posValue >= constants.IOPositionsNumber {
			err = fmt.Errorf("position is not in range from 0 to %d", constants.IOPositionsNumber-1)
			return
		}
		valuesTable, ok := streamTable.RawGetInt(4).(*lua.LTable)
		if !ok {
			err = errors.New("fourth value of stream is not an array")
			return
		}
		if valuesTable.Len() > constants.MaxStreamValuesLength {
			err = fmt.Errorf(
				"wrong stream values number: expected <=%d, got %d",
				constants.MaxStreamValuesLength,
				valuesTable.Len(),
			)
			return
		}

		streamValues := make([]int16, 0, valuesTable.Len())
		valuesTable.ForEach(func(_, value lua.LValue) {
			if err != nil {
				return
			}
			val, ok := value.(lua.LNumber)
			if !ok {
				err = errors.New("stream value is not a number")
				return
			}
			if val < constants.MinACC || val > constants.MaxACC {
				err = fmt.Errorf(
					"stream value is not in range from %d to %d",
					constants.MinACC,
					constants.MaxACC,
				)
				return
			}
			streamValues = append(streamValues, int16(val))
		})

		streams = append(streams, types.Stream{
			Type:     types.StreamType(typeValue),
			Name:     nameValue.String(),
			Position: uint8(posValue),
			Values:   streamValues,
		})
	})

	if err != nil {
		return nil, err
	}
	return streams, nil
}

func fetchLayout(L *lua.LState) ([]types.NodeType, error) {
	runResult, err := runLuaFunction(L, "GetLayout")
	if err != nil {
		return nil, err
	}
	layoutTable, ok := runResult.(*lua.LTable)
	if !ok {
		return nil, errors.New("layout is not an array")
	}
	if layoutTable.Len() != constants.NodesNumber {
		return nil, fmt.Errorf(
			"wrong nodes number: expected %d, got %d",
			constants.NodesNumber,
			layoutTable.Len(),
		)
	}

	layout := make([]types.NodeType, 0, layoutTable.Len())
	layoutTable.ForEach(func(_, value lua.LValue) {
		if err != nil {
			return
		}
		nodeType, ok := value.(lua.LNumber)
		if !ok {
			err = errors.New("layout value is not a NodeType value")
			return
		}
		if nodeType < 0 || nodeType >= constants.NodeTypesNumber {
			err = errors.New("layout value is not a NodeType value")
			return
		}
		layout = append(layout, types.NodeType(nodeType))
	})

	if err != nil {
		return nil, err
	}
	return layout, nil
}
