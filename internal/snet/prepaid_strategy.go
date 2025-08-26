package snet

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain/util"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	"google.golang.org/grpc/metadata"
)

const (
	PaymentChannelIDHeader    = "snet-payment-channel-id"
	PaymentChannelNonceHeader = "snet-payment-channel-nonce"
	PrePaidAuthTokenHeader    = "snet-prepaid-auth-token-bin"
	PrefixInSignature         = "__MPE_claim_message"
)

// PaymentChannelHandler implements payment channel call strategy for blockchain services
type PaymentChannelHandler struct {
	ethClient       blockchain.Ethereum
	grpcManager     *grpcmanager.GRPCClientManager
	serviceMetadata *db.SnetService
	privateKeyECDSA *ecdsa.PrivateKey
	signerAddress   common.Address
	Token           string
	tokenClient     TokenServiceClient
	channelID       *big.Int
	nonce           *big.Int
	signedAmount    *big.Int
	mpeAddress      common.Address
}

// NewPaymentChannelHandler creates a new payment channel call strategy
func NewPaymentChannelHandler(evm blockchain.Ethereum, grpc *grpcmanager.GRPCClientManager, serviceMetadata *db.SnetService, privateKey *ecdsa.PrivateKey, callCount uint64, database db.Service) (Strategy, error) {
	logger := log.With().
		Str("service_id", serviceMetadata.SnetID).
		Str("service_url", serviceMetadata.URL).
		Uint64("call_count", callCount).
		Logger()

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	orgGroup, err := database.GetSnetOrgGroup(serviceMetadata.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get org group: %w", err)
	}

	recipient := common.HexToAddress(orgGroup.PaymentAddress)
	mpeAddress := common.HexToAddress(serviceMetadata.MPEAddress)

	logger.Debug().
		Str("from_address", fromAddress.Hex()).
		Str("recipient_address", recipient.Hex()).
		Str("mpe_address", mpeAddress.Hex()).
		Msg("addresses configured")

	groupIDBytes, err := base64.StdEncoding.DecodeString(serviceMetadata.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode group ID: %w", err)
	}
	var groupID [32]byte
	copy(groupID[:], groupIDBytes)

	currentBlockNumber, err := evm.Client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get current block: %w", err)
	}

	transactOpts := util.GetTransactOpts(privateKey)

	opts := &blockchain.BindOpts{
		Call:     util.GetCallOpts(fromAddress, big.NewInt(int64(currentBlockNumber))),
		Transact: transactOpts,
		Watch:    util.GetWatchOpts(big.NewInt(int64(currentBlockNumber))),
		Filter:   util.GetFilterOpts(big.NewInt(int64(currentBlockNumber))),
	}

	chans := &blockchain.ChansToWatch{
		ChannelOpens:    make(chan *blockchain.MultiPartyEscrowChannelOpen),
		ChannelExtends:  make(chan *blockchain.MultiPartyEscrowChannelExtend),
		ChannelAddFunds: make(chan *blockchain.MultiPartyEscrowChannelAddFunds),
		DepositFunds:    make(chan *blockchain.MultiPartyEscrowDepositFunds),
		Err:             make(chan error),
	}

	senders := []common.Address{fromAddress}
	recipients := []common.Address{recipient}
	groupIDs := [][32]byte{groupID}

	filteredChannel, err := evm.FilterChannels(senders, recipients, groupIDs, opts.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to filter channels: %w", err)
	}

	priceInCogs := big.NewInt(int64(serviceMetadata.Price))

	paymentExpirationThreshold := orgGroup.PaymentExpirationThreshold

	newExpiration := util.GetNewExpiration(big.NewInt(int64(currentBlockNumber)), paymentExpirationThreshold)

	if filteredChannel != nil {
		if filteredChannel.Recipient != recipient {
			logger.Info().
				Str("found_recipient", filteredChannel.Recipient.Hex()).
				Str("expected_recipient", recipient.Hex()).
				Msg("found channel has different recipient, creating new one")
			filteredChannel = nil
		} else {
			_, err := GetChannelInfoFromService(grpc, context.Background(), serviceMetadata.URL, mpeAddress, filteredChannel.ChannelId, currentBlockNumber, privateKey)
			if err != nil {
				logger.Info().
					Str("channel_id", filteredChannel.ChannelId.String()).
					Err(err).
					Msg("cannot get channel state from daemon, creating new one")
				filteredChannel = nil
			}
		}
	}

	if filteredChannel == nil {
		logger.Info().Msg("creating new payment channel")

		channelID, err := evm.OpenNewChannel(priceInCogs, newExpiration, opts, chans, senders, recipients, groupIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to open new channel: %w", err)
		}

		var channelState *ChannelStateReply
		for i := 0; i < 5; i++ {
			time.Sleep(2 * time.Second)
			channelState, err = GetChannelInfoFromService(grpc, context.Background(), serviceMetadata.URL, mpeAddress, channelID, currentBlockNumber, privateKey)
			if err == nil {
				break
			}
			logger.Debug().
				Int("attempt", i+1).
				Err(err).
				Msg("failed to get channel state, retrying")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get new channel state after retries: %w", err)
		}

		nonce := new(big.Int).SetBytes(channelState.GetCurrentNonce())
		if nonce == nil {
			return nil, errors.New("error while getting current nonce")
		}

		currentSignedAmount := new(big.Int).SetBytes(channelState.GetCurrentSignedAmount())
		if currentSignedAmount == nil {
			return nil, errors.New("error while getting signed amount")
		}

		increment := new(big.Int).Mul(priceInCogs, big.NewInt(int64(callCount)))
		signedAmount := new(big.Int).Add(currentSignedAmount, increment)

		grpcClient, err := grpc.GetClient(serviceMetadata.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to get gRPC client: %w", err)
		}

		logger.Info().
			Str("channel_id", channelID.String()).
			Str("nonce", nonce.String()).
			Str("signed_amount", signedAmount.String()).
			Msg("payment channel handler created with new channel")

		return &PaymentChannelHandler{
			ethClient:       evm,
			grpcManager:     grpc,
			serviceMetadata: serviceMetadata,
			privateKeyECDSA: privateKey,
			signerAddress:   fromAddress,
			tokenClient:     NewTokenServiceClient(grpcClient.Conn),
			channelID:       channelID,
			nonce:           nonce,
			signedAmount:    signedAmount,
			mpeAddress:      mpeAddress,
		}, nil
	}

	logger.Info().
		Str("channel_id", filteredChannel.ChannelId.String()).
		Msg("using existing payment channel")

	filteredChannelState, err := GetChannelInfoFromService(grpc, context.Background(), serviceMetadata.URL, mpeAddress, filteredChannel.ChannelId, currentBlockNumber, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel state: %w", err)
	}

	currentSignedAmount := new(big.Int).SetBytes(filteredChannelState.GetCurrentSignedAmount())

	channelID, err := evm.EnsureChannelValidity(filteredChannel, currentSignedAmount, priceInCogs, newExpiration, opts, chans)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure channel validity: %w", err)
	}

	channelState, err := GetChannelInfoFromService(grpc, context.Background(), serviceMetadata.URL, mpeAddress, channelID, currentBlockNumber, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get final channel state: %w", err)
	}

	nonce := new(big.Int).SetBytes(channelState.GetCurrentNonce())
	if nonce == nil {
		return nil, errors.New("error while getting current nonce")
	}

	currentSignedAmount = new(big.Int).SetBytes(channelState.GetCurrentSignedAmount())
	if currentSignedAmount == nil {
		return nil, errors.New("error while getting signed amount")
	}

	increment := new(big.Int).Mul(priceInCogs, big.NewInt(int64(callCount)))
	signedAmount := new(big.Int).Add(currentSignedAmount, increment)

	// TODO: Add signedAmount vs channel amount check

	grpcClient, err := grpc.GetClient(serviceMetadata.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to get gRPC client: %w", err)
	}

	logger.Info().
		Str("channel_id", channelID.String()).
		Str("nonce", nonce.String()).
		Str("signed_amount", signedAmount.String()).
		Msg("payment channel handler created with existing channel")

	return &PaymentChannelHandler{
		ethClient:       evm,
		grpcManager:     grpc,
		serviceMetadata: serviceMetadata,
		privateKeyECDSA: privateKey,
		signerAddress:   fromAddress,
		tokenClient:     NewTokenServiceClient(grpcClient.Conn),
		channelID:       channelID,
		nonce:           nonce,
		signedAmount:    signedAmount,
		mpeAddress:      mpeAddress,
	}, nil
}

