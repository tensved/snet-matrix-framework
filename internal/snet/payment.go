package snet

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"maps"
	"slices"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain/util"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Strategy interface for payment strategies
type Strategy interface {
	BuildRequestMetadata(ctx context.Context) context.Context
	UpdateTokenState(ctx context.Context) error
	AvailableFreeCallCount() (uint64, error)
}

// PaymentManager manages payments and strategies
type PaymentManager struct {
	ethClient       blockchain.Ethereum
	database        db.Service
	grpcManager     *grpcmanager.GRPCClientManager
	privateKey      *ecdsa.PrivateKey
	protoFiles      map[string]string // Add proto files
	currentStrategy Strategy
}

// NewPaymentManager creates a new PaymentManager instance
func NewPaymentManager(ethClient blockchain.Ethereum, database db.Service, grpcManager *grpcmanager.GRPCClientManager, privateKey *ecdsa.PrivateKey, protoFiles map[string]string) *PaymentManager {
	return &PaymentManager{
		ethClient:   ethClient,
		database:    database,
		grpcManager: grpcManager,
		privateKey:  privateKey,
		protoFiles:  protoFiles,
	}
}

// GetStrategy selects the appropriate strategy (prepaid)
func (pm *PaymentManager) GetStrategy(snetService *db.SnetService) (Strategy, error) {
	logger := log.With().
		Str("service_id", snetService.SnetID).
		Str("service_url", snetService.URL).
		Logger()

	logger.Debug().Msg("selecting payment strategy")

	logger.Debug().Msg("using payment channel strategy")
	return pm.getPaymentChannelHandler(snetService)
}

// getPaymentChannelHandler creates a payment channel handler
func (pm *PaymentManager) getPaymentChannelHandler(snetService *db.SnetService) (Strategy, error) {
	logger := log.With().
		Str("service_id", snetService.SnetID).
		Logger()

	logger.Debug().Msg("creating payment channel handler")
	paymentHandler, err := NewPaymentChannelHandler(pm.ethClient, pm.grpcManager, snetService, pm.privateKey, 1, pm.database) // callCount = 1
	if err != nil {
		logger.Error().Err(err).Msg("failed to create payment channel handler")
		return nil, fmt.Errorf("failed to create payment channel handler: %w", err)
	}
	logger.Debug().Msg("payment channel handler created successfully")
	return paymentHandler, nil
}

// ExecuteCall executes a service call with automatic strategy selection
func (pm *PaymentManager) ExecuteCall(ctx context.Context, snetService *db.SnetService, methodName string, inputData map[string]interface{}) (interface{}, error) {
	logger := log.With().
		Str("service_id", snetService.SnetID).
		Str("method", methodName).
		Logger()

	logger.Info().Msg("executing service call")

	strategy, err := pm.GetStrategy(snetService)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get payment strategy")
		return nil, fmt.Errorf("failed to get strategy: %w", err)
	}

	err = strategy.UpdateTokenState(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update payment handler token state")
		return nil, fmt.Errorf("failed to update strategy token state: %w", err)
	}

	result, err := pm.callService(ctx, snetService, methodName, inputData, strategy)
	if err != nil {
		logger.Error().Err(err).Msg("failed to call service")
		return nil, fmt.Errorf("failed to call service: %w", err)
	}

	logger.Info().Msg("service call completed successfully")
	return result, nil
}

