package interpreter

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// Payload include some filed about RAM status on specific index
type Payload struct {
	Ptr         *uint16
	Ram         []int16
	Cursor      *int
	Instruction instruction
	output      []byte
}

// instruction is a instruction node
type instruction struct {
	operator byte
	operand  uint16
}

type BFInterpreter struct {
	operations map[byte]func(*Payload)
	program    []instruction
	ram        []int16
}

// New returns an instance of Brainfuck interpreter
func New() *BFInterpreter {
	bf := BFInterpreter{
		operations: make(map[byte]func(*Payload)),
	}

	// initiating all default Brainfuck operations
	bf.initDefaultOpt()

	return &bf
}

// initDefaultOpt initiates all Brianfuck default operations
func (bf *BFInterpreter) initDefaultOpt() {
	bf.AddOpt('>', func(p *Payload) {
		*p.Ptr++
	})

	bf.AddOpt('<', func(p *Payload) {
		*p.Ptr--
	})

	bf.AddOpt('+', func(p *Payload) {
		p.Ram[*p.Ptr]++
	})

	bf.AddOpt('-', func(p *Payload) {
		p.Ram[*p.Ptr]--
	})

	// creating output
	bf.AddOpt('.', func(p *Payload) {
		p.output = append(p.output, byte(p.Ram[*p.Ptr]))
	})

	bf.AddOpt(',', func(p *Payload) {
		reader := bufio.NewReader(os.Stdin)
		read_val, _ := reader.ReadByte()
		p.Ram[*p.Ptr] = int16(read_val)
	})

	bf.AddOpt('[', func(p *Payload) {
		if p.Ram[*p.Ptr] == 0 {
			*p.Cursor = int(p.Instruction.operand)
		}
	})

	bf.AddOpt(']', func(p *Payload) {
		if p.Ram[*p.Ptr] > 0 {
			*p.Cursor = int(p.Instruction.operand)
		}
	})
}

// RemoveOpt is usefull for removing a operation in runtime
func (bf *BFInterpreter) RemoveOpt(tag byte) error {
	_, ok := bf.operations[tag]
	if !ok {
		return errors.New("invalid operation")
	}

	delete(bf.operations, tag)
	return nil
}

// AddOpt is usefull for adding a new operation in runtime
func (bf *BFInterpreter) AddOpt(tag byte, handler func(*Payload)) error {
	_, ok := bf.operations[tag]
	if ok {
		return errors.New("duplicate operation")
	}

	bf.operations[tag] = handler
	return nil
}

func (bf *BFInterpreter) Interpret(stream io.Reader) error {
	bf.program = make([]instruction, 0)

	var pc, jmp_pc uint16 = 0, 0
	jmp_stack := make([]uint16, 0)
	input, err := io.ReadAll(stream)
	if err != nil {
		return err
	}

	for _, c := range input {
		_, ok := bf.operations[c]
		if ok {
			// adding all input operations in program stream
			if c != '[' && c != ']' {
				bf.program = append(bf.program, instruction{c, 0})
			} else if c == '[' {
				bf.program = append(bf.program, instruction{c, 0})
				// adding to stack
				jmp_stack = append(jmp_stack, pc)
			} else {
				if len(jmp_stack) == 0 {
					return errors.New("compilation error")
				}

				// linking []
				jmp_pc = jmp_stack[len(jmp_stack)-1]
				jmp_stack = jmp_stack[:len(jmp_stack)-1]
				bf.program = append(bf.program, instruction{c, jmp_pc})
				bf.program[jmp_pc].operand = pc
			}
		} else {
			pc--
		}
		pc++
	}
	if len(jmp_stack) != 0 {
		return errors.New("compilation error")
	}
	return nil
}

func (bf *BFInterpreter) Execute() ([]byte, error) {
	bf.ram = make([]int16, 32000)
	var data_ptr uint16 = 0
	pl := Payload{
		Ptr:    &data_ptr,
		Ram:    bf.ram,
		output: make([]byte, 0),
	}

	for pc := 0; pc < len(bf.program); pc++ {
		handler, ok := bf.operations[bf.program[pc].operator]
		if !ok {
			return nil, errors.New("invalid operation")
		}

		// running instructions
		pl.Cursor = &pc
		pl.Instruction = bf.program[pc]
		handler(&pl)

	}

	return pl.output, nil
}