// UpdateTokenState refreshes the payment handler state and obtains a new authentication token
func (h *PaymentChannelHandler) UpdateTokenState(ctx context.Context) error {
	logger := log.With().
		Str("service_id", h.serviceMetadata.SnetID).
		Str("channel_id", h.channelID.String()).
		Logger()

	currentBlockNumber, err := h.ethClient.Client.BlockNumber(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get current block number")
		return fmt.Errorf("failed to get current block: %w", err)
	}

	// Generate claim signature for the payment channel
	claimSignature := h.generatePaymentClaimSignature()

	// Sign the claim with current block number for time-based validation
	signature := h.signWithBlockNumber(claimSignature, big.NewInt(int64(currentBlockNumber)))

	request := TokenRequest{
		ChannelId:      h.channelID.Uint64(),
		CurrentNonce:   h.nonce.Uint64(),
		SignedAmount:   h.signedAmount.Uint64(),
		Signature:      signature,
		CurrentBlock:   currentBlockNumber,
		ClaimSignature: claimSignature,
	}

	tokenReply, err := h.tokenClient.GetToken(ctx, &request)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	h.Token = tokenReply.GetToken()
	logger.Debug().
		Str("token", h.Token).
		Uint64("channel_id", tokenReply.GetChannelId()).
		Uint64("planned_amount", tokenReply.GetPlannedAmount()).
		Uint64("used_amount", tokenReply.GetUsedAmount()).
		Msg("payment handler token state updated successfully")
	return nil
}

