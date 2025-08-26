package snet

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/tensved/bobrix/contracts"

	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain/util"
	"github.com/tensved/snet-matrix-framework/pkg/db"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/protoadapt"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type Handler struct {
	DescriptorName string
	SnetID         string
	ServiceName    string
	MethodName     string
	ETH            blockchain.Ethereum
	DB             db.Service
	GRPCManager    *grpcmanager.GRPCClientManager
	InputMsg       *dynamicpb.Message
	OutputMsg      *dynamicpb.Message
}

const (
	prefixGetChannelState = "__get_channel_state"
	paymentChannelTimeout = time.Minute * 1
)

// Deprecated: MultiPartyEscrowChannel represents a multi-party escrow channel.
type MultiPartyEscrowChannel struct {
	Sender     common.Address
	Recipient  common.Address
	GroupID    [32]byte
	Value      *big.Int
	Nonce      *big.Int
	Expiration *big.Int
	Signer     common.Address
}

// Deprecated: chansToWatch represents channels for watching blockchain events.
type chansToWatch struct {
	channelOpens    chan *blockchain.MultiPartyEscrowChannelOpen
	channelExtends  chan *blockchain.MultiPartyEscrowChannelExtend
	channelAddFunds chan *blockchain.MultiPartyEscrowChannelAddFunds
	DepositFunds    chan *blockchain.MultiPartyEscrowDepositFunds
	err             chan error
}

// Deprecated: bindOpts represents binding options for blockchain operations.
type bindOpts struct {
	call     *bind.CallOpts
	transact *bind.TransactOpts
	watch    *bind.WatchOpts
	filter   *bind.FilterOpts
}

func NewHandler(descriptorName, snetID, serviceName, methodName string, inputMsg, outputMsg *dynamicpb.Message, eth blockchain.Ethereum, db db.Service, grpc *grpcmanager.GRPCClientManager) *Handler {

	return &Handler{
		DescriptorName: descriptorName,
		SnetID:         snetID,
		ServiceName:    serviceName,
		MethodName:     methodName,
		ETH:            eth,
		DB:             db,
		GRPCManager:    grpc,
		InputMsg:       inputMsg,
		OutputMsg:      outputMsg,
	}
}

func NewService(serviceDescriptor protoreflect.ServiceDescriptor, descriptorName, snetID, serviceName string, eth blockchain.Ethereum, db db.Service, grpc *grpcmanager.GRPCClientManager) *contracts.Service {
	service := &contracts.Service{
		Name:        snetID,
		Description: map[string]string{"en": serviceName},
		Methods:     make(map[string]*contracts.Method),
	}

	methods := serviceDescriptor.Methods()
	if methods != nil {
		for j := range methods.Len() {
			if methods.Get(j) != nil {
				method := &contracts.Method{}
				var inputs []contracts.Input
				var outputs []contracts.Output

				method.Name = string(methods.Get(j).Name())

				inputFields := methods.Get(j).Input().Fields()
				outputFields := methods.Get(j).Output().Fields()

				if inputFields != nil {
					for n := range inputFields.Len() {
						input := contracts.Input{
							Name:        inputFields.Get(n).JSONName(),
							Description: map[string]string{"en": fmt.Sprintf("Field: %s", inputFields.Get(n).JSONName())},
							IsRequired:  true,
						}
						inputs = append(inputs, input)
					}
				}

				if outputFields != nil {
					for n := range outputFields.Len() {
						output := contracts.Output{
							Name: outputFields.Get(n).JSONName(),
						}
						outputs = append(outputs, output)
					}
				}

				inputType := methods.Get(j).Input()
				if inputType == nil {
					log.Error().Msgf("inputType not found in service with name %s and method with name %s", serviceName, method.Name)
					return nil
				}
				outputType := methods.Get(j).Output()
				if outputType == nil {
					log.Error().Msgf("outputType not found in service with name %s and method with name %s", serviceName, method.Name)
					return nil
				}
				inputMsg := dynamicpb.NewMessage(inputType)
				outputMsg := dynamicpb.NewMessage(outputType)

				method.Inputs = inputs
				method.Outputs = outputs

				method.Handler = &contracts.Handler{
					Name: method.Name,
					Do: func(ctx contracts.HandlerContext) error {
						handler := NewHandler(descriptorName, snetID, serviceName, method.Name, inputMsg, outputMsg, eth, db, grpc)
						inputData := make(map[string]any)
						for name, input := range ctx.Inputs() {
							inputData[name] = input.Value()
						}
						response := handler.Do(inputData)
						if response.Err != nil {
							return response.Err
						}
						// Copy outputs from response to context
						for name, output := range response.Outputs {
							if ctxOutput, ok := ctx.Outputs()[name]; ok {
								ctxOutput.SetValue(output.Value())
							}
						}
						return nil
					},
				}

				service.Methods[method.Name] = method
			}
		}
	}

	return service
}

