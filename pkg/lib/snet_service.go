package lib

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"math/big"
	escrow "matrix-ai-framework/generated"
	"matrix-ai-framework/internal/config"
	"matrix-ai-framework/internal/grpc_manager"
	"matrix-ai-framework/pkg/blockchain"
	"matrix-ai-framework/pkg/blockchain/util"
	"matrix-ai-framework/pkg/db"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	prefixGetChannelState = "__get_channel_state"
	paymentChannelTimeout = time.Minute * 1
)

type MultiPartyEscrowChannel struct {
	Sender     common.Address
	Recipient  common.Address
	GroupId    [32]byte
	Value      *big.Int
	Nonce      *big.Int
	Expiration *big.Int
	Signer     common.Address
}

var zeroAddress = common.Address{}

type SnetHandler struct {
	SnetID      string
	ServiceName string
	MethodName  string
	Descriptor  protoreflect.FileDescriptor
	eth         blockchain.Ethereum
	db          db.Service
	grpcManager *grpc_manager.GRPCClientManager
	InputMsg    *dynamicpb.Message
	OutputMsg   *dynamicpb.Message
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

func NewSnetHandler(snetID, serviceName, methodName string, inputMsg, outputMsg *dynamicpb.Message, eth blockchain.Ethereum, db db.Service, grpc *grpc_manager.GRPCClientManager) *SnetHandler {

	return &SnetHandler{
		SnetID:      snetID,
		ServiceName: serviceName,
		MethodName:  methodName,
		eth:         eth,
		db:          db,
		grpcManager: grpc,
		InputMsg:    inputMsg,
		OutputMsg:   outputMsg,
	}
}

func (h *SnetHandler) Call(c *MContext) {
	snetService, err := h.db.GetSnetService(h.SnetID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get snet service")
		return
	}

	bigIntPriceInCogs := big.NewInt(int64(snetService.Price))

	group, _ := h.db.GetSnetOrgGroup(snetService.GroupID)

	groupID, err := decodePaymentGroupID(group.GroupID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode payment group id")
	}

	recipient := common.HexToAddress(group.PaymentAddress)

	fromAddress, privateKeyECDSA, err := h.getFromAddressAndPrivateKeyECDSA()

	lastBlockNumber := h.getLastBlockNumber()

	opts := &bindOpts{
		call:     getCallOpts(fromAddress, lastBlockNumber),
		transact: getTransactOpts(privateKeyECDSA),
		watch:    getWatchOpts(lastBlockNumber),
		filter:   getFilterOpts(lastBlockNumber),
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

	mpeBalance, err := h.getMPEBalance(opts.call)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get MPE balance")
	}
	hasSufficientBalance := mpeBalance.Cmp(bigIntPriceInCogs) >= 0

	newExpiration := getNewExpiration(lastBlockNumber, group.PaymentExpirationThreshold)

	channelID, nonce, err := h.selectPaymentChannel(filteredChannel, hasSufficientBalance, chans, opts, senders, recipients, groupIDs, bigIntPriceInCogs, newExpiration)
	if err != nil {
		log.Error().Err(err).Msg("Failed to select payment channel")
	}

	target, err := removeProtocol(snetService.URL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get grpc target")
	}

	channelState, err := h.getChannelStateFromDaemon(target, channelID, lastBlockNumber, privateKeyECDSA)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get channel state from daemon")
	}

	nonce = new(big.Int).SetBytes(channelState.CurrentNonce)
	if nonce == nil {
		log.Error().Msg("Invalid nonce")
	}

	signedAmount := new(big.Int).SetBytes(channelState.CurrentSignedAmount)
	if signedAmount == nil {
		log.Error().Msg("Invalid signed amount")
	}

	totalAmount := signedAmount.Int64() + int64(snetService.Price)

	md := h.getMetadataToInvokeMethod(channelID, nonce, totalAmount, privateKeyECDSA)

	h.fillInputs(c.Params)
	inputProto := proto.MessageV2(h.InputMsg)

	result, err := h.callMethod(target, inputProto, h.OutputMsg, md)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get result from service")
	}

	c.Result <- result

	return
}

//func AgixToCog(cogs int64) *big.Int {
//	bigIntValue := new(big.Int)
//
//	bigIntValue.SetInt64(cogs * 10000000) // agix 8 decimals
//
//	return bigIntValue
//}