// BuildRequestMetadata constructs gRPC metadata for service calls with payment information
func (h *PaymentChannelHandler) BuildRequestMetadata(ctx context.Context) context.Context {
	logger := log.With().
		Str("service_id", h.serviceMetadata.SnetID).
		Str("channel_id", h.channelID.String()).
		Logger()

	logger.Debug().
		Str("token", h.Token).
		Msg("building gRPC request metadata")

	// Construct metadata with payment channel information
	md := metadata.New(map[string]string{
		"snet-payment-type":       "prepaid-call",
		PaymentChannelIDHeader:    h.channelID.String(),
		PaymentChannelNonceHeader: h.nonce.String(),
		PrePaidAuthTokenHeader:    h.Token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// AvailableFreeCallCount returns the number of available free calls for this payment handler
func (h *PaymentChannelHandler) AvailableFreeCallCount() (uint64, error) {
	// Payment channel handlers do not support free calls
	return 0, nil
}

// generatePaymentClaimSignature creates a cryptographic claim signature for payment validation
func (h *PaymentChannelHandler) generatePaymentClaimSignature() []byte {
	// Construct message by concatenating payment claim components
	message := bytes.Join([][]byte{
		[]byte(PrefixInSignature),
		h.mpeAddress.Bytes(),
		util.BigIntToBytes(h.channelID),
		util.BigIntToBytes(h.nonce),
		util.BigIntToBytes(h.signedAmount),
	}, nil)
	return util.GetSignature(message, h.privateKeyECDSA)
}

// signWithBlockNumber creates a time-bounded signature by combining payment signature with current block number
func (h *PaymentChannelHandler) signWithBlockNumber(paymentSignature []byte, currentBlockNumber *big.Int) []byte {
	// Convert block number to bytes for signature
	blockBytes := math.U256Bytes(currentBlockNumber)
	// Combine payment signature with block number for temporal validation
	message := bytes.Join([][]byte{paymentSignature, blockBytes}, nil)
	return util.GetSignature(message, h.privateKeyECDSA)
}
