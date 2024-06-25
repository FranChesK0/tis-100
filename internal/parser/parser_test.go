package parser

import (
	"fmt"
	"os"
	"testing"
)

/* TESTS */

// FetchPuzzle

/* BENCHMARKS */

// FetchPuzzle

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
