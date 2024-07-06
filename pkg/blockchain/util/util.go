package util

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"math/big"
	"matrix-ai-framework/pkg/blockchain"
)

func EstimateGas(wallet *bind.TransactOpts) (opts *bind.TransactOpts) {
	return &bind.TransactOpts{
		From:     wallet.From,
		Signer:   wallet.Signer,
		Value:    nil,
		GasLimit: 0,
	}
}

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

func BigIntToBytes(value *big.Int) []byte {
	return common.BigToHash(value).Bytes()
}

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