func NewSNETService(serviceDescriptor protoreflect.ServiceDescriptor, snetID, serviceName string, eth blockchain.Ethereum, db db.Service, grpc *grpc_manager.GRPCClientManager) *AIService {
	name := snetID + "/" + serviceName
	log.Info().Msgf("Name of ai service: %v", name)
	snetAIService := NewAIService(name, "snet")
	methods := serviceDescriptor.Methods()
	if methods != nil {
		for j := 0; j < methods.Len(); j++ {
			var inputList []MInput
			var outputList []MOutput
			if methods.Get(j) != nil {
				inputFields := methods.Get(j).Input().Fields()
				outputFields := methods.Get(j).Output().Fields()
				methodName := string(methods.Get(j).Name())
				if inputFields != nil {
					for n := 0; n < inputFields.Len(); n++ {
						inputList = append(inputList, InputString(inputFields.Get(n).JSONName(), RequiredInput()))
					}
				}
				if outputFields != nil {
					for n := 0; n < outputFields.Len(); n++ {
						outputList = append(outputList, OutputString(outputFields.Get(n).JSONName()))
					}
				}

				inputType := methods.Get(j).Input()
				if inputType == nil {
					log.Error().Msgf("inputType not found in service with name %s and method with name %s", serviceName, methodName)
				}
				outputType := methods.Get(j).Output()
				if outputType == nil {
					log.Error().Msgf("outputType not found in service with name %s and method with name %s", serviceName, methodName)
				}
				inputMsg := dynamicpb.NewMessage(inputType)
				outputMsg := dynamicpb.NewMessage(outputType)
				snetAIService.CreateMethod(
					NewAIMethod(methodName, MethodOpts{
						Inputs:  inputList,
						Outputs: outputList,
					}, SnetServiceHandler(snetID, string(serviceDescriptor.FullName()), methodName, inputMsg, outputMsg, eth, db, grpc)),
				)

				log.Debug().Msgf("Register method: %s", methodName)
				log.Debug().Msgf("For service : %s", string(serviceDescriptor.FullName()))
				log.Debug().Msgf("With snet id : %s", snetID)
				log.Debug().Msgf("With input: %v", inputList)
				log.Debug().Msgf("With ouput: %v", outputList)

			}
		}
	}

	return snetAIService
}

func SnetServiceHandler(snetID, serviceName, methodName string, inputMsg, outputMsg *dynamicpb.Message, eth blockchain.Ethereum, db db.Service, grpc *grpc_manager.GRPCClientManager) *SnetHandler {
	return NewSnetHandler(snetID,
		serviceName,
		methodName,
		inputMsg,
		outputMsg,
		eth,
		db,
		grpc,
	)
}

func removeProtocol(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(rawURL, parsedURL.Scheme+"://"), nil
}

func waitingForChannelToOpen(channelOpens <-chan *blockchain.MultiPartyEscrowChannelOpen, errChan <-chan error) *big.Int {
	select {
	case openEvent := <-channelOpens:
		log.Debug().Msgf("Channel opened: %+v", openEvent)
		return openEvent.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msgf("Error watching for OpenChannel: %v", err)
	case <-time.After(paymentChannelTimeout):
		log.Error().Msg("Timed out waiting for OpenChannel to complete")
	}
	return nil
}

func waitingForChannelToExtend(channelExtends <-chan *blockchain.MultiPartyEscrowChannelExtend, errChan <-chan error) *big.Int {
	select {
	case extendEvent := <-channelExtends:
		log.Debug().Msgf("Channel extended: %+v", extendEvent)
		return extendEvent.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msgf("Error watching for ChannelExtend: %v", err)
	case <-time.After(paymentChannelTimeout):
		log.Error().Msg("Timed out waiting for ChannelExtend to complete")
	}
	return nil
}

func waitingForChannelFundsToBeAdded(channelAddFunds <-chan *blockchain.MultiPartyEscrowChannelAddFunds, errChan <-chan error) *big.Int {
	select {
	case addFundsEvent := <-channelAddFunds:
		log.Debug().Msgf("Channel funds added: %+v", addFundsEvent)
		return addFundsEvent.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msgf("Error watching for ChannelExtendAndAddFunds: %v", err)
	case <-time.After(paymentChannelTimeout):
		log.Error().Msg("Timed out waiting for ChannelExtendAndAddFunds to complete")
	}
	return nil
}

