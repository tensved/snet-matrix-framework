package util

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"math/big"
	"matrix-ai-framework/pkg/blockchain"
	"regexp"
)

func CreatePrivateKey() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Error().Err(err).Msg("!!! can't create private key")
	}
	return hexutil.Encode(crypto.FromECDSA(privateKey))[2:]
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

func GetAddressFromPrivateKey(privateKey string) (address common.Address) {
	cryptoPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Error().Err(err)
	}

	publicKey := cryptoPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error().Err(err).Msg("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// WeiToEther wei to decimals
func WeiToEther(ivalue interface{}) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(18))
	num, err := decimal.NewFromString(value.String())
	if err != nil {
		log.Error().Err(err)
	}
	result := num.DivRound(mul, 18)

	return result
}

// EtherToWei decimals to wei
func EtherToWei(iamount interface{}) (wei *big.Int, err error) {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, err = decimal.NewFromString(v)
		if err != nil {
			log.Error().Err(err).Msg("can't convert string to decimal")
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
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(18))
	result := amount.Mul(mul)

	wei = new(big.Int)
	wei.SetString(result.String(), 10)

	return wei, err
}

// Deprecated
func weiToString(wei uint64) string {
	return fmt.Sprintf("%f", float64(wei)/1000000000000000000) // human readable без экспоненты
}

// Deprecated
func weiToDecimal(weis string) decimal.Decimal {
	wei := new(big.Int)
	wei.SetString(weis, 10)
	toDecimal := WeiToEther(wei)
	return toDecimal
}

//
//func ConvertStrAddressToCommonAddress(strAddress string) (commonAddress common.Address) {
//	// Remove the "0x" prefix if present
//	if len(strAddress) > 2 && strAddress[:2] == "0x" {
//		strAddress = strAddress[2:]
//	}
//
//	// Decode the address string from hex
//	addressBytes, err := hex.DecodeString(strAddress)
//	if err != nil {
//		panic(err)
//	}
//
//	// Hash the address bytes to ensure correct length
//	hash := sha256.Sum256(addressBytes)
//
//	// Convert the hash bytes to big.Int
//	var hashInt big.Int
//	hashInt.SetBytes(hash[:])
//
//	// Create a common.Address from the hashInt
//	commonAddress = common.BigToAddress(&hashInt)
//	return
//}

func ConvertStringToByte32(str string) [32]byte {
	log.Debug().Msgf("StringToByte32: %s", str)
	//// Ensure the input string is exactly 32 bytes long
	//if len(str) != 32 {
	//	fmt.Println("Input string must be exactly 32 bytes long")
	//	return [32]byte{}
	//}
	//
	//// Convert the string to bytes
	//stringBytes := []byte(str)
	//
	//// Create a fixed-size byte array with length 32
	//var byteArray [32]byte
	//
	//// Copy the bytes from the string into the byte array
	//copy(byteArray[:], stringBytes)
	//
	//log.Debug().Msgf("String converted to [32]byte: %x\n", byteArray)
	//return byteArray

	// The base64 encoded string
	//base64Str := "ZQyLCVXWRNbOkZ6y2KnxZ+QgPUWMKW6LVJObVnFq728="

	// Decode the base64 string to bytes
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return [32]byte{}
	}

	// Check if the length of the data is exactly 32 bytes
	if len(data) != 32 {
		fmt.Println("Decoded data is not 32 bytes long")
		return [32]byte{}
	}

	// Convert the slice to an array
	var array [32]byte
	copy(array[:], data[:32]) // Copy the data into the array

	// Print the array
	fmt.Println("Array:", array)
	return array
}

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
	//msgHash := crypto.Keccak256(message)
	//hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msgHash), msgHash)))

	signature, err := crypto.Sign(hash, privateKeyECDSA)
	if err != nil {
		panic(fmt.Sprintf("Cannot sign message: %v", err))
	}

	return signature
}

func BigIntToBytes(value *big.Int) []byte {
	return common.BigToHash(value).Bytes()
}

// AddressToHex converts Ethereum address to hex string representation.
func AddressToHex(address *common.Address) string {
	return address.Hex()
}

// BytesToBase64 converts array of bytes to base64 string.
func BytesToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// HexToBytes converts hex string to bytes array.
func HexToBytes(str string) []byte {
	return common.FromHex(str)
}

// HexToAddress converts hex string to Ethreum address.
func HexToAddress(str string) common.Address {
	return common.Address(common.BytesToAddress(HexToBytes(str)))
}

func StringToBytes32(str string) [32]byte {
	var byte32 [32]byte
	copy(byte32[:], str)
	return byte32
}

func RemoveSpecialCharactersfromHash(pString string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9=]")
	if err != nil {
		log.Panic().Err(err)
	}
	return reg.ReplaceAllString(pString, "")
}

func ConvertBase64Encoding(str string) ([32]byte, error) {
	var byte32 [32]byte
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Error().Err(err).Msgf("String Passed: %v", str)
		return byte32, err
	}
	copy(byte32[:], data[:])
	return byte32, nil
}

func GetSignerAddressFromMessage(message, signature []byte) (signer *common.Address, err error) {

	messageHash := crypto.Keccak256(
		blockchain.HashPrefix32Bytes,
		crypto.Keccak256(message),
	)
	log.Debug().Msgf("MessageHash: %s", hex.EncodeToString(messageHash))

	v, _, _, e := ParseSignature(signature)
	if e != nil {
		return nil, errors.New("incorrect signature length")
	}

	modifiedSignature := bytes.Join([][]byte{signature[0:64], {v % 27}}, nil)
	publicKey, e := crypto.SigToPub(messageHash, modifiedSignature)
	if e != nil {
		log.Warn().Msgf("Incorrect signature. modifiedSignature: %s", modifiedSignature)
		return nil, errors.New("incorrect signature data")
	}
	log.Debug().Msgf("publicKey: %s", publicKey)

	keyOwnerAddress := crypto.PubkeyToAddress(*publicKey)
	log.Debug().Msgf("Message signature parsed. keyOwnerAddress: %s", keyOwnerAddress)

	return &keyOwnerAddress, nil
}

// ParseSignature parses Ethereum signature.
func ParseSignature(jobSignatureBytes []byte) (uint8, [32]byte, [32]byte, error) {
	r := [32]byte{}
	s := [32]byte{}

	if len(jobSignatureBytes) != 65 {
		return 0, r, s, fmt.Errorf("job signature incorrect length")
	}

	v := uint8(jobSignatureBytes[64])%27 + 27
	copy(r[:], jobSignatureBytes[0:32])
	copy(s[:], jobSignatureBytes[32:64])

	return v, r, s, nil
}