// callService executes a gRPC call with the selected strategy
func (pm *PaymentManager) callService(ctx context.Context, snetService *db.SnetService, methodName string, inputData map[string]interface{}, strategy Strategy) (interface{}, error) {
	logger := log.With().
		Str("service_id", snetService.SnetID).
		Str("method", methodName).
		Logger()

	logger.Debug().Msg("preparing gRPC call")

	grpcClient, err := pm.grpcManager.GetClient(snetService.URL)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get gRPC client")
		return nil, fmt.Errorf("failed to get gRPC client: %w", err)
	}

	ctxWithMetadata := strategy.BuildRequestMetadata(ctx)
	if ctxWithMetadata == nil {
		logger.Error().Msg("failed to get gRPC metadata")
		return nil, fmt.Errorf("failed to get gRPC metadata")
	}

	inputJSON, err := json.Marshal(inputData)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal input data")
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	logger.Debug().
		Str("input_json", string(inputJSON)).
		Msg("input data prepared")

	protoFiles, err := pm.getProtoDescriptors()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get proto descriptors")
		return nil, fmt.Errorf("failed to get proto descriptors: %w", err)
	}

	fileDesc, methodDesc, err := pm.findMethod(protoFiles, methodName)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find gRPC method")
		return nil, fmt.Errorf("failed to find method: %w", err)
	}

	in := dynamicpb.NewMessage(methodDesc.Input())
	out := dynamicpb.NewMessage(methodDesc.Output())

	err = protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}.Unmarshal(inputJSON, in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal input to protobuf")
		return nil, fmt.Errorf("failed to unmarshal input: %w", err)
	}

	fullMethodName := "/" + string(fileDesc.Package()) + "." + string(methodDesc.Parent().Name()) + "/" + methodName

	logger.Debug().
		Str("method", fullMethodName).
		Msg("executing gRPC call")

	err = grpcClient.Conn.Invoke(ctxWithMetadata, fullMethodName, in, out)
	if err != nil {
		logger.Error().Err(err).Msg("failed to invoke gRPC method")
		return nil, fmt.Errorf("failed to invoke gRPC method: %w", err)
	}

	jsonBytes, err := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(out)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal response")
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &responseData); err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal response")
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := map[string]interface{}{
		"service":    snetService.SnetID,
		"method":     methodName,
		"input":      inputData,
		"url":        snetService.URL,
		"free_calls": snetService.FreeCalls,
		"price":      snetService.Price,
		"strategy":   "prepaid", // TODO: determine strategy dynamically
		"response":   responseData,
	}

	logger.Debug().
		Interface("result", result).
		Msg("gRPC call completed successfully")

	return result, nil
}

// getProtoDescriptors compiles proto files into descriptors
func (pm *PaymentManager) getProtoDescriptors() (linker.Files, error) {
	accessor := protocompile.SourceAccessorFromMap(pm.protoFiles)
	r := protocompile.WithStandardImports(&protocompile.SourceResolver{Accessor: accessor})
	compiler := protocompile.Compiler{
		Resolver:       r,
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	fds, err := compiler.Compile(context.Background(), slices.Collect(maps.Keys(pm.protoFiles))...)
	if err != nil || fds == nil {
		log.Error().Err(err).Msg("failed to compile proto files")
		return nil, fmt.Errorf("failed to compile proto files: %v", err)
	}
	return fds, nil
}

// findMethod finds the method in proto files
func (pm *PaymentManager) findMethod(files linker.Files, methodName string) (protoreflect.FileDescriptor, protoreflect.MethodDescriptor, error) {
	for _, file := range files {
		for i := 0; i < file.Services().Len(); i++ {
			service := file.Services().Get(i)
			method := service.Methods().ByName(protoreflect.Name(methodName))
			if method != nil {
				return file, method, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("method %s not found in provided proto files", methodName)
}

// GetChannelInfoFromService retrieves channel state information from the service daemon.
func GetChannelInfoFromService(grpcManager *grpcmanager.GRPCClientManager, ctx context.Context, daemonURL string, mpeAddress common.Address, channelID *big.Int, currentBlockNumber uint64, privateKey *ecdsa.PrivateKey) (*ChannelStateReply, error) {
	grpcClient, err := grpcManager.GetClient(daemonURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get gRPC client: %w", err)
	}

	message := bytes.Join([][]byte{
		[]byte("__get_channel_state"),
		mpeAddress.Bytes(),
		util.BigIntToBytes(channelID),
		math.U256Bytes(big.NewInt(int64(currentBlockNumber))),
	}, nil)

	stateClient := NewPaymentChannelStateServiceClient(grpcClient.Conn)
	stateReply, err := stateClient.GetChannelState(ctx, &ChannelStateRequest{
		ChannelId:    util.BigIntToBytes(channelID),
		Signature:    util.GetSignature(message, privateKey),
		CurrentBlock: currentBlockNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get channel state: %w", err)
	}

	return stateReply, nil
}