func (h *Handler) Do(inputData map[string]any) *contracts.MethodResponse {
	logger := log.With().
		Str("service_id", h.SnetID).
		Str("method", h.MethodName).
		Str("descriptor", h.DescriptorName).
		Logger()

	startTime := time.Now()

	snetService, err := h.DB.GetSnetService(h.SnetID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get service from database")
		return &contracts.MethodResponse{
			Err: err,
		}
	}

	logger.Debug().
		Str("service_url", snetService.URL).
		Str("service_name", snetService.DisplayName).
		Msg("retrieved service from database")

	privateKey, err := crypto.HexToECDSA(config.Blockchain.AdminPrivateKey)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse private key")
		return &contracts.MethodResponse{
			Err: err,
		}
	}

	protoFiles, err := getProtoFilesForService(snetService)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get proto files")
		return &contracts.MethodResponse{
			Err: err,
		}
	}

	logger.Debug().
		Int("proto_files_count", len(protoFiles)).
		Msg("retrieved proto files")

	paymentManager := NewPaymentManager(h.ETH, h.DB, h.GRPCManager, privateKey, protoFiles)

	result, err := paymentManager.ExecuteCall(context.Background(), snetService, h.MethodName, inputData)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute service call")
		return &contracts.MethodResponse{
			Err: err,
		}
	}

	var response interface{}
	if resultMap, ok := result.(map[string]interface{}); ok {
		if resp, exists := resultMap["response"]; exists {
			response = resp
		} else {
			response = result
		}
	} else {
		response = result
	}

	output := contracts.Output{
		Name: "answer",
		Type: contracts.IOTypeText,
	}
	output.SetValue(response)

	duration := time.Since(startTime)
	logger.Info().
		Dur("duration", duration).
		Interface("response", response).
		Msg("service request completed successfully")

	return &contracts.MethodResponse{
		Outputs: map[string]contracts.Output{
			"answer": output,
		},
	}
}

// Deprecated: getChannelStateFromBlockchain retrieves channel state from the blockchain.
func (h *Handler) getChannelStateFromBlockchain(channelID *big.Int) (channel *MultiPartyEscrowChannel, ok bool, err error) {
	logger := log.With().
		Str("channel_id", channelID.String()).
		Logger()

	logger.Debug().Msg("getting channel state from blockchain")

	ch, err := h.ETH.MPE.Channels(nil, channelID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get channel from blockchain")
		return nil, false, err
	}

	var zeroAddress = common.Address{}
	if ch.Sender == zeroAddress {
		logger.Warn().Msg("channel has zero sender address")
		return nil, false, errors.New("incorrect sender of channel")
	}

	channel = &MultiPartyEscrowChannel{
		Sender:     ch.Sender,
		Recipient:  ch.Recipient,
		GroupID:    ch.GroupId,
		Value:      ch.Value,
		Nonce:      ch.Nonce,
		Expiration: ch.Expiration,
		Signer:     ch.Signer,
	}

	logger.Debug().
		Str("sender", ch.Sender.Hex()).
		Str("recipient", ch.Recipient.Hex()).
		Str("value", ch.Value.String()).
		Str("nonce", ch.Nonce.String()).
		Str("expiration", ch.Expiration.String()).
		Msg("retrieved channel state from blockchain")

	return channel, true, nil
}

// Deprecated: getChannelStateFromDaemon retrieves channel state from daemon.
func (h *Handler) getChannelStateFromDaemon(serviceURL string, channelID, lastBlockNumber *big.Int, privateKeyECDSA *ecdsa.PrivateKey) (*ChannelStateReply, error) {
	logger := log.With().
		Str("service_url", serviceURL).
		Str("channel_id", channelID.String()).
		Str("block_number", lastBlockNumber.String()).
		Logger()

	logger.Debug().Msg("getting channel state from daemon")

	grpcService, err := h.GRPCManager.GetClient(serviceURL)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get gRPC client")
		return nil, err
	}

	if grpcService.Conn == nil {
		logger.Error().Msg("gRPC service connection is nil")
		return nil, errors.New("grpc service connection is nil")
	}

	client := NewPaymentChannelStateServiceClient(grpcService.Conn)

	signature := h.getSignatureToGetChannelStateFromDaemon(channelID, lastBlockNumber, privateKeyECDSA)

	request := &ChannelStateRequest{
		ChannelId:    util.BigIntToBytes(channelID),
		Signature:    signature,
		CurrentBlock: lastBlockNumber.Uint64(),
	}

	reply, err := client.GetChannelState(context.Background(), request)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, errors.New("channel state reply is nil")
	}
	log.Debug().Msgf("channel state reply: %v", reply)
	return reply, nil
}

