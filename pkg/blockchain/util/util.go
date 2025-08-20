package util

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
)

// Deprecated: EstimateGas returns a new bind.TransactOpts instance with zero gas limit
// based on the provided wallet's transaction options.
//
// Parameters:
//   - wallet: A pointer to bind.TransactOpts containing the wallet's transaction options.
//
// Returns:
//   - opts: A pointer to a new bind.TransactOpts with the same From and Signer values as the wallet,
//     but with a nil Value and a GasLimit of 0.
func EstimateGas(wallet *bind.TransactOpts) *bind.TransactOpts {
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
func GetSignature(message []byte, privateKeyECDSA *ecdsa.PrivateKey) []byte {
	hash := crypto.Keccak256(
		blockchain.HashPrefix32Bytes,
		crypto.Keccak256(message),
	)

	signature, err := crypto.Sign(hash, privateKeyECDSA)
	if err != nil {
		panic(fmt.Sprintf("cannot sign message: %v", err))
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

// Deprecated: AgixToCog converts an AGIX amount to its equivalent in COG (smallest unit).
//
// Parameters:
//   - iamount: The amount to be converted, which can be of type string, float64, int64, decimal.Decimal, or *decimal.Decimal.
//
// Returns:
//   - agix: A pointer to a big.Int representing the equivalent COG amount.
//   - err: An error if the conversion fails.
func AgixToCog(iamount any) (agix *big.Int, err error) {
	base := 10
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, err = decimal.NewFromString(v)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert string to decimal")
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
		log.Fatal().Msg("unsupported type")
	}
	dec, pow := float64(10), float64(8)
	mul := decimal.NewFromFloat(dec).Pow(decimal.NewFromFloat(pow))
	result := amount.Mul(mul)

	agix = new(big.Int)
	agix.SetString(result.String(), base)

	return
}

// Deprecated: CogToAgix converts a COG amount to its equivalent in AGIX.
//
// Parameters:
//   - ivalue: The value to be converted, which can be of type string, *big.Int, or int.
//
// Returns:
//   - decimal.Decimal: The equivalent AGIX amount as a decimal.Decimal.
func CogToAgix(ivalue any) decimal.Decimal {
	value := new(big.Int)
	base := 10
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, base)
	case *big.Int:
		value = v
	case int:
		value.SetInt64(int64(v))
	default:
		log.Error().Msg("unsupported type")
		return decimal.Zero
	}
	dec, pow := float64(10), float64(8)
	mul := decimal.NewFromFloat(dec).Pow(decimal.NewFromFloat(pow))
	num, err := decimal.NewFromString(value.String())
	if err != nil {
		log.Error().Err(err).Msg("failed to convert string to decimal")
	}
	precision := int32(8)
	result := num.DivRound(mul, precision)

	return result
}

// Deprecated: DecodePaymentGroupID decodes a base64-encoded payment group ID.
func DecodePaymentGroupID(encoded string) ([32]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return [32]byte{}, err
	}
	var groupID [32]byte
	copy(groupID[:], decoded)
	return groupID, nil
}

func GetCallOpts(fromAddress common.Address, lastBlockNumber *big.Int) *bind.CallOpts {
	return &bind.CallOpts{
		Pending:     false,
		From:        fromAddress,
		BlockNumber: lastBlockNumber,
		BlockHash:   common.Hash{},
		Context:     context.Background(),
	}
}

func GetWatchOpts(lastBlockNumber *big.Int) *bind.WatchOpts {
	startBlock := lastBlockNumber.Uint64()
	return &bind.WatchOpts{
		Start:   &startBlock,
		Context: context.Background(),
	}
}

func GetFilterOpts(lastBlockNumber *big.Int) *bind.FilterOpts {
	start := uint64(0)
	end := lastBlockNumber.Uint64()
	return &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}
}

