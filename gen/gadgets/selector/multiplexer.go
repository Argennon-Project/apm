// generated from multiplexer.csgo

package selector

import (
	"github.com/argennon-project/csgo/transpiled/gnark/api"
	"github.com/consensys/gnark/backend/hint"
	"github.com/consensys/gnark/frontend"
)

import "math/big"

// Map is a key value associative array: the output will be values[i] such that keys[i] == queryKey.
// If keys do not contain queryKey, no proofs can be generated. If keys has more than one key that
// equals to queryKey, the output is undefined, and the output could be a linear combination of
// all the corresponding values with that queryKey.
//
// In case, keys and values do not have the same length, this function will panic.
func Map(queryKey frontend.Variable, keys, values []frontend.Variable) frontend.Variable {
	// we don't need this check, but we added it to produce more informative errors and disallow
	// len(keys) < len(values) which is supported by generateSelector.
	if len(keys) != len(values) {
		panic("The number of keys and values must be equal")
	}
	return generateSelector(false, queryKey, keys, values)
}

// Mux is an n to 1 multiplexer: out = inputs[sel]. In other words, it selects exactly one of its
// inputs based on sel. The index of inputs starts from zero.
//
// sel needs to be between 0 and n - 1 (inclusive), where n is the number of inputs, otherwise the
// proof will fail.
func Mux(sel frontend.Variable, inputs ...frontend.Variable) frontend.Variable {
	return generateSelector(true, sel, nil, inputs)
}

// generateSelector generates a circuit for a multiplexer or an associative array (map). If wantMux
// is true, a multiplexer is generated and keys are ignored. If wantMux is false, a map is
// generated, and we must have len(keys) <= len(values), or it panics.
func generateSelector(wantMux bool, sel frontend.Variable, keys, values []frontend.Variable) frontend.Variable {
	var indicators []frontend.Variable
	if wantMux {
		if len(values) == 2 {
			return api.Api.Select(sel, values[1], values[0])
		}
		indicators, _ = api.Compiler().NewHint(muxIndicators, len(values), sel)
	} else {
		indicators, _ = api.Compiler().NewHint(mapIndicators, len(keys), append(keys, sel)...)
	}

	var indicatorsSum, out frontend.Variable = 0, 0
	for i := 0; i < len(indicators); i++ {
		// Check that all indicators for inputs that are not selected, are zero.
		if wantMux {
			api.AssertIsEqual(api.Mul(indicators[i], api.Sub(sel, i)), 0)
		} else {
			api.AssertIsEqual(api.Mul(indicators[i], api.Sub(sel, keys[i])), 0)
		}

		indicatorsSum = api.Add(indicatorsSum, indicators[i])
		out = api.Add(out, api.Mul(indicators[i], values[i]))
	}
	// We need to check that the indicator of the selected input is exactly 1. We used a sum
	// constraint, because usually it is cheap.
	api.AssertIsEqual(indicatorsSum, 1)
	return out
}

// muxIndicators is a hint function used within [Mux] function. It must be
// provided to the prover when circuit uses it.
func muxIndicators(_ *big.Int, inputs []*big.Int, results []*big.Int) error {
	sel := inputs[0]
	for i := 0; i < len(results); i++ {
		// `i` is an int which can be int32 or int64. We convert `i` to int64 then to bigInt,
		// which is safe. We should not convert `sel` to int64.
		if sel.Cmp(big.NewInt(int64(i))) == 0 {
			results[i].SetUint64(1)
		} else {
			results[i].SetUint64(0)
		}
	}
	return nil
}

// mapIndicators is a hint function used within [Map] function. It must be
// provided to the prover when circuit uses it.
func mapIndicators(_ *big.Int, inputs []*big.Int, results []*big.Int) error {
	key := inputs[len(inputs)-1]
	// We must make sure that we are initializing all elements of results
	for i := 0; i < len(results); i++ {
		if key.Cmp(inputs[i]) == 0 {
			results[i].SetUint64(1)
		} else {
			results[i].SetUint64(0)
		}
	}
	return nil
}

func init() {
	hint.Register(muxIndicators)
	hint.Register(mapIndicators)
}