// Deprecated: getSignatureToGetChannelStateFromDaemon creates signature for channel state request.
func (h *Handler) getSignatureToGetChannelStateFromDaemon(channelID, lastBlockNumber *big.Int, privateKeyECDSA *ecdsa.PrivateKey) []byte {
	message := bytes.Join([][]byte{
		[]byte(prefixGetChannelState),
		h.ETH.MPEAddress.Bytes(),
		util.BigIntToBytes(channelID),
		math.U256Bytes(lastBlockNumber),
	}, nil)

	signature := util.GetSignature(message, privateKeyECDSA)
	return signature
}

// Deprecated: getMetadataToInvokeMethod creates metadata for method invocation.
func (h *Handler) getMetadataToInvokeMethod(channelID, nonce *big.Int, totalAmount int64, privateKeyECDSA *ecdsa.PrivateKey) metadata.MD {
	message := bytes.Join([][]byte{
		[]byte(blockchain.PrefixInSignature),
		h.ETH.MPEAddress.Bytes(),
		util.BigIntToBytes(channelID),
		util.BigIntToBytes(nonce),
		util.BigIntToBytes(big.NewInt(totalAmount)),
	}, nil)

	signature := util.GetSignature(message, privateKeyECDSA)

	md := metadata.New(map[string]string{
		"snet-payment-type":                  "escrow",
		"snet-payment-channel-id":            strconv.FormatInt(channelID.Int64(), 10),
		"snet-payment-channel-nonce":         strconv.FormatInt(nonce.Int64(), 10),
		"snet-payment-channel-amount":        strconv.Itoa(int(totalAmount)),
		"snet-payment-channel-signature-bin": string(signature),
	})
	return md
}

// Deprecated: callMethod calls a service method via gRPC.
func (h *Handler) callMethod(serviceURL string, inputProto interface{}, outputMsg interface{}, md metadata.MD) (string, error) {
	log.Debug().Msgf("call method: serviceURL=%v, inputProto=%v", serviceURL, inputProto)
	client, err := h.GRPCManager.GetClient(serviceURL)
	if err != nil {
		return "", err
	}

	endpoint := "/" + h.DescriptorName + "." + h.ServiceName + "/" + h.MethodName
	log.Debug().Msgf("endpoint: %s", endpoint)

	err = client.CallMethod(endpoint, inputProto, outputMsg, md)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}

	outputProto := protoadapt.MessageV2Of(h.OutputMsg)
	if outputProto == nil {
		return "", errors.New("failed to convert output to proto message")
	}

	log.Debug().Msgf("response from snet service: %+v", outputProto)

	jsonBytes, err := protojson.Marshal(outputProto)
	if err != nil {
		return "", errors.New("failed to marshal message")
	}
	jsonString := string(jsonBytes)
	return jsonString, nil
}

// Deprecated: getMPEBalance retrieves MPE balance.
func (h *Handler) getMPEBalance(callOpts *bind.CallOpts) (*big.Int, error) {
	mpeBalance, err := h.ETH.MPE.Balances(callOpts, callOpts.From)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("MPE balance: %v", mpeBalance)
	return mpeBalance, nil
}

