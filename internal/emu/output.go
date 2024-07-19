package emu

type Output struct {
	Index  uint8
	Values []int16
}

func NewOutput(index uint8) *Output {
	return &Output{
		Index:  index,
		Values: make([]int16, 0),
	}
}

func (o *Output) AddValue(value int16) {
	o.Values = append(o.Values, value)
}
