package parser

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/FranChesK0/tis-100/internal/constants"
	"github.com/FranChesK0/tis-100/internal/types"
)

/* TESTS */

// FetchPuzzle
func TestFetchPuzzleWithCorrectScript(t *testing.T) {
	file, err := Setup(t, *NewScript(), "test_fetch_puzzle_with_correct_script.lua")
	if err != nil {
		t.Fatal(err)
	}

	puzzle, err := FetchPuzzle(file.Name())
	if err != nil {
		t.Error("unexpected error")
	}

	expectedPuzzle := types.Puzzle{
		Title:       "TEST",
		Description: []string{"TEST LINE 1", "TEST LINE 2"},
		Streams: []types.Stream{
			{
				Type:     constants.INPUT,
				Name:     "IN.TEST",
				Position: 0,
				Values:   []int16{1, 2, 3},
			},
			{
				Type:     constants.OUTPUT,
				Name:     "OUT.TEST",
				Position: 0,
				Values:   []int16{1, 2, 3},
			},
		},
		Layout: []types.NodeType{
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
			constants.COMPUTE,
		},
	}

	if !reflect.DeepEqual(*puzzle, expectedPuzzle) {
		t.Error("puzzle is not equal expected result")
	}
}

func TestFetchPuzzleWithWrongScrip(t *testing.T) {
	script := NewScript()
	script.GetTitle = []string{"func GetTitle()", "return", ";"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_script.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedErr := "unable to load lua script"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedErr, err.Error())
	}
}

func TestFetchPuzzleWithoutFunction(t *testing.T) {
	script := NewScript()
	script.GetTitle = []string{""}
	file, err := Setup(t, *script, "test_fetch_puzzle_without_function.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedErr := "error while calling GetTitle function:"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedErr, err.Error())
	}
}

func TestFetchPuzzleWithWrongTitle(t *testing.T) {
	script := NewScript()
	script.GetTitle = []string{"function GetTitle()", "return {}", "end"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_title.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedErr := "title is not a string"
	if err.Error() != expectedErr {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedErr, err.Error())
	}
}

func TestFetchPuzzleWithWrongDescriptionType(t *testing.T) {
	script := NewScript()
	script.GetDescription = []string{"function GetDescription()", "return 1", "end"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_description_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedErr := "description is not an array"
	if err.Error() != expectedErr {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedErr, err.Error())
	}
}

func TestFetchPuzzleWithWrongDescriptionLineType(t *testing.T) {
	script := NewScript()
	script.GetDescription = []string{"function GetDescription()", "return { 1, 2 }", "end"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_description_line_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedErr := "description line is not a string"
	if err.Error() != expectedErr {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedErr, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamsType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{"function GetStreams()", "return 1", "end"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_streams_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "streams is not an array"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{"function GetStreams()", "return { 1, 2 }", "end"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "stream is not an array"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamArugmentsNumber(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0 } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_arguments_number.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "wrong stream arguments number: expected 4, got 3"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamTypeValueType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { {}, \"IN.TEST\", 0, { 1, 2 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_type_value_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "first value of stream is not a StreamType value"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamTypeValue(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 5, \"IN.TEST\", 0, { 1, 2 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_type_value.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "first value of stream is not a StreamType value"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongNameType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 0, {}, 0, { 1, 2 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_name_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "second value of stream is not a string"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamPositionType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", \"0\", { 1, 2 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_position_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "third value of stream is not a number"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamPositionValue(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 5, { 1, 2 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_position_value.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := fmt.Sprintf(
		"position is not in range from 0 to %d",
		constants.IOPositionsNumber-1,
	)
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamValuesType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { {0, \"IN.TEST\", 0, 0 } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_values_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "fourth value of stream is not an array"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestFetchPuzzleWithWrongStreamValuesLength(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0, { 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_values_length.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := fmt.Sprintf(
		"wrong stream values number: expected <=%d, got 31",
		constants.MaxStreamValuesLength,
	)
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamValuesValueType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0, { \"1\", \"2\" } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_values_value_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "stream value is not a number"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamValues(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0, { 500000, 2 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_values.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := fmt.Sprintf(
		"stream value is not in range from %d to %d",
		constants.MinACC,
		constants.MaxACC,
	)
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongLayoutType(t *testing.T) {
	script := NewScript()
	script.GetLayout = []string{
		"function GetLayout()",
		"return 1",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_layout_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "layout is not an array"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongLayoutLength(t *testing.T) {
	script := NewScript()
	script.GetLayout = []string{
		"function GetLayout()",
		"return { 0, 0, 0, 0, 0, 0, 0, 0, 0 }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_layout_lenght.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := fmt.Sprintf("wrong nodes number: expected %d, got 9", constants.NodesNumber)
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongLayoutValueType(t *testing.T) {
	script := NewScript()
	script.GetLayout = []string{
		"function GetLayout()",
		"return { \"0\", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_layout_value_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "layout value is not a NodeType value"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongLayoutValues(t *testing.T) {
	script := NewScript()
	script.GetLayout = []string{
		"function GetLayout()",
		"return { 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_layout_values.lua")
	if err != nil {
		t.Fatal(err)
	}

	_, err = FetchPuzzle(file.Name())
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedError := "layout value is not a NodeType value"
	if err.Error() != expectedError {
		t.Errorf("wrong error occured. expected: %s, got: %s", expectedError, err.Error())
	}
	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

// runLuaFunction -> covered in previous tests
// fetchTitle -> covered in previous tests
// fetchgetDescription -> covered in previous tests
// fetchgetStreams -> covered in previous tests
// fetchgetLayout -> covered in previous tests

/* BENCHMARKS */

// FetchPuzzle
func BenchmarkFetchPuzzle(b *testing.B) {
	file, err := os.CreateTemp("", "bench.lua")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	defer os.Remove(file.Name())
	if err := writeScriptToFile(file, *NewScript()); err != nil {
		b.Fatal(err)
	}
	for range b.N {
		FetchPuzzle(file.Name())
	}
}

/* UTILS */
type Script struct {
	Beginning      []string
	GetTitle       []string
	GetDescription []string
	GetStreams     []string
	GetLayout      []string
}

func NewScript() *Script {
	return &Script{
		Beginning: []string{
			"local STREAM_INPUT = 0",
			"local STREAM_OUTPUT = 1",
			"local TILE_COMPUTE = 0",
			"local TILE_DAMAGED = 1",
		},
		GetTitle: []string{"function GetTitle()", "return \"TEST\"", "end"},
		GetDescription: []string{
			"function GetDescription()",
			"return { \"TEST LINE 1\", \"TEST LINE 2\" }",
			"end",
		},
		GetStreams: []string{
			"function GetStreams()",
			"return {",
			"{ STREAM_INPUT, \"IN.TEST\", 0, { 1, 2, 3 } },",
			"{ STREAM_OUTPUT, \"OUT.TEST\", 0, { 1, 2, 3 } },",
			"}",
			"end",
		},
		GetLayout: []string{
			"function GetLayout()",
			"return {",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"}",
			"end",
		},
	}
}

func (s Script) Get() []string {
	script := append(s.Beginning, s.GetTitle...)
	script = append(script, s.GetDescription...)
	script = append(script, s.GetStreams...)
	return append(script, s.GetLayout...)
}

func Setup(t *testing.T, script Script, fileName string) (*os.File, error) {
	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return file, nil
	}
	if err := writeScriptToFile(file, script); err != nil {
		return file, err
	}

	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file, nil
}

func writeScriptToFile(file *os.File, script Script) error {
	for _, line := range script.Get() {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return err
		}
	}
	return nil
}
