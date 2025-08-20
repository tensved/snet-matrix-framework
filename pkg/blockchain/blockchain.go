package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"

	contracts "github.com/singnet/snet-ecosystem-contracts"
	"github.com/tensved/snet-matrix-framework/internal/config"
)

var (
	// HashPrefix32Bytes is an Ethereum signature prefix: see https://github.com/ethereum/go-ethereum/blob/bf468a81ec261745b25206b2a596eb0ee0a24a74/internal/ethapi/api.go#L361
	HashPrefix32Bytes = []byte("\x19Ethereum Signed Message:\n32")
)

const (
	PrefixInSignature = "__MPE_claim_message"
	// Agreed constant value.
	FreeCallPrefixSignature = "__prefix_free_trial"
	// Agreed constant value.
	AllowedUserPrefixSignature = "__authorized_user"
	// PaymentTypeHeader is a type of payment used to pay for an RPC call.
	// Supported types are: "escrow".
	// Note: "job" Payment type is deprecated.
	PaymentTypeHeader = "snet-payment-type"
	// Client that calls the Daemon ( example can be "snet-cli","snet-dapp","snet-sdk").
	ClientTypeHeader = "snet-client-type"
	// Value is a user address , example "0x94d04332C4f5273feF69c4a52D24f42a3aF1F207".
	UserInfoHeader = "snet-user-info"
	// User Agent details set in on the server stream info.
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

	// Added for free call support in Daemon.

	// The user Id of the person making the call.
	FreeCallUserIdHeader = "snet-free-call-user-id"

	// Will be used to check if the Signature is still valid.
	CurrentBlockNumberHeader = "snet-current-block-number"

	// Place holder to set the free call Auth Token issued.
	FreeCallAuthTokenHeader = "snet-free-call-auth-token-bin"
	// Block number on when the Token was issued , to track the expiry of the token , which is ~ 1 Month.
	FreeCallAuthTokenExpiryBlockNumberHeader = "snet-free-call-token-expiry-block"

	// Users may decide to sign upfront and make calls .Daemon generates and Auth Token
	// Users/Clients will need to use this token to make calls for the amount signed upfront.
	PrePaidAuthTokenHeader = "snet-prepaid-auth-token-bin"

	DynamicPriceDerived = "snet-derived-dynamic-price-cost"
)

// EthereumService defines the interface for initializing the Ethereum registry.
type EthereumService interface {
	InitRegistry() (err error)
}

// BindOpts represents binding options for blockchain operations.
type BindOpts struct {
	Call     *bind.CallOpts
	Transact *bind.TransactOpts
	Watch    *bind.WatchOpts
	Filter   *bind.FilterOpts
}

// ChansToWatch represents channels for watching blockchain events.
type ChansToWatch struct {
	ChannelOpens    chan *MultiPartyEscrowChannelOpen
	ChannelExtends  chan *MultiPartyEscrowChannelExtend
	ChannelAddFunds chan *MultiPartyEscrowChannelAddFunds
	DepositFunds    chan *MultiPartyEscrowDepositFunds
	Err             chan error
}

// Ethereum represents the Ethereum client, including HTTP and WebSocket clients, registry, and MPE (MultiPartyEscrow) contracts.
type Ethereum struct {
	Client     *ethclient.Client
	WSSClient  *ethclient.Client
	Registry   *Registry
	MPE        *MultiPartyEscrow
	MPEAddress common.Address
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
}

// Init initializes the Ethereum client and connects to the blockchain via HTTPS and WSS. It also initializes the registry and MPE contracts.
//
// Returns:
//   - e: An instance of Ethereum containing the initialized clients and contracts.
func Init() (e Ethereum) {
	var err error
	log.Debug().Any("ETH_URL", config.Blockchain.EthProviderURL).Msg("ETH_URL")
	e.Client, err = ethclient.Dial(config.Blockchain.EthProviderURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to blockchain via HTTPS")
	}
	e.WSSClient, err = ethclient.Dial(config.Blockchain.EthProviderWSURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to blockchain via WSS")
	}

	e.privateKey, err = crypto.HexToECDSA(config.Blockchain.AdminPrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse private key")
	}

	e.chainID, err = e.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get chain ID")
	}

	err = e.InitRegistry()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init registry")
	}
	err = e.InitMPE()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init MPE")
	}
	return
}

