package apm

import (
	"apm/gadgets/memory"
	"github.com/argennon-project/csgo/transpiled/gnark/api"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"math/big"
	"testing"
)

type twoInstructionCircuit struct {
	Instructions                [2]frontend.Variable
	Keys, ValuesInit, ValuesOut [3]frontend.Variable
}

func (c *twoInstructionCircuit) Define(a frontend.API) error {
	api.Api = a
	s := NewState(c.Instructions[:], *memory.NewWritable(c.Keys[:], c.ValuesInit[:], len(c.ValuesInit)))
	s.Transit()
	s.Transit()
	s.AssertOutputIs(c.ValuesOut[:])
	return nil
}

type simpleProgramCircuit struct {
	Instructions                [3]frontend.Variable
	Keys, ValuesInit, ValuesOut [4]frontend.Variable
}

func (c *simpleProgramCircuit) Define(a frontend.API) error {
	api.Api = a
	s := NewState(c.Instructions[:], *memory.NewWritable(c.Keys[:], c.ValuesInit[:], len(c.ValuesInit)))
	for i := 0; i < 7; i++ {
		s.Transit()
	}
	s.AssertOutputIs(c.ValuesOut[:])
	return nil
}

func Test_Transition(t *testing.T) {
	assert := test.NewAssert(t)
	assert.ProverSucceeded(&twoInstructionCircuit{}, &twoInstructionCircuit{
		Instructions: [2]frontend.Variable{
			new(big.Int).SetBytes([]byte{0b011_00000, 0, 0, 0, 0, 0, 0, 0, 0}),
			new(big.Int).SetBytes([]byte{0b1, 0b111_01100}),
		},
		Keys:       [3]frontend.Variable{0, 15, 22},
		ValuesInit: [3]frontend.Variable{5, 4, 9},
		ValuesOut:  [3]frontend.Variable{5, 7, 9},
	})

	assert.ProverSucceeded(&twoInstructionCircuit{}, &twoInstructionCircuit{
		Instructions: [2]frontend.Variable{
			new(big.Int).SetBytes([]byte{0b1000111, 0b011_00000, 0, 0, 0, 0, 0, 0, 1, 0b011_01000}),
			new(big.Int).SetBytes([]byte{0}),
		},
		Keys:       [3]frontend.Variable{0, 9, 11},
		ValuesInit: [3]frontend.Variable{5, 4, 9},
		ValuesOut:  [3]frontend.Variable{5, 4, 571},
	})

	assert.ProverSucceeded(&twoInstructionCircuit{}, &twoInstructionCircuit{
		Instructions: [2]frontend.Variable{
			new(big.Int).SetBytes([]byte{1, 0b011_00100}),
			new(big.Int).SetBytes([]byte{1, 0b001_01111}),
		},
		Keys:       [3]frontend.Variable{0, 9, 11},
		ValuesInit: [3]frontend.Variable{5, 7, 3},
		ValuesOut:  [3]frontend.Variable{5, 4, 3},
	})

	assert.ProverSucceeded(&simpleProgramCircuit{}, &simpleProgramCircuit{
		Instructions: [3]frontend.Variable{
			new(big.Int).SetBytes([]byte{0b001_00000, 0, 0, 0, 0, 0, 0, 0, 0}),
			new(big.Int).SetBytes([]byte{0b001_10100}),
			new(big.Int).SetBytes([]byte{0b011_01100}),
		},
		Keys:       [4]frontend.Variable{0, 1, 3, 7},
		ValuesInit: [4]frontend.Variable{0, 3, 2, 10},
		ValuesOut:  [4]frontend.Variable{0, 3, 5, 10},
	})
}

type decodeCircuit struct {
	Instruction             frontend.Variable
	Op, MemRead, MemWr, Jmp frontend.Variable
	MemAddr, Operand        frontend.Variable
}

func (c *decodeCircuit) Define(a frontend.API) error {
	api.Api = a
	gotOp, gotMemRead, gotMemWr, gotJmp, gotSmall, gotBig := decodeInstruction(c.Instruction)
	a.AssertIsEqual(gotOp, c.Op)
	a.AssertIsEqual(gotMemRead, c.MemRead)
	a.AssertIsEqual(gotMemWr, c.MemWr)
	a.AssertIsEqual(gotJmp, c.Jmp)
	a.AssertIsEqual(gotSmall, c.MemAddr)
	a.AssertIsEqual(gotBig, c.Operand)
	return nil
}

func Test_decodeInstruction(t *testing.T) {
	assert := test.NewAssert(t)
	assert.ProverSucceeded(&decodeCircuit{}, &decodeCircuit{
		Instruction: 0b101_01,
		Op:          0b01,
		MemRead:     1,
		MemWr:       0,
		Jmp:         1,
		MemAddr:     0,
		Operand:     0,
	})

	assert.ProverSucceeded(&decodeCircuit{}, &decodeCircuit{
		Instruction: new(big.Int).SetBytes([]byte{0b101, 0b001_11001, 0, 0, 0, 0, 0, 0, 0b10010111, 0b110_010_10}),
		Op:          0b10,
		MemRead:     0,
		MemWr:       1,
		Jmp:         0,
		MemAddr:     new(big.Int).SetBytes([]byte{0b11001000, 0, 0, 0, 0, 0, 0b100, 0b10111110}),
		Operand:     0b101001,
	})
}
