package interpreter

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpret(t *testing.T) {
	testCases := []struct {
		Name                     string
		Input                    string
		Output                   string
		AddNewOperand            func(*BFInterpreter)
		InterpreterExpectedError error
		ExecuteExpectedError     error
	}{
		{
			Name:   "successfull",
			Input:  ">+++++++++[<++++++++>-]<.>+++++++[<++++>-]<+.+++++++..+++.[-]>++++++++[<++++>-]<.>+++++++++++[<++++++++>-]<-.--------.+++.------.--------.[-]>++++++++[<++++>-]<+.[-]++++++++++.",
			Output: "Hello world!",
			AddNewOperand: func(*BFInterpreter) {

			},
			InterpreterExpectedError: nil,
			ExecuteExpectedError:     nil,
		},
		{
			Name:  "compilation error []]",
			Input: "[]]",
			AddNewOperand: func(*BFInterpreter) {

			},
			InterpreterExpectedError: errors.New("compilation error"),
			ExecuteExpectedError:     nil,
		},
		{
			Name:  "compilation error [[]",
			Input: "[[]",
			AddNewOperand: func(*BFInterpreter) {

			},
			InterpreterExpectedError: errors.New("compilation error"),
			ExecuteExpectedError:     nil,
		},
		{
			Name:   "successfull with new operand",
			Input:  ">+++*[<++++++++>-]<.>+++++++[<++++>-]<+.+++++++..+++.[-]>++++++++[<++++>-]<.>+++++++++++[<++++++++>-]<-.--------.+++.------.--------.[-]>++++++++[<++++>-]<+.[-]++++++++++.",
			Output: "Hello world!",
			AddNewOperand: func(bf *BFInterpreter) {
				bf.AddOpt('*', func(p *Payload) {
					p.Ram[*p.Ptr] *= p.Ram[*p.Ptr]
				})
			},
			InterpreterExpectedError: nil,
			ExecuteExpectedError:     nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			inter := New()
			testCase.AddNewOperand(inter)
			err := inter.Interpret(strings.NewReader(testCase.Input))
			assert.Equal(tt, testCase.InterpreterExpectedError, err)

			result, err := inter.Execute()
			assert.Equal(tt, testCase.ExecuteExpectedError, err)
			assert.Equal(tt, testCase.Output, strings.Trim(string(result), "\n"))
		})
	}
}