// Deprecated: fillInputs fills input parameters for the handler.
func (h *Handler) fillInputs(params map[string]any) {
	cnt := 0
	for fieldName, param := range params {
		log.Debug().Msgf("set param %v with value %v to input field with name %s", param, param, fieldName)

		kind := h.InputMsg.Descriptor().Fields().Get(cnt).Kind().String()
		if kind == "float" {
			var value float64
			switch v := param.(type) {
			case string:
				var err error
				value, err = strconv.ParseFloat(v, 64)
				if err != nil {
					log.Error().Err(err).Msgf("failed to convert value %s to float64", v)
					return
				}
			case float64:
				value = v
			case float32:
				value = float64(v)
			case int:
				value = float64(v)
			case int64:
				value = float64(v)
			default:
				log.Error().Msgf("unsupported type for float field: %T", param)
				return
			}
			h.InputMsg.Set(h.InputMsg.Descriptor().Fields().ByName(protoreflect.Name(fieldName)), protoreflect.ValueOfFloat32(float32(value)))
			continue
		}
		if kind == "string" {
			var strValue string
			switch v := param.(type) {
			case string:
				strValue = v
			case float64:
				strValue = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				strValue = strconv.Itoa(v)
			case int64:
				strValue = strconv.FormatInt(v, 10)
			default:
				strValue = fmt.Sprintf("%v", v)
			}
			h.InputMsg.Set(h.InputMsg.Descriptor().Fields().ByName(protoreflect.Name(fieldName)), protoreflect.ValueOfString(strValue))
			continue
		}
		cnt++
	}
}

// Deprecated: filterChannels filters payment channels.
func (h *Handler) filterChannels(senders, recipients []common.Address, groupIDs [][32]byte, filterOpts *bind.FilterOpts) (*blockchain.MultiPartyEscrowChannelOpen, error) {
	channelOpenIterator, err := h.ETH.MPE.FilterChannelOpen(filterOpts, senders, recipients, groupIDs)
	if err != nil {
		return nil, err
	}

	defer func(channelOpenIterator *blockchain.MultiPartyEscrowChannelOpenIterator) {
		err = channelOpenIterator.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(channelOpenIterator)

	var event *blockchain.MultiPartyEscrowChannelOpen
	var filteredEvent *blockchain.MultiPartyEscrowChannelOpen

	for channelOpenIterator.Next() {
		event = channelOpenIterator.Event
		if event.Sender == senders[0] && event.Signer == senders[0] && event.Recipient == recipients[0] && event.GroupId == groupIDs[0] {
			log.Debug().Msgf("filtered event: %+v", event)
			filteredEvent = channelOpenIterator.Event
		}
	}

	if err = channelOpenIterator.Error(); err != nil {
		return nil, err
	}

	return filteredEvent, nil
}

// Deprecated: getFromAddressAndPrivateKeyECDSA retrieves address and private key.
func (h *Handler) getFromAddressAndPrivateKeyECDSA() (common.Address, *ecdsa.PrivateKey, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(config.Blockchain.AdminPrivateKey)
	if err != nil {
		return common.Address{}, nil, err
	}

	publicKey := privateKeyECDSA.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, nil, errors.New("failed to get public key")
	}

	// sender/signer address
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress, privateKeyECDSA, nil
}

// Deprecated: watchChannelOpen watches for channel open events.
func (h *Handler) watchChannelOpen(watchOpts *bind.WatchOpts, channelOpens chan *blockchain.MultiPartyEscrowChannelOpen, errChan chan error, senders, recipients []common.Address, groupIDs [][32]byte) {
	_, err := h.ETH.MPE.WatchChannelOpen(watchOpts, channelOpens, senders, recipients, groupIDs)
	if err != nil {
		errChan <- err
		return
	}
}

// Deprecated: watchDepositFunds watches for deposit fund events.
func (h *Handler) watchDepositFunds(watchOpts *bind.WatchOpts, depositFundsChan chan *blockchain.MultiPartyEscrowDepositFunds, errChan chan error, senders []common.Address) {
	_, err := h.ETH.MPE.WatchDepositFunds(watchOpts, depositFundsChan, senders)
	if err != nil {
		errChan <- err
		return
	}
}

// Deprecated: watchChannelAddFunds watches for channel add funds events.
func (h *Handler) watchChannelAddFunds(watchOpts *bind.WatchOpts, channelAddFundsChan chan *blockchain.MultiPartyEscrowChannelAddFunds, errChan chan error, channelIDs []*big.Int) {
	_, err := h.ETH.MPE.WatchChannelAddFunds(watchOpts, channelAddFundsChan, channelIDs)
	if err != nil {
		errChan <- err
		return
	}
}

// Deprecated: watchChannelExtend watches for channel extend events.
func (h *Handler) watchChannelExtend(watchOpts *bind.WatchOpts, channelExtendsChan chan *blockchain.MultiPartyEscrowChannelExtend, errChan chan error, channelIDs []*big.Int) {
	_, err := h.ETH.MPE.WatchChannelExtend(watchOpts, channelExtendsChan, channelIDs)
	if err != nil {
		errChan <- err
		return
	}
}

