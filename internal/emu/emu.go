package emu

type (
	Operation         uint8
	LocationType      uint8
	LocationDirection uint8
)

type Location struct {
	Number    int16
	Direction LocationDirection
}

type Instruction struct {
	Operation Operation
	SrcType   LocationType
	Src       Location
	DestType  LocationType
	Dest      Location
}

const (
	MOV Operation = iota
	SAV
	SWP
	SUB
	ADD
	NOP
	NEG
	JEZ
	JMP
	JNZ
	JGZ
	JLZ
	JRO
	ATA
)

const (
	NUMBER LocationType = iota
	ADDRESS
)

const (
	UP LocationDirection = iota
	RIGHT
	DOWN
	LEFT
	NIL
	ACC
	ANY
	LAST
)
