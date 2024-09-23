package snet

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
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
	escrow "github.com/tensved/snet-matrix-framework/generated"
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
	"math/big"
	"strconv"
	"strings"
	"time"
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

type MultiPartyEscrowChannel struct {
	Sender     common.Address
	Recipient  common.Address
	GroupID    [32]byte
	Value      *big.Int
	Nonce      *big.Int
	Expiration *big.Int
	Signer     common.Address
}

type chansToWatch struct {
	channelOpens    chan *blockchain.MultiPartyEscrowChannelOpen
	channelExtends  chan *blockchain.MultiPartyEscrowChannelExtend
	channelAddFunds chan *blockchain.MultiPartyEscrowChannelAddFunds
	DepositFunds    chan *blockchain.MultiPartyEscrowDepositFunds
	err             chan error
}

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
		Description: serviceName,
		Methods:     make(map[string]*contracts.Method),
	}

	methods := serviceDescriptor.Methods()
	if methods != nil {
		for j := range methods.Len() {
			if methods.Get(j) != nil {
				method := &contracts.Method{}
				var inputs []*contracts.Input
				var outputs []*contracts.Output

				method.Name = string(methods.Get(j).Name())

				inputFields := methods.Get(j).Input().Fields()
				outputFields := methods.Get(j).Output().Fields()

				if inputFields != nil {
					for n := range inputFields.Len() {
						input := &contracts.Input{
							Name:        inputFields.Get(n).JSONName(),
							Description: fmt.Sprintf("Field: %s", inputFields.Get(n).JSONName()),
							IsRequired:  true,
						}
						inputs = append(inputs, input)
					}
				}

				if outputFields != nil {
					for n := range outputFields.Len() {
						output := &contracts.Output{
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

				method.Handler = NewHandler(descriptorName, snetID, serviceName, method.Name, inputMsg, outputMsg, eth, db, grpc)

				service.Methods[method.Name] = method
			}
		}
	}

	return service
}

func (h *Handler) Do(inputData map[string]any) *contracts.MethodResponse {
	snetService, err := h.DB.GetSnetService(h.SnetID)
	if err != nil {
		log.Error().Err(err)
		return &contracts.MethodResponse{
			Error: err,
		}
	}

	bigIntPriceInCogs := big.NewInt(int64(snetService.Price))

	group, _ := h.DB.GetSnetOrgGroup(snetService.GroupID)

	groupID, err := util.DecodePaymentGroupID(group.GroupID)
	if err != nil {
		return &contracts.MethodResponse{
			Error: err,
		}
	}

	recipient := common.HexToAddress(group.PaymentAddress)

	fromAddress, privateKeyECDSA, err := h.getFromAddressAndPrivateKeyECDSA()
	if err != nil {
		return &contracts.MethodResponse{
			Error: err,
		}
	}

	lastBlockNumber := h.getLastBlockNumber()

	opts := &bindOpts{
		call:     util.GetCallOpts(fromAddress, lastBlockNumber),
		transact: util.GetTransactOpts(privateKeyECDSA),
		watch:    util.GetWatchOpts(lastBlockNumber),
		filter:   util.GetFilterOpts(lastBlockNumber),
	}

	chans := &chansToWatch{
		channelOpens:    make(chan *blockchain.MultiPartyEscrowChannelOpen),
		channelExtends:  make(chan *blockchain.MultiPartyEscrowChannelExtend),
		channelAddFunds: make(chan *blockchain.MultiPartyEscrowChannelAddFunds),
		DepositFunds:    make(chan *blockchain.MultiPartyEscrowDepositFunds),
		err:             make(chan error),
	}

	senders := []common.Address{fromAddress}
	recipients := []common.Address{recipient}
	groupIDs := [][32]byte{groupID}

	filteredChannel, err := h.filterChannels(senders, recipients, groupIDs, opts.filter)
	if err != nil {
		return &contracts.MethodResponse{
			Error: err,
		}
	}

	mpeBalance, err := h.getMPEBalance(opts.call)
	if err != nil {
		return &contracts.MethodResponse{
			Error: err,
		}
	}
	hasSufficientBalance := mpeBalance.Cmp(bigIntPriceInCogs) >= 0

	newExpiration := util.GetNewExpiration(lastBlockNumber, group.PaymentExpirationThreshold)

	channelID, _, err := h.selectPaymentChannel(filteredChannel, hasSufficientBalance, chans, opts, senders, recipients, groupIDs, bigIntPriceInCogs, newExpiration)
	if err != nil {
		log.Error().Err(err)
		return &contracts.MethodResponse{
			Error: err,
		}
	}

	channelState, err := h.getChannelStateFromDaemon(snetService.URL, channelID, lastBlockNumber, privateKeyECDSA)
	if err != nil {
		return &contracts.MethodResponse{
			Error: err,
		}
	}

	nonce := new(big.Int).SetBytes(channelState.GetCurrentNonce())
	if nonce == nil {
		return &contracts.MethodResponse{
			Error: errors.New("invalid nonce"),
		}
	}

	signedAmount := new(big.Int).SetBytes(channelState.GetCurrentSignedAmount())
	if signedAmount == nil {
		return &contracts.MethodResponse{
			Error: errors.New("invalid signed amount"),
		}
	}

	totalAmount := signedAmount.Int64() + int64(snetService.Price)

	md := h.getMetadataToInvokeMethod(channelID, nonce, totalAmount, privateKeyECDSA)

	h.fillInputs(inputData)
	inputProto := protoadapt.MessageV2Of(h.InputMsg)

	result, err := h.callMethod(snetService.URL, inputProto, h.OutputMsg, md)
	if err != nil {
		log.Error().Err(err)
		return &contracts.MethodResponse{
			Error: err,
		}
	}

	return &contracts.MethodResponse{
		Data: map[string]any{
			"answer": result,
		},
	}
}

func (h *Handler) getChannelStateFromBlockchain(channelID *big.Int) (channel *MultiPartyEscrowChannel, ok bool, err error) {
	ch, err := h.ETH.MPE.Channels(nil, channelID)
	if err != nil {
		return nil, false, err
	}
	var zeroAddress = common.Address{}
	if ch.Sender == zeroAddress {
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

	log.Debug().Msgf("channel state in blockchain: %+v", channel)

	return channel, true, nil
}

func (h *Handler) getChannelStateFromDaemon(serviceURL string, channelID, lastBlockNumber *big.Int, privateKeyECDSA *ecdsa.PrivateKey) (*escrow.ChannelStateReply, error) {
	grpcService, err := h.GRPCManager.GetClient(serviceURL)
	if err != nil {
		return nil, err
	}

	if grpcService.Conn == nil {
		return nil, errors.New("grpc service connection is nil")
	}
	client := escrow.NewPaymentChannelStateServiceClient(grpcService.Conn)

	signature := h.getSignatureToGetChannelStateFromDaemon(channelID, lastBlockNumber, privateKeyECDSA)

	request := &escrow.ChannelStateRequest{
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

func (h *Handler) getSignatureToGetChannelStateFromDaemon(channelID, lastBlockNumber *big.Int, privateKeyECDSA *ecdsa.PrivateKey) []byte {
	message := bytes.Join([][]byte{
		[]byte(prefixGetChannelState),
		h.ETH.MPEAddress.Bytes(),
		util.BigIntToBytes(channelID),
		math.U256Bytes(big.NewInt(int64(lastBlockNumber.Uint64()))),
	}, nil)

	signature := util.GetSignature(message, privateKeyECDSA)
	return signature
}

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

func (h *Handler) getMPEBalance(callOpts *bind.CallOpts) (*big.Int, error) {
	mpeBalance, err := h.ETH.MPE.Balances(callOpts, callOpts.From)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("MPE balance: %v", mpeBalance)
	return mpeBalance, nil
}

func (h *Handler) fillInputs(params map[string]any) {
	cnt := 0
	for fieldName, param := range params {
		log.Debug().Msgf("set param %s with value %s to input field with name %s", param, protoreflect.ValueOf(param).String(), fieldName)

		kind := h.InputMsg.Descriptor().Fields().Get(cnt).Kind().String()
		if kind == "float" {
			value, err := strconv.ParseFloat(param.(string), 32)
			if err != nil {
				log.Error().Err(err).Msgf("failed to convert value %s to float32", param.(string))
				return
			}
			h.InputMsg.Set(h.InputMsg.Descriptor().Fields().ByName(protoreflect.Name(fieldName)), protoreflect.ValueOfFloat32(float32(value)))
			continue
		}
		if kind == "string" {
			h.InputMsg.Set(h.InputMsg.Descriptor().Fields().ByName(protoreflect.Name(fieldName)), protoreflect.ValueOfString(param.(string)))
			continue
		}
		cnt++
	}
}

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

func (h *Handler) watchChannelOpen(watchOpts *bind.WatchOpts, channelOpens chan *blockchain.MultiPartyEscrowChannelOpen, errChan chan error, senders, recipients []common.Address, groupIDs [][32]byte) {
	_, err := h.ETH.MPE.WatchChannelOpen(watchOpts, channelOpens, senders, recipients, groupIDs)
	if err != nil {
		errChan <- err
		return
	}
}

func (h *Handler) watchDepositFunds(watchOpts *bind.WatchOpts, depositFundsChan chan *blockchain.MultiPartyEscrowDepositFunds, errChan chan error, senders []common.Address) {
	_, err := h.ETH.MPE.WatchDepositFunds(watchOpts, depositFundsChan, senders)
	if err != nil {
		errChan <- err
		return
	}
}

func (h *Handler) watchChannelAddFunds(watchOpts *bind.WatchOpts, channelAddFundsChan chan *blockchain.MultiPartyEscrowChannelAddFunds, errChan chan error, channelIDs []*big.Int) {
	_, err := h.ETH.MPE.WatchChannelAddFunds(watchOpts, channelAddFundsChan, channelIDs)
	if err != nil {
		errChan <- err
		return
	}
}

func (h *Handler) watchChannelExtend(watchOpts *bind.WatchOpts, channelExtendsChan chan *blockchain.MultiPartyEscrowChannelExtend, errChan chan error, channelIDs []*big.Int) {
	_, err := h.ETH.MPE.WatchChannelExtend(watchOpts, channelExtendsChan, channelIDs)
	if err != nil {
		errChan <- err
		return
	}
}

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
	amount.SetString("1000000000000000000", 10) // Например, 1 токен с 18 знаками после запятой

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

func (h *Handler) getLastBlockNumber() *big.Int {
	header, err := h.ETH.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to get last block number")
		return nil
	}
	return header.Number
}

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