func GetTransactOpts(privateKeyECDSA *ecdsa.PrivateKey) *bind.TransactOpts {
	chainID, _ := strconv.Atoi(config.Blockchain.ChainID)
	opts, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(int64(chainID)))
	if err != nil {
		log.Error().Err(err).Msg("failed to create transactor")
	}
	return opts
}

func GetNewExpiration(lastBlockNumber, paymentExpirationThreshold *big.Int) *big.Int {
	offset := int64(240)
	blockOffset := big.NewInt(offset)
	defaultExpiration := new(big.Int).Add(lastBlockNumber, paymentExpirationThreshold)
	return new(big.Int).Add(defaultExpiration, blockOffset)
}

// Deprecated: RemoveProtocol removes the protocol scheme from a URL.
func RemoveProtocol(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(rawURL, parsedURL.Scheme+"://"), nil
}

// Deprecated: WaitingForChannelToOpen waits for a channel to open within the specified timeout.
func WaitingForChannelToOpen(channelOpens <-chan *blockchain.MultiPartyEscrowChannelOpen, errChan <-chan error, timeout time.Duration) *big.Int {
	select {
	case openEvent := <-channelOpens:
		log.Debug().Msgf("channel opened: %+v", openEvent)
		return openEvent.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msgf("error watching for OpenChannel: %v", err)
	case <-time.After(timeout):
		log.Error().Msg("timed out waiting for OpenChannel to complete")
	}
	return nil
}

// Deprecated: WaitingForChannelToExtend waits for a channel to extend within the specified timeout.
func WaitingForChannelToExtend(channelExtends <-chan *blockchain.MultiPartyEscrowChannelExtend, errChan <-chan error, timeout time.Duration) *big.Int {
	select {
	case extendEvent := <-channelExtends:
		log.Debug().Msgf("channel extended: %+v", extendEvent)
		return extendEvent.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msgf("error watching for ChannelExtend: %v", err)
	case <-time.After(timeout):
		log.Error().Msg("timed out waiting for ChannelExtend to complete")
	}
	return nil
}

// Deprecated: WaitingForChannelFundsToBeAdded waits for funds to be added to a channel within the specified timeout.
func WaitingForChannelFundsToBeAdded(channelAddFunds <-chan *blockchain.MultiPartyEscrowChannelAddFunds, errChan <-chan error, timeout time.Duration) *big.Int {
	select {
	case addFundsEvent := <-channelAddFunds:
		log.Debug().Msgf("channel funds added: %+v", addFundsEvent)
		return addFundsEvent.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msgf("error watching for ChannelExtendAndAddFunds: %v", err)
	case <-time.After(timeout):
		log.Error().Msg("timed out waiting for ChannelExtendAndAddFunds to complete")
	}
	return nil
}

// Deprecated: WaitingToDepositFundsToMPE waits for funds to be deposited to MPE within the specified timeout.
func WaitingToDepositFundsToMPE(channelDepositFunds <-chan *blockchain.MultiPartyEscrowDepositFunds, errChan <-chan error, timeout time.Duration) bool {
	select {
	case depositFundsEvent := <-channelDepositFunds:
		log.Debug().Msgf("deposited to MPE: %+v", depositFundsEvent)
		return true
	case err := <-errChan:
		log.Error().Err(err).Msgf("error watching for deposit: %v", err)
	case <-time.After(timeout):
		log.Error().Msg("timed out waiting for deposit to complete")
	}
	return false
}

// Deprecated: IsChannelValid checks if a channel is valid based on funds and expiration.
func IsChannelValid(filteredEvent *blockchain.MultiPartyEscrowChannelOpen, bigIntPrice *big.Int, newExpiration *big.Int) (bool, bool) {
	hasSufficientFunds := filteredEvent.Amount.Cmp(bigIntPrice) >= 0
	isValidExpiration := filteredEvent.Expiration.Cmp(newExpiration) > 0
	return hasSufficientFunds, isValidExpiration
}
