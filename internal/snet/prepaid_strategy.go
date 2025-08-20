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

// PrepaidStrategy implements prepaid call strategy
type PrepaidStrategy struct {
	evmClient       blockchain.Ethereum
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

// NewPrepaidStrategy creates a new prepaid call strategy
func NewPrepaidStrategy(evm blockchain.Ethereum, grpc *grpcmanager.GRPCClientManager, serviceMetadata *db.SnetService, privateKey *ecdsa.PrivateKey, callCount uint64, database db.Service) (Strategy, error) {
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
			_, err := GetChannelStateFromDaemon(grpc, context.Background(), serviceMetadata.URL, mpeAddress, filteredChannel.ChannelId, currentBlockNumber, privateKey)
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
			channelState, err = GetChannelStateFromDaemon(grpc, context.Background(), serviceMetadata.URL, mpeAddress, channelID, currentBlockNumber, privateKey)
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
			Msg("prepaid strategy created with new channel")

		return &PrepaidStrategy{
			evmClient:       evm,
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

	filteredChannelState, err := GetChannelStateFromDaemon(grpc, context.Background(), serviceMetadata.URL, mpeAddress, filteredChannel.ChannelId, currentBlockNumber, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel state: %w", err)
	}

	currentSignedAmount := new(big.Int).SetBytes(filteredChannelState.GetCurrentSignedAmount())

	channelID, err := evm.EnsureChannelValidity(filteredChannel, currentSignedAmount, priceInCogs, newExpiration, opts, chans)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure channel validity: %w", err)
	}

	channelState, err := GetChannelStateFromDaemon(grpc, context.Background(), serviceMetadata.URL, mpeAddress, channelID, currentBlockNumber, privateKey)
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
		Msg("prepaid strategy created with existing channel")

	return &PrepaidStrategy{
		evmClient:       evm,
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

// Refresh updates the prepaid strategy state (gets token)
func (p *PrepaidStrategy) Refresh(ctx context.Context) error {
	logger := log.With().
		Str("service_id", p.serviceMetadata.SnetID).
		Str("channel_id", p.channelID.String()).
		Logger()

	currentBlockNumber, err := p.evmClient.Client.BlockNumber(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get current block number")
		return fmt.Errorf("failed to get current block: %w", err)
	}

	claimSignature := p.getClaimSignature()

	signature := p.getSignature(claimSignature, big.NewInt(int64(currentBlockNumber)))

	request := TokenRequest{
		ChannelId:      p.channelID.Uint64(),
		CurrentNonce:   p.nonce.Uint64(),
		SignedAmount:   p.signedAmount.Uint64(),
		Signature:      signature,
		CurrentBlock:   currentBlockNumber,
		ClaimSignature: claimSignature,
	}

	tokenReply, err := p.tokenClient.GetToken(ctx, &request)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	p.Token = tokenReply.GetToken()
	logger.Debug().
		Str("token", p.Token).
		Uint64("channel_id", tokenReply.GetChannelId()).
		Uint64("planned_amount", tokenReply.GetPlannedAmount()).
		Uint64("used_amount", tokenReply.GetUsedAmount()).
		Msg("prepaid strategy refreshed successfully")
	return nil
}

// GRPCMetadata creates metadata for gRPC call
func (p *PrepaidStrategy) GRPCMetadata(ctx context.Context) context.Context {
	logger := log.With().
		Str("service_id", p.serviceMetadata.SnetID).
		Str("channel_id", p.channelID.String()).
		Logger()

	logger.Debug().
		Str("token", p.Token).
		Msg("creating gRPC metadata")

	md := metadata.New(map[string]string{
		"snet-payment-type":       "prepaid-call",
		PaymentChannelIDHeader:    p.channelID.String(),
		PaymentChannelNonceHeader: p.nonce.String(),
		PrePaidAuthTokenHeader:    p.Token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// GetFreeCallsAvailable returns the number of available free calls
func (p *PrepaidStrategy) GetFreeCallsAvailable() (uint64, error) {
	return 0, nil
}

// getClaimSignature creates claim signature
func (p *PrepaidStrategy) getClaimSignature() []byte {
	message := bytes.Join([][]byte{
		[]byte(PrefixInSignature),
		p.mpeAddress.Bytes(),
		util.BigIntToBytes(p.channelID),
		util.BigIntToBytes(p.nonce),
		util.BigIntToBytes(p.signedAmount),
	}, nil)
	return util.GetSignature(message, p.privateKeyECDSA)
}

// getSignature creates signature with current block
func (p *PrepaidStrategy) getSignature(mpeSignature []byte, currentBlockNumber *big.Int) []byte {
	blockBytes := math.U256Bytes(currentBlockNumber)
	message := bytes.Join([][]byte{mpeSignature, blockBytes}, nil)
	return util.GetSignature(message, p.privateKeyECDSA)
}
