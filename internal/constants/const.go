package constants

import "github.com/FranChesK0/tis-100/internal/types"

/* PROGRAM CONSTANTS */
const (
	MaxACC                = 999
	MinACC                = -999
	StreamTypesNumber     = 2
	IOPositionsNumber     = 4
	MaxStreamValuesLength = 30
	NodesNumber           = 12
	NodeTypesNumber       = 2
)

/* ENUMS */
const (
	INPUT types.StreamType = iota
	OUTPUT
)

const (
	COMPUTE types.NodeType = iota
	DAMAGED
)
