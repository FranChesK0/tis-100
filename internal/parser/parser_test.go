package parser

import (
	"fmt"
	"os"
	"testing"
)

/* TESTS */

// FetchPuzzle
// TODO: add comparison with Puzzle struct
func TestFetchPuzzleWithCorrectScript(t *testing.T) {
	file, err := Setup(t, *NewScript(), "test_fetch_puzzle_with_correct_script.lua")
	if err != nil {
		t.Fatal(err)
	}

	puzzle, err := FetchPuzzle(file.Name())
	if err != nil {
		t.Error("expected error")
	}
	_ = puzzle
}

func TestFetchPuzzleWithWrongTitle(t *testing.T) {
	script := NewScript()
	script.GetTitle = []string{"function GetTitle()", "return {}", "end"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_title.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongDescription(t *testing.T) {
	script := NewScript()
	script.GetDescription = []string{"function GetDescription()", "return { {} }", "end"}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_description.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamTableLength(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 1, \"1\", 1 } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_table_length.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { {}, \"1\", 1, {} } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamTypeValue(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 5, \"1\", 1, {} } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_type_value.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongNameType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 1, {}, 1, {} } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_name_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamPositionType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 1, \"1\", \"1\", {} } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_position_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamPositionValue(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 1, \"1\", 5, {} } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_position_value.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamValuesType(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 1, \"1\", 1, { \"1\", \"2\" } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_values_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongStreamValues(t *testing.T) {
	script := NewScript()
	script.GetStreams = []string{
		"function GetStreams()",
		"return { { 1, \"1\", 1, { 500000, 2 } } }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_stream_values.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

// TODO: func TestFetchPuzzleWithWrongStreamValuesLength(t *testing.T)

func TestFetchPuzzleWithWrongLayoutType(t *testing.T) {
	script := NewScript()
	script.GetLayout = []string{
		"function GetLayout()",
		"return { \"1\", 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1 }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_layout_type.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongLayoutLength(t *testing.T) {
	script := NewScript()
	script.GetLayout = []string{
		"function GetLayout()",
		"return { 1, 1, 1, 1, 1, 1, 1, 1, 1 }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_layout_lenght.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

func TestFetchPuzzleWithWrongLayoutValues(t *testing.T) {
	script := NewScript()
	script.GetLayout = []string{
		"function GetLayout()",
		"return { 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1 }",
		"end",
	}
	file, err := Setup(t, *script, "test_fetch_puzzle_with_wrong_layout_values.lua")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := FetchPuzzle(file.Name()); err == nil {
		t.Error("expected to occure error")
	}
}

// getTitle -> covered in previous tests
// getDescription -> covered in previous tests
// getStreams -> covered in previous tests
// getLayout -> covered in previous tests

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
			"{ STREAM_INPUT, \"IN.A\", 0, { 1, 2, 3 } },",
			"{ STREAM_OUTPUT, \"OUT.A\", 0, { 1, 2, 3 } },",
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
