package emu

import "github.com/FranChesK0/tis-100/internal/types"

type Program struct {
	NodeList    *NodeList
	ActiveNodes *NodeList
}

func NewProgram() *Program

func (p *Program) Tick() (bool, error)

func (p *Program) LoadStreams(streams []types.Stream) error

func (p *Program) LoadCode(code types.ProgramCode) error
