package main

import (
	"apm/gen/gadgets/selector"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

//go:generate java -jar csgot.jar csgo gen

// MuxCircuit is a minimal circuit using a selector mux.
type MuxCircuit struct {
	Selector frontend.Variable     `gnark:",public"`
	In       [10]frontend.Variable `gnark:",public"`
	Expected frontend.Variable     `gnark:",public"`
}

// Define defines the arithmetic circuit.
func (c *MuxCircuit) Define(api frontend.API) error {
	selector.Configure(api)
	result := selector.Mux(c.Selector, c.In[:]...)
	api.AssertIsEqual(result, c.Expected)
	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit MuxCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := MuxCircuit{
		Selector: 2,
		In:       [10]frontend.Variable{0, 2, 4, 6, 8, 10, 12, 14, 16, 18},
		Expected: 4,
	}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify
	// proof, _ := groth16.Prove(ccs, pk, witness, backend.WithHints(selector.GetHints()...))
	proof, _ := groth16.Prove(ccs, pk, witness)
	groth16.Verify(proof, vk, publicWitness)
}