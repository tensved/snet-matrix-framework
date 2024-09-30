package util_test

import (
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain/util"
	"math/big"
	"testing"
)

// TestAgixToCog tests the AgixToCog function for various inputs and expected outputs.
// It checks the conversion of AGIX to COG for different values.
//
// Parameters:
//   - t: The testing framework instance.
func TestAgixToCog(t *testing.T) {
	type testpair struct {
		value          string
		expectedResult *big.Int
	}

	rr, isOK := new(big.Int).SetString("1111100000000", 10)
	if !isOK {
		t.Error("failed to cast to big.Int")
		log.Error().Msg("failed to cast to big.Int")
		return
	}

	var tests = []testpair{
		{value: "0.00000001", expectedResult: new(big.Int).SetUint64(1)},
		{value: "0.0000001", expectedResult: new(big.Int).SetUint64(10)},
		{value: "0.000001", expectedResult: new(big.Int).SetUint64(100)},
		{value: "0.00001", expectedResult: new(big.Int).SetUint64(1000)},
		{value: "0.0001", expectedResult: new(big.Int).SetUint64(10000)},
		{value: "0.001", expectedResult: new(big.Int).SetUint64(100000)},
		{value: "0.01", expectedResult: new(big.Int).SetUint64(1000000)},
		{value: "0.1", expectedResult: new(big.Int).SetUint64(10000000)},
		{value: "1", expectedResult: new(big.Int).SetUint64(100000000)},
		{value: "10", expectedResult: new(big.Int).SetUint64(1000000000)},
		{value: "1.2345678", expectedResult: new(big.Int).SetUint64(123456780)},
		{value: "11111", expectedResult: rr},
	}

	for _, pair := range tests {
		resultAgix, err := util.AgixToCog(pair.value)
		if resultAgix == nil || err != nil {
			t.Error("Expected", pair.expectedResult, "got", resultAgix.String())
			return
		}
		if resultAgix.Cmp(pair.expectedResult) != 0 {
			t.Error("Expected", pair.expectedResult, "got", resultAgix.String())
		}
	}
}

// TestCogToAgix tests the CogToAgix function for various inputs and expected outputs.
// It checks the conversion of COG to AGIX for different values.
//
// Parameters:
//   - t: The testing framework instance.
func TestCogToAgix(t *testing.T) {
	type testpair struct {
		value          *big.Int
		expectedResult string
	}

	var tests = []testpair{
		{value: new(big.Int).SetUint64(100000000), expectedResult: "1"},
		{value: new(big.Int).SetUint64(10000000), expectedResult: "0.1"},
		{value: new(big.Int).SetUint64(1000000), expectedResult: "0.01"},
		{value: new(big.Int).SetUint64(100000), expectedResult: "0.001"},
		{value: new(big.Int).SetUint64(10000), expectedResult: "0.0001"},
		{value: new(big.Int).SetUint64(1000), expectedResult: "0.00001"},
		{value: new(big.Int).SetUint64(100), expectedResult: "0.000001"},
		{value: new(big.Int).SetUint64(10), expectedResult: "0.0000001"},
		{value: new(big.Int).SetUint64(1), expectedResult: "0.00000001"},
	}

	for _, pair := range tests {
		result := util.CogToAgix(pair.value)
		if result.String() != pair.expectedResult {
			t.Error("Expected", pair.expectedResult, "got", result.String())
			return
		}
	}
}