func waitingToDepositFundsToMPE(channelDepositFunds <-chan *blockchain.MultiPartyEscrowDepositFunds, errChan <-chan error) bool {
	select {
	case depositFundsEvent := <-channelDepositFunds:
		log.Debug().Msgf("Deposited to MPE: %+v", depositFundsEvent)
		return true
	case err := <-errChan:
		log.Error().Err(err).Msgf("Error watching for deposit: %v", err)
	case <-time.After(paymentChannelTimeout):
		log.Error().Msg("Timed out waiting for deposit to complete")
	}
	return false
}

func waitingForChannelToBeClaimed(channelClaim <-chan *blockchain.MultiPartyEscrowChannelClaim, errChan <-chan error) bool {
	select {
	case channelClaimEvent := <-channelClaim:
		log.Debug().Msgf("Channel claim: %+v", channelClaimEvent)
		return true
	case err := <-errChan:
		log.Error().Err(err).Msgf("Error watching for channel claim: %v", err)
	case <-time.After(paymentChannelTimeout):
		log.Error().Msg("Timed out waiting for channel claim to complete")
	}
	return false
}

func isChannelValid(filteredEvent *blockchain.MultiPartyEscrowChannelOpen, bigIntPrice *big.Int, newExpiration *big.Int) (bool, bool) {
	hasSufficientFunds := filteredEvent.Amount.Cmp(bigIntPrice) >= 0
	isValidExpiration := filteredEvent.Expiration.Cmp(newExpiration) > 0
	return hasSufficientFunds, isValidExpiration
}

func (h *SnetHandler) getChannelStateFromBlockchain(channelID *big.Int) (channel *MultiPartyEscrowChannel, ok bool, err error) {
	ch, err := h.eth.MPE.Channels(nil, channelID)
	if err != nil {
		log.Error().Err(err).Msg("Error while looking up for channel id in blockchain")
		return nil, false, err
	}
	if ch.Sender == zeroAddress {
		log.Error().Err(err).Msg("Unable to find channel id in blockchain")
		return nil, false, nil
	}

	channel = &MultiPartyEscrowChannel{
		Sender:     ch.Sender,
		Recipient:  ch.Recipient,
		GroupId:    ch.GroupId,
		Value:      ch.Value,
		Nonce:      ch.Nonce,
		Expiration: ch.Expiration,
		Signer:     ch.Signer,
	}

	log.Debug().Msgf("Channel found in blockchain: %+v", channel)

	return channel, true, nil
}

func (h *SnetHandler) getChannelStateFromDaemon(target string, channelID, lastBlockNumber *big.Int, privateKeyECDSA *ecdsa.PrivateKey) (*escrow.ChannelStateReply, error) {

	grpcService, err := grpc_manager.NewGRPCService(target)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create gRPC service")
		return nil, err
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
		log.Error().Err(err).Msg("Failed to get channel state")
		return nil, err
	}

	log.Debug().Msgf("Channel state reply: %+v", reply)
	return reply, nil
}

func (h *SnetHandler) getSignatureToGetChannelStateFromDaemon(channelID, lastBlockNumber *big.Int, privateKeyECDSA *ecdsa.PrivateKey) []byte {
	message := bytes.Join([][]byte{
		[]byte(prefixGetChannelState),
		h.eth.MPEAddress.Bytes(),
		util.BigIntToBytes(channelID),
		math.U256Bytes(big.NewInt(int64(lastBlockNumber.Uint64()))),
	}, nil)

	signature := util.GetSignature(message, privateKeyECDSA)
	return signature
}

func (h *SnetHandler) getMetadataToInvokeMethod(channelID, nonce *big.Int, totalAmount int64, privateKeyECDSA *ecdsa.PrivateKey) metadata.MD {
	message := bytes.Join([][]byte{
		[]byte(blockchain.PrefixInSignature),
		h.eth.MPEAddress.Bytes(),
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

func (h *SnetHandler) callMethod(target string, inputProto interface{}, outputMsg interface{}, md metadata.MD) (string, error) {
	client, err := h.grpcManager.GetClient(target)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get grpc client")
		return "", err
	}

	endpoint := "/" + h.ServiceName + "/" + h.MethodName
	err = client.CallMethod(endpoint, inputProto, outputMsg, md)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call method")
		return "", err
	}

	outputProto := proto.MessageV2(outputMsg)
	if outputProto == nil {
		log.Error().Err(err).Msg("Failed to convert output to proto message")
	}
	log.Debug().Msgf("Response from snet service: %+v", outputProto)

	jsonBytes, err := protojson.Marshal(outputProto)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to marshal message: %v", err)
	}
	jsonString := string(jsonBytes)
	log.Info().Msgf("jsonString for matrix chat message: %s", jsonString)
	return jsonString, nil
}