// networks represents a mapping of network names to their respective addresses.
type networks map[string]struct {
	Address string `json:"address"`
}

// InitRegistry initializes the Ethereum registry contract by fetching the network address and creating a new Registry instance.
//
// Returns:
//   - err: An error if the operation fails.
func (eth *Ethereum) InitRegistry() (err error) {
	networksRaw := contracts.GetNetworks(contracts.Registry)
	var n networks
	err = json.Unmarshal(networksRaw, &n)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal")
		return
	}
	registryAddress := n[config.Blockchain.ChainID].Address
	log.Debug().Msgf("registry address: %s", registryAddress)
	eth.Registry, err = NewRegistry(common.HexToAddress(registryAddress), eth.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init registry smart contract")
	}
	return
}

// InitMPE initializes the MultiPartyEscrow contract by fetching the network address and creating a new MultiPartyEscrow instance.
//
// Returns:
//   - err: An error if the operation fails.
func (eth *Ethereum) InitMPE() (err error) {
	var n networks
	networksRaw := contracts.GetNetworks(contracts.MultiPartyEscrow)
	err = json.Unmarshal(networksRaw, &n)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal")
		return
	}
	address := n[config.Blockchain.ChainID].Address
	log.Debug().Msgf("MPE address: %s", address)
	eth.MPE, err = NewMultiPartyEscrow(common.HexToAddress(address), eth.WSSClient)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init MPE")
	}
	eth.MPEAddress = common.HexToAddress(address)
	callOpts := &bind.CallOpts{}
	tokenAddress, _ := eth.MPE.Token(callOpts)
	log.Debug().Msgf("token address: %s", tokenAddress)

	return
}

// GetOrgs retrieves a list of organization IDs from the registry contract.
//
// Returns:
//   - orgsIDs: A slice of organization IDs.
//   - err: An error if the operation fails.
func (eth Ethereum) GetOrgs() (orgsIDs [][32]byte, err error) {
	orgsIDs, err = eth.Registry.ListOrganizations(nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to get organizations from blockchain")
		return [][32]byte{}, err
	}
	return
}

// Org represents an organization with its metadata.
type Org struct {
	Found          bool
	Id             [32]byte
	OrgMetadataURI []byte
	Owner          common.Address
	Members        []common.Address
	ServiceIds     [][32]byte
}

// GetOrg retrieves an organization by its ID from the registry contract.
//
// Parameters:
//   - orgID: The ID of the organization.
//
// Returns:
//   - org: The retrieved organization.
//   - err: An error if the operation fails.
func (eth Ethereum) GetOrg(orgID [32]byte) (org Org, err error) {
	org, err = eth.Registry.GetOrganizationById(nil, orgID)
	if err != nil {
		log.Error().Err(err).Str("org_id", string(orgID[:])).Msg("failed to get organization from blockchain")
		return Org{}, err
	}
	return
}

// Service represents a service with its metadata.
type Service struct {
	Found       bool
	Id          [32]byte
	MetadataURI []byte
}

// GetService retrieves a service by its organization ID and service ID from the registry contract.
//
// Parameters:
//   - orgID: The ID of the organization.
//   - serviceID: The ID of the service.
//
// Returns:
//   - service: The retrieved service.
//   - err: An error if the operation fails.
func (eth Ethereum) GetService(orgID, serviceID [32]byte) (service Service, err error) {
	service, err = eth.Registry.GetServiceRegistrationById(nil, orgID, serviceID)
	if err != nil {
		log.Error().Err(err).Str("org_id", string(orgID[:])).Str("service_id", string(serviceID[:])).Msg("failed to get service from blockchain")
		return Service{}, err
	}
	return
}

