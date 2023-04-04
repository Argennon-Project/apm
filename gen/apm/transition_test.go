package apm

import (
	"apm/gadgets/memory"
	"github.com/argennon-project/csgo/transpiled/gnark/api"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"testing"
)

type simpleProgramCircuit struct {
	Instructions                [2]frontend.Variable
	Keys, ValuesInit, ValuesOut [3]frontend.Variable
}

func (c *simpleProgramCircuit) Define(a frontend.API) error {
	api.Api = a
	s := NewState(c.Instructions[:], *memory.NewWritable(c.Keys[:], c.ValuesInit[:], len(c.ValuesInit)))
	s.Transit()
	s.Transit()
	s.AssertOutputIs(c.ValuesOut[:])
	return nil
}

func Test_Transition(t *testing.T) {
	assert := test.NewAssert(t)
	assert.ProverSucceeded(&simpleProgramCircuit{}, &simpleProgramCircuit{
		Instructions: [2]frontend.Variable{
			0b11_0000000000000000_000_00,
			0b1111_0000000000000000_011_00,
		},
		Keys:       [3]frontend.Variable{0, 15, 22},
		ValuesInit: [3]frontend.Variable{5, 4, 9},
		ValuesOut:  [3]frontend.Variable{5, 7, 9},
	})
}

type decodeCircuit struct {
	Instruction             frontend.Variable
	Op, MemRead, MemWr, Jmp frontend.Variable
	SmallOp, BigOp          frontend.Variable
}

func (c *decodeCircuit) Define(a frontend.API) error {
	api.Api = a
	gotOp, gotMemRead, gotMemWr, gotJmp, gotSmall, gotBig := decodeInstruction(c.Instruction)
	a.AssertIsEqual(gotOp, c.Op)
	a.AssertIsEqual(gotMemRead, c.MemRead)
	a.AssertIsEqual(gotMemWr, c.MemWr)
	a.AssertIsEqual(gotJmp, c.Jmp)
	a.AssertIsEqual(gotSmall, c.SmallOp)
	a.AssertIsEqual(gotBig, c.BigOp)
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
		SmallOp:     0,
		BigOp:       0,
	})

	assert.ProverSucceeded(&decodeCircuit{}, &decodeCircuit{
		Instruction: 0b1_0000000000000000_000_00,
		Op:          0b00,
		MemRead:     0,
		MemWr:       0,
		Jmp:         0,
		SmallOp:     0,
		BigOp:       1,
	})

	assert.ProverSucceeded(&decodeCircuit{}, &decodeCircuit{
		Instruction: 0b101_1100101110001110_010_10,
		Op:          0b10,
		MemRead:     0,
		MemWr:       1,
		Jmp:         0,
		SmallOp:     0b1100_1011_1000_1110,
		BigOp:       0b101,
	})
}
