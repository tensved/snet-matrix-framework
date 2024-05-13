package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	contracts "github.com/singnet/snet-ecosystem-contracts"
	"matrix-ai-framework/internal/config"
)

var (
	// HashPrefix32Bytes is an Ethereum signature prefix: see https://github.com/ethereum/go-ethereum/blob/bf468a81ec261745b25206b2a596eb0ee0a24a74/internal/ethapi/api.go#L361
	HashPrefix32Bytes = []byte("\x19Ethereum Signed Message:\n32")
	hashPrefix42Bytes = []byte("\x19Ethereum Signed Message:\n420x")
)

const (
	PrefixInSignature = "__MPE_claim_message"
	//Agreed constant value
	FreeCallPrefixSignature = "__prefix_free_trial"
	//Agreed constant value
	AllowedUserPrefixSignature = "__authorized_user"
)

const (
	// PaymentTypeHeader is a type of payment used to pay for an RPC call.
	// Supported types are: "escrow".
	// Note: "job" Payment type is deprecated
	PaymentTypeHeader = "snet-payment-type"
	//Client that calls the Daemon ( example can be "snet-cli","snet-dapp","snet-sdk")
	ClientTypeHeader = "snet-client-type"
	//Value is a user address , example "0x94d04332C4f5273feF69c4a52D24f42a3aF1F207"
	UserInfoHeader = "snet-user-info"
	//User Agent details set in on the server stream info
	UserAgentHeader = "user-agent"
	// PaymentChannelIDHeader is a MultiPartyEscrow contract payment channel
	// id. Value is a string containing a decimal number.
	PaymentChannelIDHeader = "snet-payment-channel-id"
	// PaymentChannelNonceHeader is a payment channel nonce value. Value is a
	// string containing a decimal number.
	PaymentChannelNonceHeader = "snet-payment-channel-nonce"
	// PaymentChannelAmountHeader is an amount of payment channel value
	// which server is authorized to withdraw after handling the RPC call.
	// Value is a string containing a decimal number.
	PaymentChannelAmountHeader = "snet-payment-channel-amount"
	// PaymentChannelSignatureHeader is a signature of the client to confirm
	// amount withdrawing authorization. Value is an array of bytes.
	PaymentChannelSignatureHeader = "snet-payment-channel-signature-bin"
	// This is useful information in the header sent in by the client
	// All clients will have this information and they need this to Sign anyways
	// When Daemon is running in the block chain disabled mode , it would use this
	// header to get the MPE address. The goal here is to keep the client oblivious to the
	// Daemon block chain enabled or disabled mode and also standardize the signatures.
	// id. Value is a string containing a decimal number.
	PaymentMultiPartyEscrowAddressHeader = "snet-payment-mpe-address"

	//Added for free call support in Daemon

	//The user Id of the person making the call
	FreeCallUserIdHeader = "snet-free-call-user-id"

	//Will be used to check if the Signature is still valid
	CurrentBlockNumberHeader = "snet-current-block-number"

	//Place holder to set the free call Auth Token issued
	FreeCallAuthTokenHeader = "snet-free-call-auth-token-bin"
	//Block number on when the Token was issued , to track the expiry of the token , which is ~ 1 Month
	FreeCallAuthTokenExpiryBlockNumberHeader = "snet-free-call-token-expiry-block"

	//Users may decide to sign upfront and make calls .Daemon generates and Auth Token
	//Users/Clients will need to use this token to make calls for the amount signed upfront.
	PrePaidAuthTokenHeader = "snet-prepaid-auth-token-bin"

	DynamicPriceDerived = "snet-derived-dynamic-price-cost"
)

type EthereumService interface {
	InitRegistry() (err error)
}

type Ethereum struct {
	Client     *ethclient.Client
	Registry   *Registry
	MPE        *MultiPartyEscrow
	MPEAddress common.Address
}

func Init() (e Ethereum) {
	var err error
	log.Debug().Any("ETH_URL", config.Blockchain.EthProviderURL).Msg("ETH_URL")
	e.Client, err = ethclient.Dial(config.Blockchain.EthProviderURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to blockchain")
	}
	err = e.InitRegistry()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to InitRegistry")
	}
	err = e.InitMPE()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to InitMPE")
	}
	return e
}

type networks map[string]struct {
	Address string `json:"address"`
}

func (eth *Ethereum) InitRegistry() (err error) {
	networksRaw := contracts.GetNetworks(contracts.Registry)
	var n networks
	err = json.Unmarshal(networksRaw, &n)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal")
		return
	}
	registryAddress := n[config.Blockchain.ChainID].Address
	log.Debug().Msgf("registryAddress: %s", registryAddress)
	eth.Registry, err = NewRegistry(common.HexToAddress(registryAddress), eth.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init registry smart contract")
	}
	return
}

func (eth *Ethereum) InitMPE() (err error) {
	var n networks
	networksRaw := contracts.GetNetworks(contracts.MultiPartyEscrow)
	err = json.Unmarshal(networksRaw, &n)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal")
		return
	}
	address := n[config.Blockchain.ChainID].Address
	fmt.Println(address)
	log.Debug().Msgf("MPE address: %s", address)
	eth.MPE, err = NewMultiPartyEscrow(common.HexToAddress(address), eth.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init MPE")
	}
	eth.MPEAddress = common.HexToAddress(address)
	return
}

func (eth Ethereum) GetOrgs() (orgsIDs [][32]byte, err error) {
	orgsIDs, err = eth.Registry.ListOrganizations(nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to GetOrgs")
	}
	return
}

type Org struct {
	Found          bool
	Id             [32]byte
	OrgMetadataURI []byte
	Owner          common.Address
	Members        []common.Address
	ServiceIds     [][32]byte
}

func (eth Ethereum) GetOrg(orgID [32]byte) (org Org, err error) {
	org, err = eth.Registry.GetOrganizationById(nil, orgID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to GetOrg")
		return
	}
	return
}

type Service struct {
	Found       bool
	Id          [32]byte
	MetadataURI []byte
}

func (eth Ethereum) GetService(orgID, serviceID [32]byte) (service Service, err error) {
	service, err = eth.Registry.GetServiceRegistrationById(nil, orgID, serviceID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get service from blockchain")
	}
	return
}
