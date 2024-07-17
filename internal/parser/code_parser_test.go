package parser_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/FranChesK0/tis-100/internal/constants"
	"github.com/FranChesK0/tis-100/internal/parser"
	"github.com/FranChesK0/tis-100/internal/types"
)

/* TESTS */

// SaveCode
func TestSaveCodeWithCorrectInput(t *testing.T) {
	dir, err := SetupDir(t, "test_parser")
	if err != nil {
		t.Fatal(err)
	}

	code := newProgramCode("TEST-SAVE-WITH-CORRECT-INPUT")
	filePath, err := parser.SaveCode(dir, code)
	if err != nil {
		t.Error("unexpected error")
	}

	if filePath != filepath.Join(dir, "test-save-with-correct-input.tis") {
		t.Error("unexpected path for file")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	strContent := string(content)

	if strContent != codeToString(code.NodesCode) {
		t.Error("data in file is not equal expected result")
	}
}

func TestSaveCodeWithNotExistingDirectory(t *testing.T) {
	code := newProgramCode("TEST-SAVE-WITH-NOT-EXISTING-DIRECTORY")
	_, err := parser.SaveCode("/notexistingdirectory", code)

	if err == nil {
		t.Error("expected to occure error")
	}
	expectedErr := "unable to read directory /notexistingdirectory"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("wrong error occurred. expected: %s, got: %s", expectedErr, err.Error())
	}
}

// FetchCode
func TestFetchCodeWithCorrectInput(t *testing.T) {
	expectedCode := newProgramCode("TEST-FETCH-CODE-WITH-CORRECT-INPUT")
	file, err := SetupCode(t, expectedCode, "test-fetch-code-with-correct-input")
	if err != nil {
		t.Fatal(err)
	}

	code, err := parser.FetchCode(file.Name())
	if err != nil {
		t.Error("unexpected error")
	}

	expectedCode.Title = strings.TrimSuffix(filepath.Base(file.Name()), ".tis")
	expectedCode.Title = strings.ToUpper(expectedCode.Title)
	if !reflect.DeepEqual(code, expectedCode) {
		t.Error("code is not equal expected result")
	}
}

func TestFetchCodeWithNotExistingFile(t *testing.T) {
	_, err := parser.FetchCode("/notexistingdirectory/notexistingfile.tis")
	if err == nil {
		t.Error("expected to occure error")
	}
	expectedErr := "unable to open file /notexistingdirectory/notexistingfile.tis"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("wrong error occurred. expected: %s, got: %s", expectedErr, err.Error())
	}
}

/* BENCHMARKS */

// SaveCode

// FetchCode

/* UTILS */
func SetupDir(t *testing.T, dirName string) (string, error) {
	dir, err := os.MkdirTemp("", dirName)
	if err != nil {
		return "", err
	}

	t.Cleanup(func() { os.RemoveAll(dir) })

	return dir, err
}

func SetupCode(t *testing.T, code *types.ProgramCode, fileName string) (*os.File, error) {
	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return file, err
	}

	strCode := codeToString(code.NodesCode)
	_, err = fmt.Fprint(file, strCode)
	if err != nil {
		return file, err
	}

	return file, nil
}

func newProgramCode(title string) *types.ProgramCode {
	nodesCode := make([][]string, 0)
	for range constants.NodesNumber {
		nodesCode = append(nodesCode, []string{"MOV UP DOWN", "MOV UP DOWN"})
	}
	return &types.ProgramCode{
		Title:     title,
		NodesCode: nodesCode,
	}
}

func codeToString(nodesCode [][]string) string {
	res := ""
	for i, node := range nodesCode {
		res += fmt.Sprintf("@%d\n", i+1)
		for _, line := range node {
			res += line + "\n"
		}
		res += "\n"
	}
	return res
}