// Deprecated: Bytes32ArrayToStrings converts an array of [32]byte to an array of strings,
// trimming all trailing null bytes.
//
// Parameters:
//   - arr: An array of [32]byte to convert.
//
// Returns:
//   - []string: An array of cleaned strings.
func Bytes32ArrayToStrings(arr [][32]byte) []string {
	result := make([]string, len(arr))
	for i, b := range arr {
		clean := bytes.TrimRight(b[:], "\x00")
		result[i] = string(clean)
	}
	return result
}

// Deprecated: StringToBytes32 converts a string to [32]byte, padding with null bytes if necessary.
//
// Parameters:
//   - str: A string to convert.
//
// Returns:
//   - [32]byte: A 32-byte array representation of the string.
func StringToBytes32(str string) [32]byte {
	var byte32 [32]byte
	copy(byte32[:], str)
	return byte32
}

// Deprecated: GetMPEBalance returns the balance of the specified address in the MPE contract.
//
// Parameters:
//   - address: The address to check balance for.
//
// Returns:
//   - *big.Int: The balance amount.
//   - error: An error if the operation fails.
func (eth Ethereum) GetMPEBalance(address common.Address) (*big.Int, error) {
	balance, err := eth.MPE.Balances(nil, address)
	if err != nil {
		log.Error().Err(err).Str("address", address.Hex()).Msg("failed to get MPE balance")
		return nil, err
	}
	return balance, nil
}

// GetTransactOpts creates transaction options for blockchain operations.
func GetTransactOpts(chainID *big.Int, privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	return auth, nil
}

// GetCallOpts creates call options for blockchain operations.
func GetCallOpts(from common.Address, blockNumber *big.Int, ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		From:    from,
		Context: ctx,
	}
}

// GetWatchOpts creates watch options for blockchain operations.
func GetWatchOpts(blockNumber *big.Int, ctx context.Context) *bind.WatchOpts {
	start := blockNumber.Uint64()
	return &bind.WatchOpts{
		Start:   &start,
		Context: ctx,
	}
}

// GetFilterOpts creates filter options for blockchain operations.
func GetFilterOpts(blockNumber *big.Int, ctx context.Context) *bind.FilterOpts {
	return &bind.FilterOpts{
		Start:   blockNumber.Uint64(),
		Context: ctx,
	}
}

// GetNewExpiration calculates new expiration block number.
func GetNewExpiration(currentBlock *big.Int, threshold uint64) *big.Int {
	return new(big.Int).Add(currentBlock, big.NewInt(int64(threshold)))
}

