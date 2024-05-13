package lib

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"math/big"
	"matrix-ai-framework/internal/config"
	"matrix-ai-framework/internal/grpc_manager"
	"matrix-ai-framework/pkg/blockchain"
	"matrix-ai-framework/pkg/blockchain/util"
	"matrix-ai-framework/pkg/db"
	"net/url"
	"strconv"
	"strings"
)

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

func NewSnetHandler(snetID, serviceName, methodName string, inputMsg, outputMsg *dynamicpb.Message) *SnetHandler {

	return &SnetHandler{
		SnetID:      snetID,
		ServiceName: serviceName,
		MethodName:  methodName,
		eth:         blockchain.Init(),
		db:          db.New(),
		grpcManager: grpc_manager.NewGRPCClientManager(),
		InputMsg:    inputMsg,
		OutputMsg:   outputMsg,
	}
}

func (h *SnetHandler) Call(c *MContext) {
	log.Debug().Msgf("snetID: %v", h.SnetID)
	log.Debug().Msgf("serviceName: %v", h.ServiceName)
	log.Debug().Msgf("methodName: %v", h.MethodName)

	snetService, err := h.db.GetSnetService(h.SnetID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get snet service")
		return
	}
	log.Debug().Msgf("snetService: %+v", snetService)

	group, _ := h.db.GetSnetOrgGroup(snetService.GroupID)
	log.Debug().Msgf("group: %+v", group)
	log.Debug().Msgf("groupID from group: %s", group.GroupID)
	log.Debug().Msgf("groupID from snet service: %s", snetService.GroupID)

	decodedGroupID, err := base64.StdEncoding.DecodeString(group.GroupID)
	var groupID [32]byte
	copy(groupID[:], decodedGroupID)
	log.Debug().Msgf("groupID in bytes: %v", groupID)

	recipient := common.HexToAddress(group.PaymentAddress)
	log.Debug().Msgf("recipient: %v", recipient)

	// getting ecds pub priv keys
	privateKeyECDSA, err := crypto.HexToECDSA(config.Blockchain.PrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get private key")
	}
	log.Debug().Msgf("privateKeyECDSA: %+v", privateKeyECDSA)

	publicKey := privateKeyECDSA.Public()
	log.Debug().Msgf("publicKey: %+v", publicKey)

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal().Msg("Failed to get public key")
	}
	log.Debug().Msgf("publicKeyECDSA: %+v", publicKeyECDSA)

	// signer address
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Debug().Msgf("fromAddress: %v", fromAddress)

	chainID, _ := strconv.Atoi(config.Blockchain.ChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(int64(chainID)))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create transactor")
	}

	nextChannelID, err := h.eth.MPE.NextChannelId(&bind.CallOpts{})
	if err != nil {
		return
	}
	log.Debug().Msgf("Next channel id: %v", nextChannelID)
	header, err := h.eth.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block number")
	}
	blockNumber := header.Number
	log.Debug().Msgf("Last block number: %v", blockNumber)
	expiration := new(big.Int).Add(blockNumber, big.NewInt(10000000))
	channel, err := h.eth.MPE.DepositAndOpenChannel(util.EstimateGas(auth), fromAddress, recipient, groupID, big.NewInt(10000000), expiration)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to deposit and open payment channel")
	}
	log.Debug().Msgf("Opened channel: %v", channel.Hash().Hex())

	target, err := removeProtocol(snetService.URL)
	log.Info().Msgf("Target: %v", target)
	client, err := h.grpcManager.GetClient(target)
	if err != nil {
		log.Error().Err(err)
	}

	inputMsg := h.InputMsg
	outputMsg := h.OutputMsg

	cnt := 0
	for fieldName, param := range c.Params {
		log.Info().Msgf("Set param %s with value %s to input field with name %s", param, protoreflect.ValueOf(param).String(), fieldName)

		kind := inputMsg.Descriptor().Fields().Get(cnt).Kind().String()
		if kind == "float" {
			value, err := strconv.ParseFloat(param.(string), 32)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to convert value %s to float32", param.(string))
				return
			}
			inputMsg.Set(inputMsg.Descriptor().Fields().ByName(protoreflect.Name(fieldName)), protoreflect.ValueOfFloat32(float32(value)))
			continue
		}
		if kind == "string" {
			inputMsg.Set(inputMsg.Descriptor().Fields().ByName(protoreflect.Name(fieldName)), protoreflect.ValueOfString(param.(string)))
			continue
		}
		cnt++
	}

	inputProto := proto.MessageV2(inputMsg)

	price := snetService.Price

	message := bytes.Join([][]byte{
		[]byte(blockchain.PrefixInSignature),                      // prefix
		h.eth.MPEAddress.Bytes(),                                  // mpe address
		util.BigIntToBytes(big.NewInt(nextChannelID.Int64() - 1)), // channel id
		util.BigIntToBytes(big.NewInt(0)),                         // nonce
		util.BigIntToBytes(big.NewInt(int64(price))),              // amount
	}, nil)

	signature := util.GetSignature(message, privateKeyECDSA)

	md := metadata.New(map[string]string{
		"snet-payment-type":                  "escrow",
		"snet-payment-channel-id":            strconv.FormatInt(nextChannelID.Int64()-1, 10),
		"snet-payment-channel-nonce":         "0",
		"snet-payment-channel-amount":        strconv.Itoa(price),
		"snet-payment-channel-signature-bin": string(signature),
	})

	endpoint := "/" + h.ServiceName + "/" + h.MethodName
	err = client.CallMethod(endpoint, inputProto, outputMsg, md)
	if err != nil {
		log.Error().Err(err).Msg("Failed to call method")
		return
	}

	outputProto := proto.MessageV2(outputMsg)
	if outputProto == nil {
		log.Error().Err(err).Msg("Failed to convert output to proto message")
	}
	log.Debug().Msgf("Response from snet service: %v", outputProto)

	jsonBytes, err := protojson.Marshal(outputProto)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to marshal message: %v", err)
	}
	jsonString := string(jsonBytes)
	log.Info().Msgf("jsonString from server for matrix: %v", jsonString)

	c.Result <- jsonString

	return
}

func NewSNETService(serviceDescriptor protoreflect.ServiceDescriptor, snetID, serviceName string) *AIService {
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
					}, SnetServiceHandler(snetID, string(serviceDescriptor.FullName()), methodName, inputMsg, outputMsg)),
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

func SnetServiceHandler(snetID, serviceName, methodName string, inputMsg, outputMsg *dynamicpb.Message) *SnetHandler {
	return NewSnetHandler(snetID,
		serviceName,
		methodName,
		inputMsg,
		outputMsg,
	)
}

func removeProtocol(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(rawURL, parsedURL.Scheme+"://"), nil
}
