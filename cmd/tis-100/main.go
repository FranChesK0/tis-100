package main

import (
	"fmt"
	"os"

	"github.com/FranChesK0/tis-100/internal/tui"
)

func main() {
	if err := tui.ProgramRun(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
