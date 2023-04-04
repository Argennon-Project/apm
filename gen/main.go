package main

import (
	"github.com/argennon-project/csgo/transpiled/gnark/api"
	"github.com/argennon-project/csgo/transpiled/selector"
	"github.com/consensys/gnark/frontend"
)

//go:generate java -jar ../csgot.jar ../src ../gen

// MuxCircuit is a minimal circuit using a selector mux.
type MuxCircuit struct {
	Selector frontend.Variable    `gnark:",public"`
	In       [2]frontend.Variable `gnark:",public"`
	Expected frontend.Variable    `gnark:",public"`
}

// Define defines the arithmetic circuit.
func (c *MuxCircuit) Define(a frontend.API) error {
	api.Api = a
	result, _ := selector.Mux(c.Selector, c.In[:]...)
	api.AssertIsEqual(result, c.Expected)
	return nil
}

type TemporalCircuit struct {
	In    frontend.Variable
	Out   frontend.Variable
	steps int
}

/*
func (tc *TemporalCircuit) Define(a frontend.API) error {
	api.Api = a

	state := apm.NewState(tc.In)

	tc.steps = 10
	for i := 0; i < tc.steps; i++ {
		state.Transit()
	}

	state.AssertOutputIs(tc.Out)

	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit TemporalCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := TemporalCircuit{
		In:    2,
		Out:   1024,
		steps: 10,
	}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify
	// proof, _ := groth16.Prove(ccs, pk, witness, backend.WithHints(selector.GetHints()...))
	proof, _ := groth16.Prove(ccs, pk, witness)
	_ = groth16.Verify(proof, vk, publicWitness)
}

type InputState struct {
	Y frontend.Variable
}

type OutputState struct {
	X frontend.Variable
}

type State struct {
	x, y frontend.Variable
}

func (s *State) Transit() {
	s.x = api.Mul(s.x, s.y)
}

func (s *State) Init(input InputState) {
	s.y = input.Y
	s.x = 1
}

func (s *State) AssertOutputIs(output OutputState) {
	api.AssertIsEqual(s.x, output.X)
}

type TemporalState[In, Out any] interface {
	Transit()
	Init(input In)
	AssertOutputIs(output Out)
}*/
