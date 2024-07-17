package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FranChesK0/tis-100/internal/types"
)

func SaveCode(dirPath string, code *types.ProgramCode) error {
	_, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("unable to read directory %s: %w", dirPath, err)
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.tis", code.Title))

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