func (h *SnetHandler) getMPEBalance(callOpts *bind.CallOpts) (*big.Int, error) {
	mpeBalance, err := h.eth.MPE.Balances(callOpts, callOpts.From)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get MPE balance")
		return nil, err
	}
	log.Debug().Msgf("MPE balance: %v", mpeBalance)
	return mpeBalance, nil
}

func getCallOpts(fromAddress common.Address, lastBlockNumber *big.Int) *bind.CallOpts {
	return &bind.CallOpts{
		Pending:     false,
		From:        fromAddress,
		BlockNumber: lastBlockNumber,
		BlockHash:   common.Hash{},
		Context:     context.Background(),
	}
}

func getWatchOpts(lastBlockNumber *big.Int) *bind.WatchOpts {
	var startBlock uint64
	startBlock = lastBlockNumber.Uint64()
	return &bind.WatchOpts{
		Start:   &startBlock,
		Context: context.Background(),
	}
}

func getFilterOpts(lastBlockNumber *big.Int) *bind.FilterOpts {
	start := uint64(0)
	end := lastBlockNumber.Uint64()
	return &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}
}

func getTransactOpts(privateKeyECDSA *ecdsa.PrivateKey) *bind.TransactOpts {
	chainID, _ := strconv.Atoi(config.Blockchain.ChainID)
	opts, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(int64(chainID)))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create transactor")
	}
	return opts
}

func getNewExpiration(lastBlockNumber, paymentExpirationThreshold *big.Int) *big.Int {
	blockOffset := big.NewInt(240)
	defaultExpiration := new(big.Int).Add(lastBlockNumber, paymentExpirationThreshold)
	return new(big.Int).Add(defaultExpiration, blockOffset)
}

func (h *SnetHandler) getLastBlockNumber() *big.Int {
	header, err := h.eth.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get last block number")
	}
	return header.Number
}

func decodePaymentGroupID(encoded string) ([32]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode payment group id")
		return [32]byte{}, err
	}
	var groupID [32]byte
	copy(groupID[:], decoded)
	return groupID, nil
}