// Deprecated: GetAddressFromPrivateKeyECDSA extracts address from private key.
func GetAddressFromPrivateKeyECDSA(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// DecodePaymentGroupID decodes payment group ID from base64.
func DecodePaymentGroupID(groupID string) ([32]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(groupID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode group ID: %w", err)
	}
	var result [32]byte
	copy(result[:], decoded)
	return result, nil
}

// FilterChannels filters channels by sender, recipient and group
func (eth Ethereum) FilterChannels(senders, recipients []common.Address, groupIDs [][32]byte, filterOpts *bind.FilterOpts) (*MultiPartyEscrowChannelOpen, error) {
	logger := log.With().
		Str("sender", senders[0].Hex()).
		Str("recipient", recipients[0].Hex()).
		Logger()

	logger.Debug().Msg("filtering existing payment channels")

	iter, err := eth.MPE.FilterChannelOpen(filterOpts, senders, recipients, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to filter channels: %w", err)
	}
	defer iter.Close()

	if iter.Next() {
		event := iter.Event
		logger.Info().
			Str("channel_id", event.ChannelId.String()).
			Str("amount", event.Amount.String()).
			Str("expiration", event.Expiration.String()).
			Msg("found existing payment channel")
		return event, nil
	}

	logger.Debug().Msg("no existing payment channels found")
	return nil, nil
}

// EnsureChannelValidity ensures the channel is valid and has sufficient funds.
func (eth Ethereum) EnsureChannelValidity(opened *MultiPartyEscrowChannelOpen, currentSigned, price, newExpiration *big.Int, opts *BindOpts, chans *ChansToWatch) (*big.Int, error) {
	logger := log.With().
		Str("channel_id", opened.ChannelId.String()).
		Str("sender", opened.Sender.Hex()).
		Str("recipient", opened.Recipient.Hex()).
		Logger()

	avail := new(big.Int).Sub(opened.Amount, currentSigned)

	needFunds := avail.Cmp(price) < 0
	needExtend := opened.Expiration.Cmp(newExpiration) <= 0

	logger.Debug().
		Str("current_expiration", opened.Expiration.String()).
		Str("new_expiration", newExpiration.String()).
		Str("available_amount", avail.String()).
		Str("price", price.String()).
		Bool("need_funds", needFunds).
		Bool("need_extend", needExtend).
		Msg("checking channel validity")

	if !needFunds && !needExtend {
		logger.Debug().Msg("channel is valid")
		return opened.ChannelId, nil
	}

	if needFunds {
		logger.Info().Msg("adding funds to payment channel")

		missing := new(big.Int).Sub(price, avail)

		mpeBal, err := eth.MPE.Balances(opts.Call, opened.Sender)
		if err != nil {
			return nil, fmt.Errorf("failed to get MPE balance: %w", err)
		}

		if mpeBal.Cmp(missing) < 0 {
			logger.Info().
				Str("mpe_balance", mpeBal.String()).
				Str("missing_amount", missing.String()).
				Msg("depositing to MPE")

			go eth.watchDepositFunds(opts.Watch, chans.DepositFunds, chans.Err, []common.Address{opened.Sender})

			tx, err := eth.MPE.Deposit(estimateGas(opts.Transact), missing)
			if err != nil {
				return nil, fmt.Errorf("failed to deposit to MPE: %w", err)
			}

			logger.Info().
				Str("tx_hash", tx.Hash().Hex()).
				Str("amount", missing.String()).
				Msg("MPE deposit transaction sent")

			deposited := waitingToDepositToMPE(chans.DepositFunds, chans.Err, 30*time.Second)
			if !deposited {

				return nil, fmt.Errorf("deposit to MPE timeout")

			}
		}

		logger.Info().
			Str("amount", missing.String()).
			Msg("adding funds to payment channel")

		go eth.watchChannelAddFunds(opts.Watch, chans.ChannelAddFunds, chans.Err, []*big.Int{opened.ChannelId})

		tx, err := eth.MPE.ChannelAddFunds(estimateGas(opts.Transact), opened.ChannelId, missing)
		if err != nil {
			return nil, fmt.Errorf("failed to add funds to channel: %w", err)
		}

		logger.Info().
			Str("tx_hash", tx.Hash().Hex()).
			Str("amount", missing.String()).
			Msg("channel add funds transaction sent")

		addedFundsID := waitingForChannelFundsToBeAdded(chans.ChannelAddFunds, chans.Err, 30*time.Second)
		if addedFundsID == nil {

			return nil, fmt.Errorf("add funds timeout")

		}

		logger.Info().Msg("funds added to payment channel successfully")
	}

	if needExtend {
		logger.Info().
			Str("new_expiration", newExpiration.String()).
			Msg("extending payment channel")

		go eth.watchChannelExtend(opts.Watch, chans.ChannelExtends, chans.Err, []*big.Int{opened.ChannelId})

		tx, err := eth.MPE.ChannelExtend(estimateGas(opts.Transact), opened.ChannelId, newExpiration)
		if err != nil {
			return nil, fmt.Errorf("failed to extend channel: %w", err)
		}

		logger.Info().
			Str("tx_hash", tx.Hash().Hex()).
			Str("new_expiration", newExpiration.String()).
			Msg("channel extend transaction sent")

		extendedID := waitingForChannelToExtend(chans.ChannelExtends, chans.Err, 30*time.Second)
		if extendedID == nil {

			return nil, fmt.Errorf("channel extend timeout")

		}

		logger.Info().Msg("payment channel extended successfully")
	}

	return opened.ChannelId, nil
}

// OpenNewChannel opens a new payment channel with the specified parameters.
//
// Parameters:
//   - price: The amount to deposit in the channel.
//   - desiredExpiration: The expiration block number.
//   - opts: Transaction options.
//   - chans: Channels for watching events.
//   - senders: Array of sender addresses.
//   - recipients: Array of recipient addresses.
//   - groupIDs: Array of group IDs.
//
// Returns:
//   - *big.Int: The channel ID.
//   - error: An error if the operation fails.
func (eth Ethereum) OpenNewChannel(price, desiredExpiration *big.Int, opts *BindOpts, chans *ChansToWatch, senders, recipients []common.Address, groupIDs [][32]byte) (*big.Int, error) {
	mpeBal, err := eth.MPE.Balances(opts.Call, senders[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get MPE balance: %w", err)
	}

	log.Info().
		Str("from", senders[0].Hex()).
		Str("recipient", recipients[0].Hex()).
		Str("price", price.String()).
		Str("balance", mpeBal.String()).
		Msg("OpenNewChannel: checking balance")

	open := func() (*big.Int, error) {
		go eth.watchChannelOpen(opts.Watch, chans.ChannelOpens, chans.Err, senders, recipients, groupIDs)

		_, err := eth.MPE.OpenChannel(estimateGas(opts.Transact), senders[0], recipients[0], groupIDs[0], price, desiredExpiration)
		if err != nil {
			return nil, fmt.Errorf("failed to open channel: %w", err)
		}

		id := waitingForChannelToOpen(chans.ChannelOpens, chans.Err, 30*time.Second)
		if id == nil {
			return nil, fmt.Errorf("channel open timeout")
		}
		return id, nil
	}

	if mpeBal.Cmp(price) >= 0 {
		log.Info().Msg("OpenNewChannel: sufficient balance, opening channel")
		return open()
	}

	log.Info().Msg("OpenNewChannel: insufficient balance, depositing and opening channel")

	go eth.watchChannelOpen(opts.Watch, chans.ChannelOpens, chans.Err, senders, recipients, groupIDs)
	go eth.watchDepositFunds(opts.Watch, chans.DepositFunds, chans.Err, senders)

	transactOpts := estimateGas(opts.Transact)
	tx, err := eth.MPE.DepositAndOpenChannel(transactOpts, senders[0], recipients[0], groupIDs[0], price, desiredExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to deposit and open channel: %w", err)
	}

	log.Info().Str("txHash", tx.Hash().Hex()).Msg("OpenNewChannel: transaction sent")

	id := waitingForChannelToOpen(chans.ChannelOpens, chans.Err, 30*time.Second)
	if id == nil {
		return nil, fmt.Errorf("channel open timeout")
	}

	deposited := waitingToDepositToMPE(chans.DepositFunds, chans.Err, 30*time.Second)
	if !deposited {
		return nil, fmt.Errorf("deposit timeout")
	}

	log.Info().Msg("OpenNewChannel: channel created and funded successfully")
	return id, nil
}

// estimateGas estimates gas for transaction
func estimateGas(wallet *bind.TransactOpts) *bind.TransactOpts {
	return wallet
}

// waitingForChannelToOpen waits for channel open event
func waitingForChannelToOpen(channelOpens chan *MultiPartyEscrowChannelOpen, errChan chan error, timeout time.Duration) *big.Int {
	log.Debug().Msg("waiting for channel open event")

	select {
	case event := <-channelOpens:
		log.Info().
			Str("channel_id", event.ChannelId.String()).
			Str("sender", event.Sender.Hex()).
			Str("recipient", event.Recipient.Hex()).
			Str("amount", event.Amount.String()).
			Msg("payment channel opened successfully")
		return event.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msg("error waiting for channel open event")
		return nil
	case <-time.After(timeout):
		log.Error().Msg("timeout waiting for channel open event")
		return nil
	}
}

// waitingToDepositToMPE waits for deposit event in MPE
func waitingToDepositToMPE(channelDepositFunds chan *MultiPartyEscrowDepositFunds, errChan chan error, timeout time.Duration) bool {
	log.Debug().Msg("waiting for deposit event")

	select {
	case event := <-channelDepositFunds:
		log.Info().
			Str("sender", event.Sender.Hex()).
			Str("amount", event.Amount.String()).
			Msg("deposit completed successfully")
		return true
	case err := <-errChan:
		log.Error().Err(err).Msg("error waiting for deposit event")
		return false
	case <-time.After(timeout):
		log.Error().Msg("timeout waiting for deposit event")
		return false
	}
}

// waitingForChannelToExtend waits for channel extend event
func waitingForChannelToExtend(channelExtends chan *MultiPartyEscrowChannelExtend, errChan chan error, timeout time.Duration) *big.Int {
	log.Debug().Msg("waiting for channel extend event")

	select {
	case event := <-channelExtends:
		log.Info().
			Str("channel_id", event.ChannelId.String()).
			Str("new_expiration", event.NewExpiration.String()).
			Msg("payment channel extended successfully")
		return event.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msg("error waiting for channel extend event")
		return nil
	case <-time.After(timeout):
		log.Error().Msg("timeout waiting for channel extend event")
		return nil
	}
}

// waitingForChannelFundsToBeAdded waits for add funds event to channel
func waitingForChannelFundsToBeAdded(channelAddFunds chan *MultiPartyEscrowChannelAddFunds, errChan chan error, timeout time.Duration) *big.Int {
	log.Debug().Msg("waiting for channel add funds event")

	select {
	case event := <-channelAddFunds:
		log.Info().
			Str("channel_id", event.ChannelId.String()).
			Str("additional_funds", event.AdditionalFunds.String()).
			Msg("funds added to payment channel successfully")
		return event.ChannelId
	case err := <-errChan:
		log.Error().Err(err).Msg("error waiting for channel add funds event")
		return nil
	case <-time.After(timeout):
		log.Error().Msg("timeout waiting for channel add funds event")
		return nil
	}
}

// watchChannelOpen watches for channel open events
func (eth Ethereum) watchChannelOpen(watchOpts *bind.WatchOpts, channelOpens chan *MultiPartyEscrowChannelOpen, errChan chan error, senders, recipients []common.Address, groupIDs [][32]byte) {
	log.Debug().Msg("watching for channel open events")

	_, err := eth.MPE.WatchChannelOpen(watchOpts, channelOpens, senders, recipients, groupIDs)
	if err != nil {
		errChan <- fmt.Errorf("failed to watch channel open events: %w", err)
		return
	}
}

// watchDepositFunds watches for deposit events
func (eth Ethereum) watchDepositFunds(watchOpts *bind.WatchOpts, depositFundsChan chan *MultiPartyEscrowDepositFunds, errChan chan error, senders []common.Address) {
	log.Debug().Msg("watching for deposit events")

	_, err := eth.MPE.WatchDepositFunds(watchOpts, depositFundsChan, senders)
	if err != nil {
		errChan <- fmt.Errorf("failed to watch deposit events: %w", err)
		return
	}
}

// watchChannelExtend watches for channel extend events
func (eth Ethereum) watchChannelExtend(watchOpts *bind.WatchOpts, channelExtendsChan chan *MultiPartyEscrowChannelExtend, errChan chan error, channelIDs []*big.Int) {
	log.Debug().Msg("watching for channel extend events")

	_, err := eth.MPE.WatchChannelExtend(watchOpts, channelExtendsChan, channelIDs)
	if err != nil {
		errChan <- fmt.Errorf("failed to watch channel extend events: %w", err)
		return
	}
}

// watchChannelAddFunds watches for add funds events to channel
func (eth Ethereum) watchChannelAddFunds(watchOpts *bind.WatchOpts, channelAddFundsChan chan *MultiPartyEscrowChannelAddFunds, errChan chan error, channelIDs []*big.Int) {
	log.Debug().Msg("watching for channel add funds events")

	_, err := eth.MPE.WatchChannelAddFunds(watchOpts, channelAddFundsChan, channelIDs)
	if err != nil {
		errChan <- fmt.Errorf("failed to watch channel add funds events: %w", err)
		return
	}
}
