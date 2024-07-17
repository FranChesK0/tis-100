package emu

import (
	"errors"
	"strconv"
	"strings"

	"github.com/FranChesK0/tis-100/internal/constants"
)

type Node struct {
	Index          uint8
	Blocked        bool
	CursorPosition uint8
	Instructions   []*Instruction
	ACC            int16
	BAK            int16
	OutputPort     *Node
	Last           *Node
	OutputValue    int16
	Ports          [4]*Node
}

type ReadResult struct {
	Blocked bool
	Value   int16
}

func NewNode() *Node {
	return &Node{
		Instructions: make([]*Instruction, 0),
		Ports:        [4]*Node{nil, nil, nil, nil},
	}
}

func (n *Node) CreateInstruction(op Operation) *Instruction {
	ins := &Instruction{
		Operation: op,
	}
	n.Instructions = append(n.Instructions, ins)
	return ins
}

func (n *Node) ParseCode(ic *InputCode) error {
	for i, line := range ic.Lines {
		if ind := strings.Index(line, ":"); ind != -1 {
			label := line[:ind]
			ic.Labels[label] = uint8(i)

			rem := strings.TrimSpace(line[ind+1:])
			if len(rem) == 0 {
				rem = "NOP"
			}
			ic.Lines[i] = rem
		}
	}

	var err error
	for _, line := range ic.Lines {
		if err = n.ParseLine(ic, line); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) ParseLine(ic *InputCode, line string) error {
	if len(line) <= 2 {
		return errors.New("invalid line length")
	}

	strIns := line[:3]
	insMap := map[string]Operation{
		"SUB": SUB,
		"ADD": ADD,
		"JEZ": JEZ,
		"JMP": JMP,
		"JNZ": JNZ,
		"JGZ": JGZ,
		"JLZ": JLZ,
		"JRO": JRO,
		"SAV": SAV,
		"SWP": SWP,
		"NOP": NOP,
		"NEG": NEG,
		"ATA": ATA,
	}
	var err error
	switch strIns {
	case "MOV":
		if err = n.parseMov(line); err != nil {
			return err
		}
	case "SUB", "ADD", "JEZ", "JMP", "JNZ", "JGZ", "JLZ", "JRO":
		if err = n.parseOneArg(ic, line, insMap[strIns]); err != nil {
			return err
		}
	case "SAV", "SWP", "NOP", "NEG", "ATA":
		n.CreateInstruction(insMap[strIns])
	default:
		return errors.New("invalid instruction")
	}

	return nil
}

func (n *Node) Read(locType LocationType, loc Location) (ReadResult, error) {
	res := ReadResult{}

	if n.OutputPort != nil {
		return res, nil
	}
	if locType == NUMBER {
		res.Value = loc.Number
		return res, nil
	}

	switch loc.Direction {
	case NIL:
		res.Value = 0
	case ACC:
		res.Value = n.ACC
	case UP, RIGHT, DOWN, LEFT, ANY, LAST:
		readFrom := n.getInputPort(loc.Direction)
		if readFrom != nil && readFrom.OutputPort == n {
			res.Value = readFrom.OutputValue
			res.Blocked = false

			readFrom.OutputValue = 0
			readFrom.OutputPort = nil
			readFrom.MoveCursor()

			if loc.Direction == ANY {
				n.Last = readFrom
			}
		} else if readFrom != nil && loc.Direction == LAST {
			res.Value = 0
		} else {
			res.Blocked = true
		}
	default:
		return ReadResult{}, errors.New("unknown direction")
	}

	return res, nil
}

func (n *Node) Write(dir LocationDirection, value int16) (bool, error) {
	switch dir {
	case ACC:
		n.ACC = value
	case UP, RIGHT, DOWN, LEFT, ANY, LAST:
		dest := n.getOutputPort(dir)
		if dest != nil && n.OutputPort == nil {
			n.OutputPort = dest
			n.OutputValue = value
			if dir == ANY {
				n.Last = dest
			}
		}
		return true, nil
	case NIL:
		return false, errors.New("unable to write")
	default:
		return false, errors.New("nowhere to write")
	}

	return false, nil
}

func (n *Node) MoveCursor() {
	n.CursorPosition++
}

func (n *Node) Tick() error {
	n.Blocked = true

	if n.CursorPosition >= uint8(len(n.Instructions)) {
		n.CursorPosition = 0
	}
	ins := n.Instructions[n.CursorPosition]

	switch ins.Operation {
	case MOV:
		read, err := n.Read(ins.SrcType, ins.Src)
		if err != nil {
			return err
		}
		if read.Blocked {
			return nil
		}

		blocked, err := n.Write(ins.Dest.Direction, read.Value)
		if err != nil {
			return err
		}
		if blocked {
			return nil
		}
	case ADD:
		read, err := n.Read(ins.SrcType, ins.Src)
		if err != nil {
			return err
		}
		if read.Blocked {
			return nil
		}

		n.ACC += read.Value
		n.normalizeACC()
	case SUB:
		read, err := n.Read(ins.SrcType, ins.Src)
		if err != nil {
			return err
		}
		if read.Blocked {
			return nil
		}

		n.ACC -= read.Value
		n.normalizeACC()
	case JMP:
		n.setCursorPosition(ins.Src.Number)
		return nil
	case JRO:
		n.setCursorPosition(int16(n.CursorPosition) + ins.Src.Number)
		return nil
	case JEZ:
		if n.ACC == 0 {
			n.setCursorPosition(ins.Src.Number)
			return nil
		}
	case JGZ:
		if n.ACC > 0 {
			n.setCursorPosition(ins.Src.Number)
			return nil
		}
	case JLZ:
		if n.ACC < 0 {
			n.setCursorPosition(ins.Src.Number)
			return nil
		}
	case JNZ:
		if n.ACC != 0 {
			n.setCursorPosition(ins.Src.Number)
			return nil
		}
	case SWP:
		tmp := n.BAK
		n.BAK = n.ACC
		n.ACC = tmp
	case SAV:
		n.BAK = n.ACC
	case NEG:
		n.ACC *= -1
	case NOP:
	case ATA:
		// TODO: process the output values
	default:
		return errors.New("unknown operation")
	}

	n.Blocked = false
	n.MoveCursor()
	return nil
}

func (n *Node) parseMov(line string) error {
	if len(line) <= 3 {
		return errors.New("wrong mov instruction format")
	}
	rem := line[4:]
	var tokens []string
	if strings.Contains(rem, ", ") {
		tokens = strings.Split(rem, ", ")
	} else if strings.Contains(rem, ",") {
		tokens = strings.Split(rem, ",")
	} else {
		tokens = strings.Split(rem, " ")
	}
	if len(tokens) != 2 {
		return errors.New("wrong mov instruction format")
	}

	ins := n.CreateInstruction(MOV)
	var err error
	if err = parseLocation(tokens[0], &ins.SrcType, &ins.Src); err != nil {
		return err
	}
	if err = parseLocation(tokens[1], &ins.DestType, &ins.Dest); err != nil {
		return err
	}

	return nil
}

func (n *Node) parseOneArg(ic *InputCode, line string, op Operation) error {
	if len(line) <= 3 {
		return errors.New("wrong one arg instruction format")
	}
	rem := line[4:]
	ins := n.CreateInstruction(op)

	switch op {
	case JEZ, JMP, JNZ, JGZ, JLZ:
		for label, pos := range ic.Labels {
			if rem == label {
				ins.SrcType = NUMBER
				ins.Src.Number = int16(pos)
			}
		}
	default:
		if err := parseLocation(rem, &ins.SrcType, &ins.Src); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) getInputPort(dir LocationDirection) *Node {
	switch dir {
	case ANY:
		dirs := []LocationDirection{LEFT, RIGHT, UP, DOWN}
		for _, d := range dirs {
			port := n.Ports[d]
			if port != nil && port.OutputPort == n {
				return port
			}
		}
	case LAST:
		return n.Last
	default:
		return n.Ports[dir]
	}
	return nil
}

func (n *Node) getOutputPort(dir LocationDirection) *Node {
	switch dir {
	case ANY:
		dirs := []LocationDirection{UP, LEFT, RIGHT, DOWN}
		for _, d := range dirs {
			port := n.Ports[d]
			if port != nil {
				ins := port.Instructions[port.CursorPosition]
				if ins.Operation == MOV && ins.SrcType == ADDRESS &&
					(ins.Src.Direction == ANY || port.Ports[ins.Src.Direction] == n) {
					return port
				}
			}
		}
	case LAST:
		return n.Last
	default:
		return n.Ports[dir]
	}
	return nil
}

func (n *Node) setCursorPosition(pos int16) {
	if pos >= int16(len(n.Instructions)) || pos < 0 {
		pos = 0
	}
	n.CursorPosition = uint8(pos)
}

func (n *Node) normalizeACC() {
	if n.ACC > constants.MaxACC {
		n.ACC = constants.MaxACC
	}
	if n.ACC < constants.MinACC {
		n.ACC = constants.MinACC
	}
}

func parseLocation(strLoc string, locType *LocationType, loc *Location) error {
	if strLoc == "" {
		return errors.New("no source was found")
	}

	*locType = ADDRESS
	switch strLoc {
	case "UP":
		loc.Direction = UP
	case "DOWN":
		loc.Direction = DOWN
	case "LEFT":
		loc.Direction = LEFT
	case "RIGHT":
		loc.Direction = RIGHT
	case "ACC":
		loc.Direction = ACC
	case "NIL":
		loc.Direction = NIL
	case "ANY":
		loc.Direction = ANY
	case "LAST":
		loc.Direction = LAST
	default:
		num, err := strconv.Atoi(strLoc)
		if err != nil {
			return err
		}

		*locType = NUMBER
		loc.Number = int16(num)
	}

	return nil
}