func (h *SnetHandler) fillInputs(params map[string]any) {
	cnt := 0
	for fieldName, param := range params {
		log.Info().Msgf("Set param %s with value %s to input field with name %s", param, protoreflect.ValueOf(param).String(), fieldName)

		kind := h.InputMsg.Descriptor().Fields().Get(cnt).Kind().String()
		if kind == "float" {
			value, err := strconv.ParseFloat(param.(string), 32)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to convert value %s to float32", param.(string))
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

func (h *SnetHandler) filterChannels(senders, recipients []common.Address, groupIDs [][32]byte, filterOpts *bind.FilterOpts) (*blockchain.MultiPartyEscrowChannelOpen, error) {
	channelOpenIterator, err := h.eth.MPE.FilterChannelOpen(filterOpts, senders, recipients, groupIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to filter channel open")
		return nil, err
	}

	defer func(channelOpenIterator *blockchain.MultiPartyEscrowChannelOpenIterator) {
		err = channelOpenIterator.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close channel open iterator")
		}
	}(channelOpenIterator)

	var event *blockchain.MultiPartyEscrowChannelOpen
	var filteredEvent *blockchain.MultiPartyEscrowChannelOpen

	for channelOpenIterator.Next() {
		event = channelOpenIterator.Event
		if event.Sender == senders[0] && event.Signer == senders[0] && event.Recipient == recipients[0] && event.GroupId == groupIDs[0] {
			log.Debug().Msgf("Filtered event: %+v", event)
			filteredEvent = channelOpenIterator.Event
		}
	}

	if err = channelOpenIterator.Error(); err != nil {
		log.Error().Err(err).Msg("Error during iteration")
		return nil, err
	}

	return filteredEvent, nil
}

func (h *SnetHandler) getFromAddressAndPrivateKeyECDSA() (common.Address, *ecdsa.PrivateKey, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(config.Blockchain.AdminPrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get private key")
		return common.Address{}, nil, err
	}

	publicKey := privateKeyECDSA.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error().Msg("Failed to get public key")
		return common.Address{}, nil, errors.New("failed to get public key")
	}

	// sender/signer address
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Debug().Msgf("fromAddress: %s", fromAddress.Hex())
	return fromAddress, privateKeyECDSA, nil
}

func (h *SnetHandler) watchChannelOpen(watchOpts *bind.WatchOpts, channelOpens chan *blockchain.MultiPartyEscrowChannelOpen, errChan chan error, senders, recipients []common.Address, groupIDs [][32]byte) {
	_, err := h.eth.MPE.WatchChannelOpen(watchOpts, channelOpens, senders, recipients, groupIDs)
	if err != nil {
		errChan <- err
		return
	}
}

func (h *SnetHandler) watchDepositFunds(watchOpts *bind.WatchOpts, DepositFundsChan chan *blockchain.MultiPartyEscrowDepositFunds, errChan chan error, senders []common.Address) {
	_, err := h.eth.MPE.WatchDepositFunds(watchOpts, DepositFundsChan, senders)
	if err != nil {
		errChan <- err
		return
	}
}

func (h *SnetHandler) watchChannelAddFunds(watchOpts *bind.WatchOpts, channelAddFundsChan chan *blockchain.MultiPartyEscrowChannelAddFunds, errChan chan error, channelIDs []*big.Int) {
	_, err := h.eth.MPE.WatchChannelAddFunds(watchOpts, channelAddFundsChan, channelIDs)
	if err != nil {
		errChan <- err
		return
	}
}

func (h *SnetHandler) watchChannelExtend(watchOpts *bind.WatchOpts, channelExtendsChan chan *blockchain.MultiPartyEscrowChannelExtend, errChan chan error, channelIDs []*big.Int) {
	_, err := h.eth.MPE.WatchChannelExtend(watchOpts, channelExtendsChan, channelIDs)
	if err != nil {
		errChan <- err
		return
	}
}

func (h *SnetHandler) selectPaymentChannel(openedChannel *blockchain.MultiPartyEscrowChannelOpen, hasSufficientBalance bool, chans *chansToWatch, opts *bindOpts, senders, recipients []common.Address, groupIDs [][32]byte, price, newExpiration *big.Int) (*big.Int, *big.Int, error) {
	var channelID *big.Int
	var nonce *big.Int
	var channelIDs []*big.Int
	if openedChannel == nil {
		if hasSufficientBalance {
			log.Debug().Msg("MPE has sufficient balance")
			go h.watchChannelOpen(opts.watch, chans.channelOpens, chans.err, senders, recipients, groupIDs)

			channel, err := h.eth.MPE.OpenChannel(util.EstimateGas(opts.transact), senders[0], recipients[0], groupIDs[0], price, newExpiration)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to open payment channel")
				return nil, nil, err
			}
			log.Debug().Msgf("Opened channel: %v", channel.Hash().Hex())

			channelID = waitingForChannelToOpen(chans.channelOpens, chans.err)
		} else {
			log.Debug().Msg("MPE does not have sufficient balance")
			go h.watchChannelOpen(opts.watch, chans.channelOpens, chans.err, senders, recipients, groupIDs)
			go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
			channel, err := h.eth.MPE.DepositAndOpenChannel(util.EstimateGas(opts.transact), senders[0], recipients[0], groupIDs[0], price, newExpiration)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to deposit and open payment channel")
				return nil, nil, err
			}
			log.Debug().Msgf("Deposited amount: %v", price)
			log.Debug().Msgf("Opened channel: %v", channel.Hash().Hex())
			channelID = waitingForChannelToOpen(chans.channelOpens, chans.err)
			deposited := waitingToDepositFundsToMPE(chans.DepositFunds, chans.err)
			if !deposited {
				log.Error().Msg("Failed to deposit funds to MPE")
				return nil, nil, err
			}
		}

		nonce = big.NewInt(0)
	} else {
		channelID = openedChannel.ChannelId

		channel, _, _ := h.getChannelStateFromBlockchain(channelID)
		log.Debug().Msgf("Opened channel: %+v", channel)

		channelIDs = []*big.Int{channelID}
		nonce = openedChannel.Nonce

		hasSufficientFunds, isValidExpiration := isChannelValid(openedChannel, price, newExpiration)

		// todo check threshold in newExpiration of existing opened channel

		switch {
		case hasSufficientFunds && isValidExpiration:
			log.Debug().Msg("The channel has sufficient funds and has not expired")
		case !hasSufficientFunds && !isValidExpiration:
			log.Debug().Msg("The channel doesn't have enough funds and has expired")

			if !hasSufficientBalance {
				log.Debug().Msg("MPE does not have sufficient balance")
				go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
				txHash, err := h.eth.MPE.Deposit(util.EstimateGas(opts.transact), price)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to deposit to MPE")
					return nil, nil, err
				}
				log.Debug().Msgf("Deposited amount: %v", price)
				log.Debug().Msgf("Deposit transaction: %v", txHash.Hash().Hex())
				deposited := waitingToDepositFundsToMPE(chans.DepositFunds, chans.err)
				if !deposited {
					log.Error().Msg("Failed to deposit funds to MPE")
					return nil, nil, err
				}
			}

			go h.watchChannelAddFunds(opts.watch, chans.channelAddFunds, chans.err, channelIDs)
			go h.watchChannelExtend(opts.watch, chans.channelExtends, chans.err, channelIDs)
			extendAndAddFunds, err := h.eth.MPE.ChannelExtendAndAddFunds(util.EstimateGas(opts.transact), openedChannel.ChannelId, newExpiration, price)
			if err != nil {
				log.Error().Err(err).Msg("Failed to extend and add funds to channel")
				return nil, nil, err
			}
			log.Debug().Msgf("extendAndAddFunds transaction: %+v", extendAndAddFunds)

			_ = waitingForChannelToExtend(chans.channelExtends, chans.err)
			_ = waitingForChannelFundsToBeAdded(chans.channelAddFunds, chans.err)

		case !hasSufficientFunds:
			log.Debug().Msg("The channel does not have enough funds and has not expired")

			if !hasSufficientBalance {
				log.Debug().Msg("MPE does not have sufficient balance")
				go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
				txHash, err := h.eth.MPE.Deposit(util.EstimateGas(opts.transact), price)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to deposit to MPE")
					return nil, nil, err
				}
				log.Debug().Msgf("Deposited amount: %v", price)
				log.Debug().Msgf("Deposit transaction: %v", txHash.Hash().Hex())
				deposited := waitingToDepositFundsToMPE(chans.DepositFunds, chans.err)
				if !deposited {
					log.Error().Msg("Failed to deposit funds to MPE")
					return nil, nil, err
				}
			}

			go h.watchChannelAddFunds(opts.watch, chans.channelAddFunds, chans.err, channelIDs)
			funds, err := h.eth.MPE.ChannelAddFunds(util.EstimateGas(opts.transact), openedChannel.ChannelId, price)
			if err != nil {
				log.Error().Err(err).Msg("Failed to add funds to channel")
			}
			log.Debug().Msgf("funds transaction: %s", funds.Hash().Hex())
			_ = waitingForChannelFundsToBeAdded(chans.channelAddFunds, chans.err)

		default:
			if !hasSufficientBalance {
				log.Debug().Msg("MPE does not have sufficient balance")
				go h.watchDepositFunds(opts.watch, chans.DepositFunds, chans.err, senders)
				txHash, err := h.eth.MPE.Deposit(util.EstimateGas(opts.transact), price)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to deposit to MPE")
					return nil, nil, err
				}
				log.Debug().Msgf("Deposited amount: %v", price)
				log.Debug().Msgf("Deposit transaction: %v", txHash.Hash().Hex())
				deposited := waitingToDepositFundsToMPE(chans.DepositFunds, chans.err)
				if !deposited {
					log.Error().Msg("Failed to deposit funds to MPE")
					return nil, nil, err
				}
			}

			go h.watchChannelExtend(opts.watch, chans.channelExtends, chans.err, channelIDs)
			go h.watchChannelAddFunds(opts.watch, chans.channelAddFunds, chans.err, channelIDs)
			extendAndAddFunds, err := h.eth.MPE.ChannelExtendAndAddFunds(util.EstimateGas(opts.transact), openedChannel.ChannelId, newExpiration, price)
			if err != nil {
				log.Error().Err(err).Msg("Failed to extend and add funds to channel")
				return nil, nil, err
			}
			log.Debug().Msgf("extendAndAddFunds transaction: %+v", extendAndAddFunds.Hash().Hex())

			_ = waitingForChannelToExtend(chans.channelExtends, chans.err)
			_ = waitingForChannelFundsToBeAdded(chans.channelAddFunds, chans.err)
		}
	}

	return channelID, nonce, nil
}