// Deprecated: selectPaymentChannel selects appropriate payment channel.
func (h *Handler) selectPaymentChannel(openedChannel *blockchain.MultiPartyEscrowChannelOpen, hasSufficientBalance bool, chans *chansToWatch, opts *bindOpts, senders, recipients []common.Address, groupIDs [][32]byte, price, newExpiration *big.Int) (*big.Int, *big.Int, error) {
	abiString := "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	tokenABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, nil, err
	}

	tokenAddress, err := h.ETH.MPE.Token(&bind.CallOpts{})
	if err != nil {
		return nil, nil, err
	}

	spenderAddress := h.ETH.MPEAddress

	client := h.ETH.Client

	tokenContract := bind.NewBoundContract(tokenAddress, tokenABI, client, client, client)

	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10)

	privateKeyECDSA, err := crypto.HexToECDSA(config.Blockchain.AdminPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	chainID, err := h.ETH.Client.NetworkID(context.Background())
	if err != nil {
		return nil, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)
	if err != nil {
		return nil, nil, err
	}

	var channelID *big.Int
	var nonce *big.Int
	var channelIDs []*big.Int
	if openedChannel == nil {
		if hasSufficientBalance {
			log.Debug().Msg("MPE has sufficient balance")
			go h.watchChannelOpen(opts.watch, chans.channelOpens, chans.err, senders, recipients, groupIDs)

			channel, err := h.ETH.MPE.OpenChannel(util.EstimateGas(opts.transact), senders[0], recipients[0], groupIDs[0], price, newExpiration)
			if err != nil {
				return nil, nil, err
			}
			log.Debug().Msgf("opened channel: %v", channel.Hash().Hex())

			channelID = util.WaitingForChannelToOpen(chans.channelOpens, chans.err, paymentChannelTimeout)
		} else {
			log.Debug().Msg("MPE does not have sufficient balance")
			go h.watchChannelOpen(opts.watch, chans.channelOpens, chans.err, senders, recipients, groupIDs)
			go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
			var channel *types.Transaction
			channel, err = h.ETH.MPE.DepositAndOpenChannel(util.EstimateGas(opts.transact), senders[0], recipients[0], groupIDs[0], price, newExpiration)
			if err != nil {
				log.Error().Err(err)
				if strings.Contains(err.Error(), "execution reverted: ERC20: transfer amount exceeds allowance") {
					tx, err := tokenContract.Transact(auth, "approve", spenderAddress, amount)
					if err != nil {
						return nil, nil, err
					}

					log.Info().Msgf("Transaction approve sent: %s", tx.Hash().Hex())

					_, err = waitForTransaction(h.ETH.Client, tx)
					if err != nil {
						return nil, nil, err
					}

					channel, err = h.ETH.MPE.DepositAndOpenChannel(util.EstimateGas(opts.transact), senders[0], recipients[0], groupIDs[0], price, newExpiration)
					if err != nil {
						log.Error().Err(err)
						return nil, nil, err
					}
				}
			}
			log.Debug().Msgf("deposited amount: %v", price)
			log.Debug().Msgf("opened channel: %v", channel.Hash().Hex())
			channelID = util.WaitingForChannelToOpen(chans.channelOpens, chans.err, paymentChannelTimeout)
			deposited := util.WaitingToDepositFundsToMPE(chans.DepositFunds, chans.err, paymentChannelTimeout)
			if !deposited {
				return nil, nil, errors.New("failed to deposit funds to MPE")
			}
		}

		nonce = big.NewInt(0)
	} else {
		channelID = openedChannel.ChannelId

		_, _, err = h.getChannelStateFromBlockchain(channelID)
		if err != nil {
			return nil, nil, err
		}

		channelIDs = []*big.Int{channelID}
		nonce = openedChannel.Nonce

		hasSufficientFunds, isValidExpiration := util.IsChannelValid(openedChannel, price, newExpiration)

		switch {
		case hasSufficientFunds && isValidExpiration:
			log.Debug().Msg("channel has sufficient funds and has not expired")
		case !hasSufficientFunds && !isValidExpiration:
			log.Debug().Msg("channel doesn't have enough funds and has expired")

			if !hasSufficientBalance {
				log.Debug().Msg("MPE does not have sufficient balance")
				go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
				txHash, err := h.ETH.MPE.Deposit(util.EstimateGas(opts.transact), price)
				if err != nil {
					return nil, nil, err
				}
				log.Debug().Msgf("deposited amount: %v", price)
				log.Debug().Msgf("deposit transaction: %v", txHash.Hash().Hex())
				deposited := util.WaitingToDepositFundsToMPE(chans.DepositFunds, chans.err, paymentChannelTimeout)
				if !deposited {
					return nil, nil, errors.New("failed to deposit funds to MPE")
				}
			}

			go h.watchChannelAddFunds(opts.watch, chans.channelAddFunds, chans.err, channelIDs)
			go h.watchChannelExtend(opts.watch, chans.channelExtends, chans.err, channelIDs)
			extendAndAddFunds, err := h.ETH.MPE.ChannelExtendAndAddFunds(util.EstimateGas(opts.transact), openedChannel.ChannelId, newExpiration, price)
			if err != nil {
				return nil, nil, err
			}
			log.Debug().Msgf("extendAndAddFunds transaction: %+v", extendAndAddFunds)

			_ = util.WaitingForChannelToExtend(chans.channelExtends, chans.err, paymentChannelTimeout)
			_ = util.WaitingForChannelFundsToBeAdded(chans.channelAddFunds, chans.err, paymentChannelTimeout)

		case !hasSufficientFunds:
			log.Debug().Msg("channel does not have enough funds and has not expired")
			if !hasSufficientBalance {
				log.Debug().Msg("MPE does not have sufficient balance")
				go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
				txHash, err := h.ETH.MPE.Deposit(util.EstimateGas(opts.transact), price)
				if err != nil {
					return nil, nil, err
				}
				log.Debug().Msgf("deposited amount: %v", price)
				log.Debug().Msgf("deposit transaction: %v", txHash.Hash().Hex())
				deposited := util.WaitingToDepositFundsToMPE(chans.DepositFunds, chans.err, paymentChannelTimeout)
				if !deposited {
					return nil, nil, errors.New("failed to deposit funds to MPE")
				}
			}

			go h.watchChannelAddFunds(opts.watch, chans.channelAddFunds, chans.err, channelIDs)
			funds, err := h.ETH.MPE.ChannelAddFunds(util.EstimateGas(opts.transact), openedChannel.ChannelId, price)
			if err != nil {
				return nil, nil, err
			}
			log.Debug().Msgf("funds transaction: %s", funds.Hash().Hex())
			_ = util.WaitingForChannelFundsToBeAdded(chans.channelAddFunds, chans.err, paymentChannelTimeout)

		default:
			if !hasSufficientBalance {
				log.Debug().Msg("MPE does not have sufficient balance")
				go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
				txHash, err := h.ETH.MPE.Deposit(util.EstimateGas(opts.transact), price)
				if err != nil {
					return nil, nil, err
				}
				log.Debug().Msgf("deposited amount: %v", price)
				log.Debug().Msgf("deposit transaction: %v", txHash.Hash().Hex())
				deposited := util.WaitingToDepositFundsToMPE(chans.DepositFunds, chans.err, paymentChannelTimeout)
				if !deposited {
					return nil, nil, errors.New("failed to deposit funds to MPE")
				}
			}

			go h.watchChannelExtend(opts.watch, chans.channelExtends, chans.err, channelIDs)
			go h.watchChannelAddFunds(opts.watch, chans.channelAddFunds, chans.err, channelIDs)
			extendAndAddFunds, err := h.ETH.MPE.ChannelExtendAndAddFunds(util.EstimateGas(opts.transact), openedChannel.ChannelId, newExpiration, price)
			if err != nil {
				return nil, nil, err
			}
			log.Debug().Msgf("extendAndAddFunds transaction: %+v", extendAndAddFunds.Hash().Hex())

			_ = util.WaitingForChannelToExtend(chans.channelExtends, chans.err, paymentChannelTimeout)
			_ = util.WaitingForChannelFundsToBeAdded(chans.channelAddFunds, chans.err, paymentChannelTimeout)
		}
	}

	return channelID, nonce, nil
}

// Deprecated: getLastBlockNumber retrieves the last block number.
func (h *Handler) getLastBlockNumber() *big.Int {
	header, err := h.ETH.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to get last block number")
		return nil
	}
	return header.Number
}

// Deprecated: waitForTransaction waits for transaction confirmation.
func waitForTransaction(client *ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
	ctx := context.Background()
	txHash := tx.Hash()

	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if errors.Is(err, ethereum.NotFound) {
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			return nil, err
		}

		return receipt, nil
	}
}
