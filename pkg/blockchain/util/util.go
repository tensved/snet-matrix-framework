package util

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"math/big"
)

// EstimateGas returns a new bind.TransactOpts instance with zero gas limit
// based on the provided wallet's transaction options.
//
// Parameters:
//   - wallet: A pointer to bind.TransactOpts containing the wallet's transaction options.
//
// Returns:
//   - opts: A pointer to a new bind.TransactOpts with the same From and Signer values as the wallet,
//     but with a nil Value and a GasLimit of 0.
func EstimateGas(wallet *bind.TransactOpts) (opts *bind.TransactOpts) {
	return &bind.TransactOpts{
		From:     wallet.From,
		Signer:   wallet.Signer,
		Value:    nil,
		GasLimit: 0,
	}
}

// GetSignature generates a cryptographic signature for the provided message using the given ECDSA private key.
//
// Parameters:
//   - message: A byte slice containing the message to be signed.
//   - privateKeyECDSA: A pointer to an ecdsa.PrivateKey used to sign the message.
//
// Returns:
//   - signature: A byte slice containing the generated signature.
func GetSignature(message []byte, privateKeyECDSA *ecdsa.PrivateKey) (signature []byte) {
	hash := crypto.Keccak256(
		blockchain.HashPrefix32Bytes,
		crypto.Keccak256(message),
	)

	signature, err := crypto.Sign(hash, privateKeyECDSA)
	if err != nil {
		panic(fmt.Sprintf("Cannot sign message: %v", err))
	}

	return signature
}

// BigIntToBytes converts a big.Int value to a byte slice.
//
// Parameters:
//   - value: A pointer to a big.Int to be converted.
//
// Returns:
//   - []byte: A byte slice representing the big.Int value.
func BigIntToBytes(value *big.Int) []byte {
	return common.BigToHash(value).Bytes()
}

// AgixToCog converts an AGIX amount to its equivalent in COG (smallest unit).
//
// Parameters:
//   - iamount: The amount to be converted, which can be of type string, float64, int64, decimal.Decimal, or *decimal.Decimal.
//
// Returns:
//   - agix: A pointer to a big.Int representing the equivalent COG amount.
//   - err: An error if the conversion fails.
func AgixToCog(iamount any) (agix *big.Int, err error) {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, err = decimal.NewFromString(v)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert string to decimal")
			return nil, err
		}
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	default:
		log.Fatal().Msg("Unsupported type")
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(8))
	result := amount.Mul(mul)

	agix = new(big.Int)
	agix.SetString(result.String(), 10)

	return agix, err
}

// CogToAgix converts a COG amount to its equivalent in AGIX.
//
// Parameters:
//   - ivalue: The value to be converted, which can be of type string, *big.Int, or int.
//
// Returns:
//   - decimal.Decimal: The equivalent AGIX amount as a decimal.Decimal.
func CogToAgix(ivalue any) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	case int:
		value.SetInt64(int64(v))
	default:
		log.Error().Msg("Unsupported type")
		return decimal.Zero
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(8))
	num, err := decimal.NewFromString(value.String())
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert string to decimal")
	}
	result := num.DivRound(mul, 8)

	return result
}
