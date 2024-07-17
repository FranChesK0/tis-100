package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FranChesK0/tis-100/internal/constants"
	"github.com/FranChesK0/tis-100/internal/types"
)

func SaveCode(dirPath string, code *types.ProgramCode) error {
	_, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("unable to read directory %s: %w", dirPath, err)
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.tis", strings.ToLower(code.Title)))

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file with name %s: %w", filePath, err)
	}
	defer file.Close()

	for i, node := range code.NodesCode {
		if _, err = file.WriteString(fmt.Sprintf("@%d\n", i+1)); err != nil {
			return fmt.Errorf("error while writing data to file %s: %w", filePath, err)
		}
		for _, str := range node {
			if _, err = file.WriteString(str + "\n"); err != nil {
				return fmt.Errorf("error while writing data to file %s: %w", filePath, err)
			}
		}
		if _, err = file.WriteString("\n"); err != nil {
			return fmt.Errorf("error while writing data to file %s: %w", filePath, err)
		}
	}

	return nil
}

func FetchCode(fileName string) (*types.ProgramCode, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %w", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nodesCode := make([][]string, constants.NodesNumber)

	var line string
	curNode := -1
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "@") {
			curNode++
			continue
		}

		nodesCode[curNode] = append(nodesCode[curNode], line)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading file %s: %w", fileName, err)
	}

	title := filepath.Base(file.Name())
	title = strings.ToUpper(strings.TrimSuffix(title, ".tis"))

	return &types.ProgramCode{
		Title:     title,
		NodesCode: nodesCode,
	}, nil
}
