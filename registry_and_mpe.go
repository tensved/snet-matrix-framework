// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MultiPartyEscrowMetaData contains all meta data concerning the MultiPartyEscrow contract.
var MultiPartyEscrowMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"additionalFunds\",\"type\":\"uint256\"}],\"name\":\"ChannelAddFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"claimAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"plannedAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sendBackAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"keepAmount\",\"type\":\"uint256\"}],\"name\":\"ChannelClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newExpiration\",\"type\":\"uint256\"}],\"name\":\"ChannelExtend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"groupId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"ChannelOpen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"claimAmount\",\"type\":\"uint256\"}],\"name\":\"ChannelSenderClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawFunds\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"channels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"groupId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"nextChannelId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"usedMessages\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"groupId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"openChannel\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"groupId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageNonce\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"openChannelByThirdParty\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"groupId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"depositAndOpenChannel\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"channelIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"actualAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"plannedAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bool[]\",\"name\":\"isSendbacks\",\"type\":\"bool[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"multiChannelClaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"plannedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isSendback\",\"type\":\"bool\"}],\"name\":\"channelClaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newExpiration\",\"type\":\"uint256\"}],\"name\":\"channelExtend\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"channelAddFunds\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newExpiration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"channelExtendAndAddFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"}],\"name\":\"channelClaimTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051611d56380380611d568339818101604052602081101561003357600080fd5b5051600380546001600160a01b0319166001600160a01b03909216919091179055611cf3806100636000396000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c8063aa5f510a116100a2578063da2a5b4f11610071578063da2a5b4f146106e5578063e3b3925014610708578063e5949b5d1461074a578063f4606f00146107ae578063fc0c546a146107b65761010b565b8063aa5f510a14610298578063b6b55f2514610645578063b8da092214610662578063baea65b5146106c85761010b565b80632e1a7d4d116100de5780632e1a7d4d1461020f57806345059a5d1461022c5780635a0284001461024f578063a9059cbb1461026c5761010b565b8063047df8f9146101105780630c19d0ec146101665780631d41f87c1461019157806327e235e3146101d7575b600080fd5b610152600480360360a081101561012657600080fd5b506001600160a01b038135811691602081013590911690604081013590606081013590608001356107da565b604080519115158252519081900360200190f35b61018f6004803603606081101561017c57600080fd5b508035906020810135906040013561088a565b005b61018f600480360360e08110156101a757600080fd5b5080359060208101359060408101359060ff6060820135169060808101359060a08101359060c001351515610945565b6101fd600480360360208110156101ed57600080fd5b50356001600160a01b0316610cbc565b60408051918252519081900360200190f35b6101526004803603602081101561022557600080fd5b5035610cce565b6101526004803603604081101561024257600080fd5b5080359060200135610e48565b6101526004803603602081101561026557600080fd5b5035610f50565b6101526004803603604081101561028257600080fd5b506001600160a01b038135169060200135610f65565b61018f600480360360e08110156102ae57600080fd5b810190602081018135600160201b8111156102c857600080fd5b8201836020820111156102da57600080fd5b803590602001918460208302840111600160201b831117156102fb57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561034a57600080fd5b82018360208201111561035c57600080fd5b803590602001918460208302840111600160201b8311171561037d57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156103cc57600080fd5b8201836020820111156103de57600080fd5b803590602001918460208302840111600160201b831117156103ff57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561044e57600080fd5b82018360208201111561046057600080fd5b803590602001918460208302840111600160201b8311171561048157600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156104d057600080fd5b8201836020820111156104e257600080fd5b803590602001918460208302840111600160201b8311171561050357600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561055257600080fd5b82018360208201111561056457600080fd5b803590602001918460208302840111600160201b8311171561058557600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156105d457600080fd5b8201836020820111156105e657600080fd5b803590602001918460208302840111600160201b8311171561060757600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929550611063945050505050565b6101526004803603602081101561065b57600080fd5b50356111a8565b610152600480360361014081101561067957600080fd5b506001600160a01b03813581169160208101358216916040820135169060608101359060808101359060a08101359060c08101359060ff60e082013516906101008101359061012001356112d9565b61018f600480360360208110156106de57600080fd5b5035611595565b610152600480360360408110156106fb57600080fd5b50803590602001356116be565b610152600480360360a081101561071e57600080fd5b506001600160a01b038135811691602081013590911690604081013590606081013590608001356117ae565b6107676004803603602081101561076057600080fd5b5035611867565b604080519788526001600160a01b03968716602089015294861687860152929094166060860152608085015260a084019290925260c0830191909152519081900360e00190f35b6101fd6118b2565b6107be6118b8565b604080516001600160a01b039092168252519081900360200190f35b60006107e5836111a8565b6108205760405162461bcd60e51b8152600401808060200182810382526028815260200180611c4d6028913960400191505060405180910390fd5b61082d86868686866117ae565b61087e576040805162461bcd60e51b815260206004820152601760248201527f556e61626c6520746f206f70656e206368616e6e656c2e000000000000000000604482015290519081900360640190fd5b50600195945050505050565b6108948383610e48565b6108e5576040805162461bcd60e51b815260206004820152601d60248201527f556e61626c6520746f20657874656e6420746865206368616e6e656c2e000000604482015290519081900360640190fd5b6108ef83826116be565b610940576040805162461bcd60e51b815260206004820152601f60248201527f556e61626c6520746f206164642066756e647320746f206368616e6e656c2e00604482015290519081900360640190fd5b505050565b600087815260208190526040902060058101548711156109ac576040805162461bcd60e51b815260206004820152601b60248201527f496e73756666696369656e74206368616e6e656c20616d6f756e740000000000604482015290519081900360640190fd5b60038101546001600160a01b03163314610a01576040805162461bcd60e51b8152602060048201526011602482015270125b9d985b1a59081c9958da5c1a595b9d607a1b604482015290519081900360640190fd5b85871115610a4e576040805162461bcd60e51b8152602060048201526015602482015274125b9d985b1a59081858dd1d585b08185b5bdd5b9d605a1b604482015290519081900360640190fd5b805460408051725f5f4d50455f636c61696d5f6d65737361676560681b6020808301919091523060601b6033830152604782018c9052606782019390935260878082018a90528251808303909101815260a79091019091528051910120600090610ab7906118c7565b9050600060018288888860405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610b18573d6000803e3d6000fd5b5050604051601f19015160028501549092506001600160a01b03808416911614905080610b54575060018301546001600160a01b038281169116145b610b99576040805162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b604482015290519081900360640190fd5b6005830154610bae908a63ffffffff61191816565b600584015533600090815260016020526040902054610bd3908a63ffffffff61196116565b336000908152600160205260409020558315610c5057610bf28a6119bb565b8254600584015460408051928352602083018c90528281018b90526060830191909152600060808301525133918c917f77c3504a57863d978ba4c28ea297490f1f4814365f5ed32b35cbf5b695db003c9181900360a00190a3610cb0565b8254600101808455600584015460408051928352602083018c90528281018b90526000606084015260808301919091525133918c917f77c3504a57863d978ba4c28ea297490f1f4814365f5ed32b35cbf5b695db003c9181900360a00190a35b50505050505050505050565b60016020526000908152604090205481565b33600090815260016020526040812054821115610d1c5760405162461bcd60e51b8152600401808060200182810382526025815260200180611c756025913960400191505060405180910390fd5b6003546040805163a9059cbb60e01b81523360048201526024810185905290516001600160a01b039092169163a9059cbb916044808201926020929091908290030181600087803b158015610d7057600080fd5b505af1158015610d84573d6000803e3d6000fd5b505050506040513d6020811015610d9a57600080fd5b5051610dd75760405162461bcd60e51b8152600401808060200182810382526029815260200180611c246029913960400191505060405180910390fd5b33600090815260016020526040902054610df7908363ffffffff61191816565b33600081815260016020908152604091829020939093558051858152905191927f21901fa892c430ea8bd38b9390225ac8e67eac75ee10ffba16feefc539a288f992918290030190a2506001919050565b600082815260208190526040812060018101546001600160a01b03163314610eaf576040805162461bcd60e51b815260206004820152601560248201527414d95b99195c881b9bdd08185d5d1a1bdc9a5e9959605a1b604482015290519081900360640190fd5b8060060154831015610efe576040805162461bcd60e51b815260206004820152601360248201527224b73b30b634b21032bc3834b930ba34b7b71760691b604482015290519081900360640190fd5b600084815260208181526040918290206006018590558151858152915186927ff8d4e64f6b2b3db6aaf38b319e259285a48ecd0c5bc0115c9928aba297c7342092908290030190a25060019392505050565b60046020526000908152604090205460ff1681565b33600090815260016020526040812054821115610fb35760405162461bcd60e51b8152600401808060200182810382526024815260200180611c9a6024913960400191505060405180910390fd5b33600090815260016020526040902054610fd3908363ffffffff61191816565b33600090815260016020526040808220929092556001600160a01b03851681522054611005908363ffffffff61196116565b6001600160a01b0384166000818152600160209081526040918290209390935580518581529051919233927f5a0155838afb0f859197785e575b9ad1afeb456c6e522b6f632ee8465941315e9281900390910190a350600192915050565b86518551811480156110755750808751145b80156110815750808551145b801561108d5750808451145b80156110995750808351145b80156110a55750808251145b6110f6576040805162461bcd60e51b815260206004820152601c60248201527f496e76616c69642066756e6374696f6e20706172616d65746572732e00000000604482015290519081900360640190fd5b60005b8181101561119d5761119589828151811061111057fe5b602002602001015189838151811061112457fe5b602002602001015189848151811061113857fe5b602002602001015188858151811061114c57fe5b602002602001015188868151811061116057fe5b602002602001015188878151811061117457fe5b60200260200101518c888151811061118857fe5b6020026020010151610945565b6001016110f9565b505050505050505050565b600354604080516323b872dd60e01b81523360048201523060248201526044810184905290516000926001600160a01b0316916323b872dd91606480830192602092919082900301818787803b15801561120157600080fd5b505af1158015611215573d6000803e3d6000fd5b505050506040513d602081101561122b57600080fd5b50516112685760405162461bcd60e51b8152600401808060200182810382526029815260200180611c246029913960400191505060405180910390fd5b33600090815260016020526040902054611288908363ffffffff61196116565b33600081815260016020908152604091829020939093558051858152905191927fd241e73300212f6df233a8e6d3146b88a9d4964e06621d54b5ff6afeba7b1b8892918290030190a2506001919050565b33600090815260016020526040812054871115611334576040805162461bcd60e51b8152602060048201526014602482015273496e73756666696369656e742062616c616e636560601b604482015290519081900360640190fd5b604080517f5f5f6f70656e4368616e6e656c4279546869726450617274790000000000000060208083019190915230606090811b603984015233811b604d8401526bffffffffffffffffffffffff198e821b81166061850152908d901b166075830152608982018b905260a982018a905260c9820189905260e98083018990528351808403909101815261010990920190925280519101206000906113d8906118c7565b60008181526004602052604090205490915060ff161561143f576040805162461bcd60e51b815260206004820152601f60248201527f5369676e61747572652068617320616c7265616479206265656e207573656400604482015290519081900360640190fd5b60016004600083815260200190815260200160002060006101000a81548160ff0219169083151502179055508b6001600160a01b031660018287878760405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156114d2573d6000803e3d6000fd5b505050602060405103516001600160a01b03161461152b576040805162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b604482015290519081900360640190fd5b6115398c8c8c8c8c8c611a2b565b611583576040805162461bcd60e51b8152602060048201526016602482015275155b98589b19481d1bc81bdc195b8818da185b9b995b60521b604482015290519081900360640190fd5b5060019b9a5050505050505050505050565b6000818152602081905260409020600101546001600160a01b031633146115fc576040805162461bcd60e51b815260206004820152601660248201527529b2b73232b9103737ba1030baba3437b934bd32b21760511b604482015290519081900360640190fd5b600081815260208190526040902060060154431015611662576040805162461bcd60e51b815260206004820152601760248201527f436c61696d2063616c6c656420746f6f206561726c792e000000000000000000604482015290519081900360640190fd5b61166b816119bb565b600081815260208181526040918290208054600590910154835191825291810191909152815183927f592ba8545b0ef2ef56ac54c4db27df2bdbb2a60acc1c5a4ac134eccc20cb8096928290030190a250565b3360009081526001602052604081205482111561170c5760405162461bcd60e51b8152600401808060200182810382526024815260200180611c9a6024913960400191505060405180910390fd5b3360009081526001602052604090205461172c908363ffffffff61191816565b33600090815260016020908152604080832093909355858252819052206005015461175d908363ffffffff61196116565b60008481526020818152604091829020600501929092558051848152905185927fb0e2286f86435d8f98d9cf1c908b693792eb905dd03cd40d2b1d23a3e5311a40928290030190a250600192915050565b336000908152600160205260408120548311156117fc5760405162461bcd60e51b8152600401808060200182810382526025815260200180611c756025913960400191505060405180910390fd5b6001600160a01b03861661180f57600080fd5b61181d338787878787611a2b565b61087e576040805162461bcd60e51b8152602060048201526016602482015275155b98589b19481d1bc81bdc195b8818da185b9b995b60521b604482015290519081900360640190fd5b600060208190529081526040902080546001820154600283015460038401546004850154600586015460069096015494956001600160a01b0394851695938516949092169290919087565b60025481565b6003546001600160a01b031681565b604080517f19457468657265756d205369676e6564204d6573736167653a0a333200000000602080830191909152603c8083019490945282518083039094018452605c909101909152815191012090565b600061195a83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f770000815250611b8c565b9392505050565b60008282018381101561195a576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b60008181526020818152604080832060058101546001808301546001600160a01b031686529093529220546119f59163ffffffff61196116565b6001828101546001600160a01b031660009081526020829052604081209290925560058301829055825401825560069091015550565b6040805160e08101825260008082526001600160a01b03808a1660208085019182528a83168587019081528a841660608701908152608087018b815260a088018b815260c089018b8152600280548a528987528b8a209a518b5596516001808c018054928b166001600160a01b03199384161790559551978b018054988a1698821698909817909755925160038a01805491909816961695909517909555935160048701559151600586015591516006909401939093553382529190915290812054611afd908463ffffffff61191816565b336000908152600160209081526040808320939093556002548351908152908101919091526001600160a01b038881168284015260608201869052608082018590529151869288811692908b16917f172899db3034d5e4e68a2873998cc66a59bad4610fa6319a51f31f75e84452b79181900360a00190a4506002805460019081019091559695505050505050565b60008184841115611c1b5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611be0578181015183820152602001611bc8565b50505050905090810190601f168015611c0d5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50505090039056fe556e61626c6520746f207472616e7366657220746f6b656e20746f2074686520636f6e74726163742e556e61626c6520746f206465706f73697420746f6b656e20746f2074686520636f6e74726163742e496e73756666696369656e742062616c616e636520696e2074686520636f6e74726163742e496e73756666696369656e742062616c616e636520696e2074686520636f6e7472616374a2646970667358221220cff22c9e6287b881a51a9f49de0bee58907447bdedbf6765d8901045b802bded64736f6c63430006020033",
}

// MultiPartyEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiPartyEscrowMetaData.ABI instead.
var MultiPartyEscrowABI = MultiPartyEscrowMetaData.ABI

// MultiPartyEscrowBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MultiPartyEscrowMetaData.Bin instead.
var MultiPartyEscrowBin = MultiPartyEscrowMetaData.Bin

// DeployMultiPartyEscrow deploys a new Ethereum contract, binding an instance of MultiPartyEscrow to it.
func DeployMultiPartyEscrow(auth *bind.TransactOpts, backend bind.ContractBackend, _token common.Address) (common.Address, *types.Transaction, *MultiPartyEscrow, error) {
	parsed, err := MultiPartyEscrowMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultiPartyEscrowBin), backend, _token)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiPartyEscrow{MultiPartyEscrowCaller: MultiPartyEscrowCaller{contract: contract}, MultiPartyEscrowTransactor: MultiPartyEscrowTransactor{contract: contract}, MultiPartyEscrowFilterer: MultiPartyEscrowFilterer{contract: contract}}, nil
}

// MultiPartyEscrow is an auto generated Go binding around an Ethereum contract.
type MultiPartyEscrow struct {
	MultiPartyEscrowCaller     // Read-only binding to the contract
	MultiPartyEscrowTransactor // Write-only binding to the contract
	MultiPartyEscrowFilterer   // Log filterer for contract events
}

// MultiPartyEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiPartyEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiPartyEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiPartyEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiPartyEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiPartyEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiPartyEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiPartyEscrowSession struct {
	Contract     *MultiPartyEscrow // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultiPartyEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiPartyEscrowCallerSession struct {
	Contract *MultiPartyEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// MultiPartyEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiPartyEscrowTransactorSession struct {
	Contract     *MultiPartyEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// MultiPartyEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiPartyEscrowRaw struct {
	Contract *MultiPartyEscrow // Generic contract binding to access the raw methods on
}

// MultiPartyEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiPartyEscrowCallerRaw struct {
	Contract *MultiPartyEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// MultiPartyEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiPartyEscrowTransactorRaw struct {
	Contract *MultiPartyEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiPartyEscrow creates a new instance of MultiPartyEscrow, bound to a specific deployed contract.
func NewMultiPartyEscrow(address common.Address, backend bind.ContractBackend) (*MultiPartyEscrow, error) {
	contract, err := bindMultiPartyEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrow{MultiPartyEscrowCaller: MultiPartyEscrowCaller{contract: contract}, MultiPartyEscrowTransactor: MultiPartyEscrowTransactor{contract: contract}, MultiPartyEscrowFilterer: MultiPartyEscrowFilterer{contract: contract}}, nil
}

// NewMultiPartyEscrowCaller creates a new read-only instance of MultiPartyEscrow, bound to a specific deployed contract.
func NewMultiPartyEscrowCaller(address common.Address, caller bind.ContractCaller) (*MultiPartyEscrowCaller, error) {
	contract, err := bindMultiPartyEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowCaller{contract: contract}, nil
}

// NewMultiPartyEscrowTransactor creates a new write-only instance of MultiPartyEscrow, bound to a specific deployed contract.
func NewMultiPartyEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiPartyEscrowTransactor, error) {
	contract, err := bindMultiPartyEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowTransactor{contract: contract}, nil
}

// NewMultiPartyEscrowFilterer creates a new log filterer instance of MultiPartyEscrow, bound to a specific deployed contract.
func NewMultiPartyEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiPartyEscrowFilterer, error) {
	contract, err := bindMultiPartyEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowFilterer{contract: contract}, nil
}

// bindMultiPartyEscrow binds a generic wrapper to an already deployed contract.
func bindMultiPartyEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiPartyEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiPartyEscrow *MultiPartyEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiPartyEscrow.Contract.MultiPartyEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiPartyEscrow *MultiPartyEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.MultiPartyEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiPartyEscrow *MultiPartyEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.MultiPartyEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiPartyEscrow *MultiPartyEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiPartyEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiPartyEscrow *MultiPartyEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiPartyEscrow *MultiPartyEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_MultiPartyEscrow *MultiPartyEscrowCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MultiPartyEscrow.contract.Call(opts, &out, "balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_MultiPartyEscrow *MultiPartyEscrowSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _MultiPartyEscrow.Contract.Balances(&_MultiPartyEscrow.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_MultiPartyEscrow *MultiPartyEscrowCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _MultiPartyEscrow.Contract.Balances(&_MultiPartyEscrow.CallOpts, arg0)
}

// Channels is a free data retrieval call binding the contract method 0xe5949b5d.
//
// Solidity: function channels(uint256 ) view returns(uint256 nonce, address sender, address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration)
func (_MultiPartyEscrow *MultiPartyEscrowCaller) Channels(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Nonce      *big.Int
	Sender     common.Address
	Signer     common.Address
	Recipient  common.Address
	GroupId    [32]byte
	Value      *big.Int
	Expiration *big.Int
}, error) {
	var out []interface{}
	err := _MultiPartyEscrow.contract.Call(opts, &out, "channels", arg0)

	outstruct := new(struct {
		Nonce      *big.Int
		Sender     common.Address
		Signer     common.Address
		Recipient  common.Address
		GroupId    [32]byte
		Value      *big.Int
		Expiration *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Nonce = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Sender = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Signer = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Recipient = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.GroupId = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.Value = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Expiration = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Channels is a free data retrieval call binding the contract method 0xe5949b5d.
//
// Solidity: function channels(uint256 ) view returns(uint256 nonce, address sender, address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration)
func (_MultiPartyEscrow *MultiPartyEscrowSession) Channels(arg0 *big.Int) (struct {
	Nonce      *big.Int
	Sender     common.Address
	Signer     common.Address
	Recipient  common.Address
	GroupId    [32]byte
	Value      *big.Int
	Expiration *big.Int
}, error) {
	return _MultiPartyEscrow.Contract.Channels(&_MultiPartyEscrow.CallOpts, arg0)
}

// Channels is a free data retrieval call binding the contract method 0xe5949b5d.
//
// Solidity: function channels(uint256 ) view returns(uint256 nonce, address sender, address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration)
func (_MultiPartyEscrow *MultiPartyEscrowCallerSession) Channels(arg0 *big.Int) (struct {
	Nonce      *big.Int
	Sender     common.Address
	Signer     common.Address
	Recipient  common.Address
	GroupId    [32]byte
	Value      *big.Int
	Expiration *big.Int
}, error) {
	return _MultiPartyEscrow.Contract.Channels(&_MultiPartyEscrow.CallOpts, arg0)
}

// NextChannelId is a free data retrieval call binding the contract method 0xf4606f00.
//
// Solidity: function nextChannelId() view returns(uint256)
func (_MultiPartyEscrow *MultiPartyEscrowCaller) NextChannelId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiPartyEscrow.contract.Call(opts, &out, "nextChannelId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextChannelId is a free data retrieval call binding the contract method 0xf4606f00.
//
// Solidity: function nextChannelId() view returns(uint256)
func (_MultiPartyEscrow *MultiPartyEscrowSession) NextChannelId() (*big.Int, error) {
	return _MultiPartyEscrow.Contract.NextChannelId(&_MultiPartyEscrow.CallOpts)
}

// NextChannelId is a free data retrieval call binding the contract method 0xf4606f00.
//
// Solidity: function nextChannelId() view returns(uint256)
func (_MultiPartyEscrow *MultiPartyEscrowCallerSession) NextChannelId() (*big.Int, error) {
	return _MultiPartyEscrow.Contract.NextChannelId(&_MultiPartyEscrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_MultiPartyEscrow *MultiPartyEscrowCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiPartyEscrow.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_MultiPartyEscrow *MultiPartyEscrowSession) Token() (common.Address, error) {
	return _MultiPartyEscrow.Contract.Token(&_MultiPartyEscrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_MultiPartyEscrow *MultiPartyEscrowCallerSession) Token() (common.Address, error) {
	return _MultiPartyEscrow.Contract.Token(&_MultiPartyEscrow.CallOpts)
}

// UsedMessages is a free data retrieval call binding the contract method 0x5a028400.
//
// Solidity: function usedMessages(bytes32 ) view returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowCaller) UsedMessages(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _MultiPartyEscrow.contract.Call(opts, &out, "usedMessages", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedMessages is a free data retrieval call binding the contract method 0x5a028400.
//
// Solidity: function usedMessages(bytes32 ) view returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) UsedMessages(arg0 [32]byte) (bool, error) {
	return _MultiPartyEscrow.Contract.UsedMessages(&_MultiPartyEscrow.CallOpts, arg0)
}

// UsedMessages is a free data retrieval call binding the contract method 0x5a028400.
//
// Solidity: function usedMessages(bytes32 ) view returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowCallerSession) UsedMessages(arg0 [32]byte) (bool, error) {
	return _MultiPartyEscrow.Contract.UsedMessages(&_MultiPartyEscrow.CallOpts, arg0)
}

// ChannelAddFunds is a paid mutator transaction binding the contract method 0xda2a5b4f.
//
// Solidity: function channelAddFunds(uint256 channelId, uint256 amount) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) ChannelAddFunds(opts *bind.TransactOpts, channelId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "channelAddFunds", channelId, amount)
}

// ChannelAddFunds is a paid mutator transaction binding the contract method 0xda2a5b4f.
//
// Solidity: function channelAddFunds(uint256 channelId, uint256 amount) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) ChannelAddFunds(channelId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelAddFunds(&_MultiPartyEscrow.TransactOpts, channelId, amount)
}

// ChannelAddFunds is a paid mutator transaction binding the contract method 0xda2a5b4f.
//
// Solidity: function channelAddFunds(uint256 channelId, uint256 amount) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) ChannelAddFunds(channelId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelAddFunds(&_MultiPartyEscrow.TransactOpts, channelId, amount)
}

// ChannelClaim is a paid mutator transaction binding the contract method 0x1d41f87c.
//
// Solidity: function channelClaim(uint256 channelId, uint256 actualAmount, uint256 plannedAmount, uint8 v, bytes32 r, bytes32 s, bool isSendback) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) ChannelClaim(opts *bind.TransactOpts, channelId *big.Int, actualAmount *big.Int, plannedAmount *big.Int, v uint8, r [32]byte, s [32]byte, isSendback bool) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "channelClaim", channelId, actualAmount, plannedAmount, v, r, s, isSendback)
}

// ChannelClaim is a paid mutator transaction binding the contract method 0x1d41f87c.
//
// Solidity: function channelClaim(uint256 channelId, uint256 actualAmount, uint256 plannedAmount, uint8 v, bytes32 r, bytes32 s, bool isSendback) returns()
func (_MultiPartyEscrow *MultiPartyEscrowSession) ChannelClaim(channelId *big.Int, actualAmount *big.Int, plannedAmount *big.Int, v uint8, r [32]byte, s [32]byte, isSendback bool) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelClaim(&_MultiPartyEscrow.TransactOpts, channelId, actualAmount, plannedAmount, v, r, s, isSendback)
}

// ChannelClaim is a paid mutator transaction binding the contract method 0x1d41f87c.
//
// Solidity: function channelClaim(uint256 channelId, uint256 actualAmount, uint256 plannedAmount, uint8 v, bytes32 r, bytes32 s, bool isSendback) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) ChannelClaim(channelId *big.Int, actualAmount *big.Int, plannedAmount *big.Int, v uint8, r [32]byte, s [32]byte, isSendback bool) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelClaim(&_MultiPartyEscrow.TransactOpts, channelId, actualAmount, plannedAmount, v, r, s, isSendback)
}

// ChannelClaimTimeout is a paid mutator transaction binding the contract method 0xbaea65b5.
//
// Solidity: function channelClaimTimeout(uint256 channelId) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) ChannelClaimTimeout(opts *bind.TransactOpts, channelId *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "channelClaimTimeout", channelId)
}

// ChannelClaimTimeout is a paid mutator transaction binding the contract method 0xbaea65b5.
//
// Solidity: function channelClaimTimeout(uint256 channelId) returns()
func (_MultiPartyEscrow *MultiPartyEscrowSession) ChannelClaimTimeout(channelId *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelClaimTimeout(&_MultiPartyEscrow.TransactOpts, channelId)
}

// ChannelClaimTimeout is a paid mutator transaction binding the contract method 0xbaea65b5.
//
// Solidity: function channelClaimTimeout(uint256 channelId) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) ChannelClaimTimeout(channelId *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelClaimTimeout(&_MultiPartyEscrow.TransactOpts, channelId)
}

// ChannelExtend is a paid mutator transaction binding the contract method 0x45059a5d.
//
// Solidity: function channelExtend(uint256 channelId, uint256 newExpiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) ChannelExtend(opts *bind.TransactOpts, channelId *big.Int, newExpiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "channelExtend", channelId, newExpiration)
}

// ChannelExtend is a paid mutator transaction binding the contract method 0x45059a5d.
//
// Solidity: function channelExtend(uint256 channelId, uint256 newExpiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) ChannelExtend(channelId *big.Int, newExpiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelExtend(&_MultiPartyEscrow.TransactOpts, channelId, newExpiration)
}

// ChannelExtend is a paid mutator transaction binding the contract method 0x45059a5d.
//
// Solidity: function channelExtend(uint256 channelId, uint256 newExpiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) ChannelExtend(channelId *big.Int, newExpiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelExtend(&_MultiPartyEscrow.TransactOpts, channelId, newExpiration)
}

// ChannelExtendAndAddFunds is a paid mutator transaction binding the contract method 0x0c19d0ec.
//
// Solidity: function channelExtendAndAddFunds(uint256 channelId, uint256 newExpiration, uint256 amount) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) ChannelExtendAndAddFunds(opts *bind.TransactOpts, channelId *big.Int, newExpiration *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "channelExtendAndAddFunds", channelId, newExpiration, amount)
}

// ChannelExtendAndAddFunds is a paid mutator transaction binding the contract method 0x0c19d0ec.
//
// Solidity: function channelExtendAndAddFunds(uint256 channelId, uint256 newExpiration, uint256 amount) returns()
func (_MultiPartyEscrow *MultiPartyEscrowSession) ChannelExtendAndAddFunds(channelId *big.Int, newExpiration *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelExtendAndAddFunds(&_MultiPartyEscrow.TransactOpts, channelId, newExpiration, amount)
}

// ChannelExtendAndAddFunds is a paid mutator transaction binding the contract method 0x0c19d0ec.
//
// Solidity: function channelExtendAndAddFunds(uint256 channelId, uint256 newExpiration, uint256 amount) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) ChannelExtendAndAddFunds(channelId *big.Int, newExpiration *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.ChannelExtendAndAddFunds(&_MultiPartyEscrow.TransactOpts, channelId, newExpiration, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) Deposit(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "deposit", value)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) Deposit(value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.Deposit(&_MultiPartyEscrow.TransactOpts, value)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) Deposit(value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.Deposit(&_MultiPartyEscrow.TransactOpts, value)
}

// DepositAndOpenChannel is a paid mutator transaction binding the contract method 0x047df8f9.
//
// Solidity: function depositAndOpenChannel(address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) DepositAndOpenChannel(opts *bind.TransactOpts, signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "depositAndOpenChannel", signer, recipient, groupId, value, expiration)
}

// DepositAndOpenChannel is a paid mutator transaction binding the contract method 0x047df8f9.
//
// Solidity: function depositAndOpenChannel(address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) DepositAndOpenChannel(signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.DepositAndOpenChannel(&_MultiPartyEscrow.TransactOpts, signer, recipient, groupId, value, expiration)
}

// DepositAndOpenChannel is a paid mutator transaction binding the contract method 0x047df8f9.
//
// Solidity: function depositAndOpenChannel(address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) DepositAndOpenChannel(signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.DepositAndOpenChannel(&_MultiPartyEscrow.TransactOpts, signer, recipient, groupId, value, expiration)
}

// MultiChannelClaim is a paid mutator transaction binding the contract method 0xaa5f510a.
//
// Solidity: function multiChannelClaim(uint256[] channelIds, uint256[] actualAmounts, uint256[] plannedAmounts, bool[] isSendbacks, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) MultiChannelClaim(opts *bind.TransactOpts, channelIds []*big.Int, actualAmounts []*big.Int, plannedAmounts []*big.Int, isSendbacks []bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "multiChannelClaim", channelIds, actualAmounts, plannedAmounts, isSendbacks, v, r, s)
}

// MultiChannelClaim is a paid mutator transaction binding the contract method 0xaa5f510a.
//
// Solidity: function multiChannelClaim(uint256[] channelIds, uint256[] actualAmounts, uint256[] plannedAmounts, bool[] isSendbacks, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_MultiPartyEscrow *MultiPartyEscrowSession) MultiChannelClaim(channelIds []*big.Int, actualAmounts []*big.Int, plannedAmounts []*big.Int, isSendbacks []bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.MultiChannelClaim(&_MultiPartyEscrow.TransactOpts, channelIds, actualAmounts, plannedAmounts, isSendbacks, v, r, s)
}

// MultiChannelClaim is a paid mutator transaction binding the contract method 0xaa5f510a.
//
// Solidity: function multiChannelClaim(uint256[] channelIds, uint256[] actualAmounts, uint256[] plannedAmounts, bool[] isSendbacks, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) MultiChannelClaim(channelIds []*big.Int, actualAmounts []*big.Int, plannedAmounts []*big.Int, isSendbacks []bool, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.MultiChannelClaim(&_MultiPartyEscrow.TransactOpts, channelIds, actualAmounts, plannedAmounts, isSendbacks, v, r, s)
}

// OpenChannel is a paid mutator transaction binding the contract method 0xe3b39250.
//
// Solidity: function openChannel(address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) OpenChannel(opts *bind.TransactOpts, signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "openChannel", signer, recipient, groupId, value, expiration)
}

// OpenChannel is a paid mutator transaction binding the contract method 0xe3b39250.
//
// Solidity: function openChannel(address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) OpenChannel(signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.OpenChannel(&_MultiPartyEscrow.TransactOpts, signer, recipient, groupId, value, expiration)
}

// OpenChannel is a paid mutator transaction binding the contract method 0xe3b39250.
//
// Solidity: function openChannel(address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) OpenChannel(signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.OpenChannel(&_MultiPartyEscrow.TransactOpts, signer, recipient, groupId, value, expiration)
}

// OpenChannelByThirdParty is a paid mutator transaction binding the contract method 0xb8da0922.
//
// Solidity: function openChannelByThirdParty(address sender, address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration, uint256 messageNonce, uint8 v, bytes32 r, bytes32 s) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) OpenChannelByThirdParty(opts *bind.TransactOpts, sender common.Address, signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int, messageNonce *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "openChannelByThirdParty", sender, signer, recipient, groupId, value, expiration, messageNonce, v, r, s)
}

// OpenChannelByThirdParty is a paid mutator transaction binding the contract method 0xb8da0922.
//
// Solidity: function openChannelByThirdParty(address sender, address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration, uint256 messageNonce, uint8 v, bytes32 r, bytes32 s) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) OpenChannelByThirdParty(sender common.Address, signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int, messageNonce *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.OpenChannelByThirdParty(&_MultiPartyEscrow.TransactOpts, sender, signer, recipient, groupId, value, expiration, messageNonce, v, r, s)
}

// OpenChannelByThirdParty is a paid mutator transaction binding the contract method 0xb8da0922.
//
// Solidity: function openChannelByThirdParty(address sender, address signer, address recipient, bytes32 groupId, uint256 value, uint256 expiration, uint256 messageNonce, uint8 v, bytes32 r, bytes32 s) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) OpenChannelByThirdParty(sender common.Address, signer common.Address, recipient common.Address, groupId [32]byte, value *big.Int, expiration *big.Int, messageNonce *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.OpenChannelByThirdParty(&_MultiPartyEscrow.TransactOpts, sender, signer, recipient, groupId, value, expiration, messageNonce, v, r, s)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address receiver, uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) Transfer(opts *bind.TransactOpts, receiver common.Address, value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "transfer", receiver, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address receiver, uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) Transfer(receiver common.Address, value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.Transfer(&_MultiPartyEscrow.TransactOpts, receiver, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address receiver, uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) Transfer(receiver common.Address, value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.Transfer(&_MultiPartyEscrow.TransactOpts, receiver, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactor) Withdraw(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.contract.Transact(opts, "withdraw", value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowSession) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.Withdraw(&_MultiPartyEscrow.TransactOpts, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns(bool)
func (_MultiPartyEscrow *MultiPartyEscrowTransactorSession) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _MultiPartyEscrow.Contract.Withdraw(&_MultiPartyEscrow.TransactOpts, value)
}

// MultiPartyEscrowChannelAddFundsIterator is returned from FilterChannelAddFunds and is used to iterate over the raw logs and unpacked data for ChannelAddFunds events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelAddFundsIterator struct {
	Event *MultiPartyEscrowChannelAddFunds // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowChannelAddFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowChannelAddFunds)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowChannelAddFunds)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowChannelAddFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowChannelAddFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowChannelAddFunds represents a ChannelAddFunds event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelAddFunds struct {
	ChannelId       *big.Int
	AdditionalFunds *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterChannelAddFunds is a free log retrieval operation binding the contract event 0xb0e2286f86435d8f98d9cf1c908b693792eb905dd03cd40d2b1d23a3e5311a40.
//
// Solidity: event ChannelAddFunds(uint256 indexed channelId, uint256 additionalFunds)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterChannelAddFunds(opts *bind.FilterOpts, channelId []*big.Int) (*MultiPartyEscrowChannelAddFundsIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "ChannelAddFunds", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowChannelAddFundsIterator{contract: _MultiPartyEscrow.contract, event: "ChannelAddFunds", logs: logs, sub: sub}, nil
}

// WatchChannelAddFunds is a free log subscription operation binding the contract event 0xb0e2286f86435d8f98d9cf1c908b693792eb905dd03cd40d2b1d23a3e5311a40.
//
// Solidity: event ChannelAddFunds(uint256 indexed channelId, uint256 additionalFunds)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchChannelAddFunds(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowChannelAddFunds, channelId []*big.Int) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "ChannelAddFunds", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowChannelAddFunds)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelAddFunds", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelAddFunds is a log parse operation binding the contract event 0xb0e2286f86435d8f98d9cf1c908b693792eb905dd03cd40d2b1d23a3e5311a40.
//
// Solidity: event ChannelAddFunds(uint256 indexed channelId, uint256 additionalFunds)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseChannelAddFunds(log types.Log) (*MultiPartyEscrowChannelAddFunds, error) {
	event := new(MultiPartyEscrowChannelAddFunds)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelAddFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiPartyEscrowChannelClaimIterator is returned from FilterChannelClaim and is used to iterate over the raw logs and unpacked data for ChannelClaim events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelClaimIterator struct {
	Event *MultiPartyEscrowChannelClaim // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowChannelClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowChannelClaim)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowChannelClaim)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowChannelClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowChannelClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowChannelClaim represents a ChannelClaim event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelClaim struct {
	ChannelId      *big.Int
	Nonce          *big.Int
	Recipient      common.Address
	ClaimAmount    *big.Int
	PlannedAmount  *big.Int
	SendBackAmount *big.Int
	KeepAmount     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterChannelClaim is a free log retrieval operation binding the contract event 0x77c3504a57863d978ba4c28ea297490f1f4814365f5ed32b35cbf5b695db003c.
//
// Solidity: event ChannelClaim(uint256 indexed channelId, uint256 nonce, address indexed recipient, uint256 claimAmount, uint256 plannedAmount, uint256 sendBackAmount, uint256 keepAmount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterChannelClaim(opts *bind.FilterOpts, channelId []*big.Int, recipient []common.Address) (*MultiPartyEscrowChannelClaimIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "ChannelClaim", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowChannelClaimIterator{contract: _MultiPartyEscrow.contract, event: "ChannelClaim", logs: logs, sub: sub}, nil
}

// WatchChannelClaim is a free log subscription operation binding the contract event 0x77c3504a57863d978ba4c28ea297490f1f4814365f5ed32b35cbf5b695db003c.
//
// Solidity: event ChannelClaim(uint256 indexed channelId, uint256 nonce, address indexed recipient, uint256 claimAmount, uint256 plannedAmount, uint256 sendBackAmount, uint256 keepAmount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchChannelClaim(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowChannelClaim, channelId []*big.Int, recipient []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "ChannelClaim", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowChannelClaim)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelClaim", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelClaim is a log parse operation binding the contract event 0x77c3504a57863d978ba4c28ea297490f1f4814365f5ed32b35cbf5b695db003c.
//
// Solidity: event ChannelClaim(uint256 indexed channelId, uint256 nonce, address indexed recipient, uint256 claimAmount, uint256 plannedAmount, uint256 sendBackAmount, uint256 keepAmount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseChannelClaim(log types.Log) (*MultiPartyEscrowChannelClaim, error) {
	event := new(MultiPartyEscrowChannelClaim)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelClaim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiPartyEscrowChannelExtendIterator is returned from FilterChannelExtend and is used to iterate over the raw logs and unpacked data for ChannelExtend events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelExtendIterator struct {
	Event *MultiPartyEscrowChannelExtend // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowChannelExtendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowChannelExtend)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowChannelExtend)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowChannelExtendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowChannelExtendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowChannelExtend represents a ChannelExtend event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelExtend struct {
	ChannelId     *big.Int
	NewExpiration *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterChannelExtend is a free log retrieval operation binding the contract event 0xf8d4e64f6b2b3db6aaf38b319e259285a48ecd0c5bc0115c9928aba297c73420.
//
// Solidity: event ChannelExtend(uint256 indexed channelId, uint256 newExpiration)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterChannelExtend(opts *bind.FilterOpts, channelId []*big.Int) (*MultiPartyEscrowChannelExtendIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "ChannelExtend", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowChannelExtendIterator{contract: _MultiPartyEscrow.contract, event: "ChannelExtend", logs: logs, sub: sub}, nil
}

// WatchChannelExtend is a free log subscription operation binding the contract event 0xf8d4e64f6b2b3db6aaf38b319e259285a48ecd0c5bc0115c9928aba297c73420.
//
// Solidity: event ChannelExtend(uint256 indexed channelId, uint256 newExpiration)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchChannelExtend(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowChannelExtend, channelId []*big.Int) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "ChannelExtend", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowChannelExtend)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelExtend", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelExtend is a log parse operation binding the contract event 0xf8d4e64f6b2b3db6aaf38b319e259285a48ecd0c5bc0115c9928aba297c73420.
//
// Solidity: event ChannelExtend(uint256 indexed channelId, uint256 newExpiration)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseChannelExtend(log types.Log) (*MultiPartyEscrowChannelExtend, error) {
	event := new(MultiPartyEscrowChannelExtend)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelExtend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiPartyEscrowChannelOpenIterator is returned from FilterChannelOpen and is used to iterate over the raw logs and unpacked data for ChannelOpen events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelOpenIterator struct {
	Event *MultiPartyEscrowChannelOpen // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowChannelOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowChannelOpen)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowChannelOpen)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowChannelOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowChannelOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowChannelOpen represents a ChannelOpen event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelOpen struct {
	ChannelId  *big.Int
	Nonce      *big.Int
	Sender     common.Address
	Signer     common.Address
	Recipient  common.Address
	GroupId    [32]byte
	Amount     *big.Int
	Expiration *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterChannelOpen is a free log retrieval operation binding the contract event 0x172899db3034d5e4e68a2873998cc66a59bad4610fa6319a51f31f75e84452b7.
//
// Solidity: event ChannelOpen(uint256 channelId, uint256 nonce, address indexed sender, address signer, address indexed recipient, bytes32 indexed groupId, uint256 amount, uint256 expiration)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterChannelOpen(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, groupId [][32]byte) (*MultiPartyEscrowChannelOpenIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var groupIdRule []interface{}
	for _, groupIdItem := range groupId {
		groupIdRule = append(groupIdRule, groupIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "ChannelOpen", senderRule, recipientRule, groupIdRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowChannelOpenIterator{contract: _MultiPartyEscrow.contract, event: "ChannelOpen", logs: logs, sub: sub}, nil
}

// WatchChannelOpen is a free log subscription operation binding the contract event 0x172899db3034d5e4e68a2873998cc66a59bad4610fa6319a51f31f75e84452b7.
//
// Solidity: event ChannelOpen(uint256 channelId, uint256 nonce, address indexed sender, address signer, address indexed recipient, bytes32 indexed groupId, uint256 amount, uint256 expiration)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchChannelOpen(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowChannelOpen, sender []common.Address, recipient []common.Address, groupId [][32]byte) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var groupIdRule []interface{}
	for _, groupIdItem := range groupId {
		groupIdRule = append(groupIdRule, groupIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "ChannelOpen", senderRule, recipientRule, groupIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowChannelOpen)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelOpen", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelOpen is a log parse operation binding the contract event 0x172899db3034d5e4e68a2873998cc66a59bad4610fa6319a51f31f75e84452b7.
//
// Solidity: event ChannelOpen(uint256 channelId, uint256 nonce, address indexed sender, address signer, address indexed recipient, bytes32 indexed groupId, uint256 amount, uint256 expiration)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseChannelOpen(log types.Log) (*MultiPartyEscrowChannelOpen, error) {
	event := new(MultiPartyEscrowChannelOpen)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelOpen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiPartyEscrowChannelSenderClaimIterator is returned from FilterChannelSenderClaim and is used to iterate over the raw logs and unpacked data for ChannelSenderClaim events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelSenderClaimIterator struct {
	Event *MultiPartyEscrowChannelSenderClaim // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowChannelSenderClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowChannelSenderClaim)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowChannelSenderClaim)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowChannelSenderClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowChannelSenderClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowChannelSenderClaim represents a ChannelSenderClaim event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowChannelSenderClaim struct {
	ChannelId   *big.Int
	Nonce       *big.Int
	ClaimAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChannelSenderClaim is a free log retrieval operation binding the contract event 0x592ba8545b0ef2ef56ac54c4db27df2bdbb2a60acc1c5a4ac134eccc20cb8096.
//
// Solidity: event ChannelSenderClaim(uint256 indexed channelId, uint256 nonce, uint256 claimAmount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterChannelSenderClaim(opts *bind.FilterOpts, channelId []*big.Int) (*MultiPartyEscrowChannelSenderClaimIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "ChannelSenderClaim", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowChannelSenderClaimIterator{contract: _MultiPartyEscrow.contract, event: "ChannelSenderClaim", logs: logs, sub: sub}, nil
}

// WatchChannelSenderClaim is a free log subscription operation binding the contract event 0x592ba8545b0ef2ef56ac54c4db27df2bdbb2a60acc1c5a4ac134eccc20cb8096.
//
// Solidity: event ChannelSenderClaim(uint256 indexed channelId, uint256 nonce, uint256 claimAmount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchChannelSenderClaim(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowChannelSenderClaim, channelId []*big.Int) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "ChannelSenderClaim", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowChannelSenderClaim)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelSenderClaim", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelSenderClaim is a log parse operation binding the contract event 0x592ba8545b0ef2ef56ac54c4db27df2bdbb2a60acc1c5a4ac134eccc20cb8096.
//
// Solidity: event ChannelSenderClaim(uint256 indexed channelId, uint256 nonce, uint256 claimAmount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseChannelSenderClaim(log types.Log) (*MultiPartyEscrowChannelSenderClaim, error) {
	event := new(MultiPartyEscrowChannelSenderClaim)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "ChannelSenderClaim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiPartyEscrowDepositFundsIterator is returned from FilterDepositFunds and is used to iterate over the raw logs and unpacked data for DepositFunds events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowDepositFundsIterator struct {
	Event *MultiPartyEscrowDepositFunds // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowDepositFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowDepositFunds)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowDepositFunds)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowDepositFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowDepositFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowDepositFunds represents a DepositFunds event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowDepositFunds struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDepositFunds is a free log retrieval operation binding the contract event 0xd241e73300212f6df233a8e6d3146b88a9d4964e06621d54b5ff6afeba7b1b88.
//
// Solidity: event DepositFunds(address indexed sender, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterDepositFunds(opts *bind.FilterOpts, sender []common.Address) (*MultiPartyEscrowDepositFundsIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "DepositFunds", senderRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowDepositFundsIterator{contract: _MultiPartyEscrow.contract, event: "DepositFunds", logs: logs, sub: sub}, nil
}

// WatchDepositFunds is a free log subscription operation binding the contract event 0xd241e73300212f6df233a8e6d3146b88a9d4964e06621d54b5ff6afeba7b1b88.
//
// Solidity: event DepositFunds(address indexed sender, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchDepositFunds(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowDepositFunds, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "DepositFunds", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowDepositFunds)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "DepositFunds", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDepositFunds is a log parse operation binding the contract event 0xd241e73300212f6df233a8e6d3146b88a9d4964e06621d54b5ff6afeba7b1b88.
//
// Solidity: event DepositFunds(address indexed sender, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseDepositFunds(log types.Log) (*MultiPartyEscrowDepositFunds, error) {
	event := new(MultiPartyEscrowDepositFunds)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "DepositFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiPartyEscrowTransferFundsIterator is returned from FilterTransferFunds and is used to iterate over the raw logs and unpacked data for TransferFunds events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowTransferFundsIterator struct {
	Event *MultiPartyEscrowTransferFunds // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowTransferFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowTransferFunds)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowTransferFunds)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowTransferFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowTransferFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowTransferFunds represents a TransferFunds event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowTransferFunds struct {
	Sender   common.Address
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferFunds is a free log retrieval operation binding the contract event 0x5a0155838afb0f859197785e575b9ad1afeb456c6e522b6f632ee8465941315e.
//
// Solidity: event TransferFunds(address indexed sender, address indexed receiver, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterTransferFunds(opts *bind.FilterOpts, sender []common.Address, receiver []common.Address) (*MultiPartyEscrowTransferFundsIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "TransferFunds", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowTransferFundsIterator{contract: _MultiPartyEscrow.contract, event: "TransferFunds", logs: logs, sub: sub}, nil
}

// WatchTransferFunds is a free log subscription operation binding the contract event 0x5a0155838afb0f859197785e575b9ad1afeb456c6e522b6f632ee8465941315e.
//
// Solidity: event TransferFunds(address indexed sender, address indexed receiver, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchTransferFunds(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowTransferFunds, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "TransferFunds", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowTransferFunds)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "TransferFunds", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferFunds is a log parse operation binding the contract event 0x5a0155838afb0f859197785e575b9ad1afeb456c6e522b6f632ee8465941315e.
//
// Solidity: event TransferFunds(address indexed sender, address indexed receiver, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseTransferFunds(log types.Log) (*MultiPartyEscrowTransferFunds, error) {
	event := new(MultiPartyEscrowTransferFunds)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "TransferFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultiPartyEscrowWithdrawFundsIterator is returned from FilterWithdrawFunds and is used to iterate over the raw logs and unpacked data for WithdrawFunds events raised by the MultiPartyEscrow contract.
type MultiPartyEscrowWithdrawFundsIterator struct {
	Event *MultiPartyEscrowWithdrawFunds // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MultiPartyEscrowWithdrawFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiPartyEscrowWithdrawFunds)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MultiPartyEscrowWithdrawFunds)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MultiPartyEscrowWithdrawFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiPartyEscrowWithdrawFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiPartyEscrowWithdrawFunds represents a WithdrawFunds event raised by the MultiPartyEscrow contract.
type MultiPartyEscrowWithdrawFunds struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawFunds is a free log retrieval operation binding the contract event 0x21901fa892c430ea8bd38b9390225ac8e67eac75ee10ffba16feefc539a288f9.
//
// Solidity: event WithdrawFunds(address indexed sender, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) FilterWithdrawFunds(opts *bind.FilterOpts, sender []common.Address) (*MultiPartyEscrowWithdrawFundsIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.FilterLogs(opts, "WithdrawFunds", senderRule)
	if err != nil {
		return nil, err
	}
	return &MultiPartyEscrowWithdrawFundsIterator{contract: _MultiPartyEscrow.contract, event: "WithdrawFunds", logs: logs, sub: sub}, nil
}

// WatchWithdrawFunds is a free log subscription operation binding the contract event 0x21901fa892c430ea8bd38b9390225ac8e67eac75ee10ffba16feefc539a288f9.
//
// Solidity: event WithdrawFunds(address indexed sender, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) WatchWithdrawFunds(opts *bind.WatchOpts, sink chan<- *MultiPartyEscrowWithdrawFunds, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MultiPartyEscrow.contract.WatchLogs(opts, "WithdrawFunds", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiPartyEscrowWithdrawFunds)
				if err := _MultiPartyEscrow.contract.UnpackLog(event, "WithdrawFunds", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawFunds is a log parse operation binding the contract event 0x21901fa892c430ea8bd38b9390225ac8e67eac75ee10ffba16feefc539a288f9.
//
// Solidity: event WithdrawFunds(address indexed sender, uint256 amount)
func (_MultiPartyEscrow *MultiPartyEscrowFilterer) ParseWithdrawFunds(log types.Log) (*MultiPartyEscrowWithdrawFunds, error) {
	event := new(MultiPartyEscrowWithdrawFunds)
	if err := _MultiPartyEscrow.contract.UnpackLog(event, "WithdrawFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"}],\"name\":\"OrganizationCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"}],\"name\":\"OrganizationDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"}],\"name\":\"OrganizationModified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"metadataURI\",\"type\":\"bytes\"}],\"name\":\"ServiceCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"}],\"name\":\"ServiceDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"metadataURI\",\"type\":\"bytes\"}],\"name\":\"ServiceMetadataModified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"}],\"name\":\"ServiceTagsModified\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"orgMetadataURI\",\"type\":\"bytes\"},{\"internalType\":\"address[]\",\"name\":\"members\",\"type\":\"address[]\"}],\"name\":\"createOrganization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOrganizationOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"orgMetadataURI\",\"type\":\"bytes\"}],\"name\":\"changeOrganizationMetadataURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"newMembers\",\"type\":\"address[]\"}],\"name\":\"addOrganizationMembers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"existingMembers\",\"type\":\"address[]\"}],\"name\":\"removeOrganizationMembers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"}],\"name\":\"deleteOrganization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"metadataURI\",\"type\":\"bytes\"}],\"name\":\"createServiceRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"metadataURI\",\"type\":\"bytes\"}],\"name\":\"updateServiceRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"}],\"name\":\"deleteServiceRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listOrganizations\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"orgIds\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"}],\"name\":\"getOrganizationById\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"found\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"orgMetadataURI\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"members\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"serviceIds\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"}],\"name\":\"listServicesForOrganization\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"found\",\"type\":\"bool\"},{\"internalType\":\"bytes32[]\",\"name\":\"serviceIds\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orgId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"serviceId\",\"type\":\"bytes32\"}],\"name\":\"getServiceRegistrationById\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"found\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"metadataURI\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]",
	Bin: "0x608060405234801561001057600080fd5b5061002a6301ffc9a760e01b6001600160e01b0361004816565b610043631f91217560e11b6001600160e01b0361004816565b6100cc565b6001600160e01b031980821614156100a7576040805162461bcd60e51b815260206004820152601c60248201527f4552433136353a20696e76616c696420696e7465726661636520696400000000604482015290519081900360640190fd5b6001600160e01b0319166000908152602081905260409020805460ff19166001179055565b611b9a806100db6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063a4123f0f1161008c578063d2e544f711610066578063d2e544f71461059c578063d9219d2814610611578063e443852d14610686578063ef72a9af146106b2576100ea565b8063a4123f0f14610485578063bcb43444146104ff578063ca97758c14610579576100ea565b80636928848d116100c85780636928848d146101fe5780638895303f1461021d578063987b3f4d146102ca5780639a017e0a1461033f576100ea565b806301ffc9a7146100ef578063340d84e51461012a5780633a2cb86014610182575b600080fd5b6101166004803603602081101561010557600080fd5b50356001600160e01b031916610777565b604080519115158252519081900360200190f35b610132610796565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561016e578181015183820152602001610156565b505050509050019250505060405180910390f35b61019f6004803603602081101561019857600080fd5b50356107ef565b604051808315151515815260200180602001828103825283818151815260200191508051906020019060200280838360005b838110156101e95781810151838201526020016101d1565b50505050905001935050505060405180910390f35b61021b6004803603602081101561021457600080fd5b5035610876565b005b6102406004803603604081101561023357600080fd5b5080359060200135610a77565b604051808415151515815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561028d578181015183820152602001610275565b50505050905090810190601f1680156102ba5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b61021b600480360360408110156102e057600080fd5b81359190810190604081016020820135600160201b81111561030157600080fd5b82018360208201111561031357600080fd5b803590602001918460018302840111600160201b8311171561033457600080fd5b509092509050610b78565b61035c6004803603602081101561035557600080fd5b5035610bdb565b604051808715151515815260200186815260200180602001856001600160a01b03166001600160a01b031681526020018060200180602001848103845288818151815260200191508051906020019080838360005b838110156103c95781810151838201526020016103b1565b50505050905090810190601f1680156103f65780820380516001836020036101000a031916815260200191505b508481038352865181528651602091820191808901910280838360005b8381101561042b578181015183820152602001610413565b50505050905001848103825285818151815260200191508051906020019060200280838360005b8381101561046a578181015183820152602001610452565b50505050905001995050505050505050505060405180910390f35b61021b6004803603606081101561049b57600080fd5b813591602081013591810190606081016040820135600160201b8111156104c157600080fd5b8201836020820111156104d357600080fd5b803590602001918460018302840111600160201b831117156104f457600080fd5b509092509050610d95565b61021b6004803603606081101561051557600080fd5b813591602081013591810190606081016040820135600160201b81111561053b57600080fd5b82018360208201111561054d57600080fd5b803590602001918460018302840111600160201b8311171561056e57600080fd5b509092509050610ed0565b61021b6004803603604081101561058f57600080fd5b5080359060200135610f85565b61021b600480360360408110156105b257600080fd5b81359190810190604081016020820135600160201b8111156105d357600080fd5b8201836020820111156105e557600080fd5b803590602001918460208302840111600160201b8311171561060657600080fd5b509092509050610fe2565b61021b6004803603604081101561062757600080fd5b81359190810190604081016020820135600160201b81111561064857600080fd5b82018360208201111561065a57600080fd5b803590602001918460208302840111600160201b8311171561067b57600080fd5b509092509050611030565b61021b6004803603604081101561069c57600080fd5b50803590602001356001600160a01b03166110b3565b61021b600480360360608110156106c857600080fd5b81359190810190604081016020820135600160201b8111156106e957600080fd5b8201836020820111156106fb57600080fd5b803590602001918460018302840111600160201b8311171561071c57600080fd5b919390929091602081019035600160201b81111561073957600080fd5b82018360208201111561074b57600080fd5b803590602001918460208302840111600160201b8311171561076c57600080fd5b509092509050611120565b6001600160e01b03191660009081526020819052604090205460ff1690565b606060018054806020026020016040519081016040528092919081815260200182805480156107e457602002820191906000526020600020905b8154815260200190600101908083116107d0575b505050505090505b90565b60008181526002602052604081205460609061080e5760009150610871565b600083815260026020908152604091829020600501805483518184028101840190945280845260019550909183018282801561086957602002820191906000526020600020905b815481526020019060010190808311610855575b505050505090505b915091565b6108818160016112b7565b61088c816000611371565b6000818152600260205260409020600501545b80156108e557600082815260026020526040902060050180546108dc91849160001985019081106108cc57fe5b906000526020600020015461140e565b6000190161089f565b506000818152600260205260409020600301545b80156109495760008281526002602052604090206003018054610940918491600019850190811061092657fe5b6000918252602090912001546001600160a01b0316611546565b600019016108f9565b5060008181526002602052604081206007015460018054919291600019810190811061097157fe5b90600052602060002001549050806001838154811061098c57fe5b9060005260206000200154146109ce5780600183815481106109aa57fe5b60009182526020808320909101929092558281526002909152604090206007018290555b60018054806109d957fe5b600082815260208082208301600019908101839055909201909255848252600290526040812081815590610a1060018301826118d4565b6002820180546001600160a01b0319169055610a3060038301600061191b565b610a3e60058301600061191b565b50600060079190910181905560405184917fb1dbc279d80967cb8113073651e5919753c82fad7b002581eb7a020eaf89e9a791a2505050565b6000828152600260205260408120548190606090610a985760009250610b71565b6000858152600260209081526040808320878452600601909152902054610ac25760009250610b71565b6000858152600260208181526040808420888552600601825292839020805460019182018054865181851615610100026000190190911695909504601f81018590048502860185019096528585529197509550919290830182828015610b695780601f10610b3e57610100808354040283529160200191610b69565b820191906000526020600020905b815481529060010190602001808311610b4c57829003601f168201915b505050505090505b9250925092565b610b838360016112b7565b610b8e836000611371565b6000838152600260205260409020610baa906001018383611939565b5060405183907f06ccb920be65231f5c9d04dd4883d3c7648ebe5f5317cc7177ee4f4a7cc2d03890600090a2505050565b6000818152600260205260408120548190606090829082908190610c025760009550610d8c565b600087815260026020818152604092839020805460019182018054865181851615610100026000190190911695909504601f8101859004850286018501909652858552919a509850919290830182828015610c9e5780601f10610c7357610100808354040283529160200191610c9e565b820191906000526020600020905b815481529060010190602001808311610c8157829003601f168201915b50505060008a8152600260208181526040928390209182015460039092018054845181840281018401909552808552969a506001600160a01b03909216985091949093509150830182828015610d1d57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610cff575b5050505050915060026000888152602001908152602001600020600501805480602002602001604051908101604052809291908181526020018280548015610d8457602002820191906000526020600020905b815481526020019060010190808311610d70575b505050505090505b91939550919395565b610da08460016112b7565b610dab846001611371565b610db7848460006116dc565b610dbf6119b7565b838152604080516020601f8501819004810282018101909252838152908490849081908401838280828437600092018290525060208087019586528a82526002815260408083206005810154828a01528b84526006018252909120865181559451805187969550610e3994506001860193509101906119db565b5060409182015160029182015560008781526020918252828120600501805460018101825590825290829020018690558151818152908101849052859187917f3229c7118d95880b1bb8abc6231f47f1a63d7b1e7e22fbd91a8ccffc9fa75df9918791879181908101848480828437600083820152604051601f909101601f19169092018290039550909350505050a35050505050565b610edb8460016112b7565b610ee6846001611371565b610ef2848460016116dc565b60008481526002602090815260408083208684526006019091529020610f1c906001018383611939565b5082847fb7b13a2b2a9c0147b27815cbec2c7e5ed10588c9d5812211619614c379174c5a848460405180806020018281038252848482818152602001925080828437600083820152604051601f909101601f19169092018290039550909350505050a350505050565b610f908260016112b7565b610f9b826001611371565b610fa7828260016116dc565b610fb1828261140e565b604051819083907f3caed2ddcd24bae20b0075a02d974ba3a229dace98271af93ce8c8f3ebd9d27290600090a35050565b610fed8360016112b7565b610ff8836001611371565b60005b81811015610baa576110288484848481811061101357fe5b905060200201356001600160a01b0316611546565b600101610ffb565b61103b8360016112b7565b611046836001611371565b611083838383808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506117bb92505050565b60405183907f06ccb920be65231f5c9d04dd4883d3c7648ebe5f5317cc7177ee4f4a7cc2d03890600090a2505050565b6110be8260016112b7565b6110c9826000611371565b600082815260026020819052604080832090910180546001600160a01b0319166001600160a01b0385161790555183917f06ccb920be65231f5c9d04dd4883d3c7648ebe5f5317cc7177ee4f4a7cc2d03891a25050565b61112b8560006112b7565b611133611a49565b60008681526002602090815260409091208251815581830151805184936111619260018501929101906119db565b5060408201516002820180546001600160a01b0319166001600160a01b03909216919091179055606082015180516111a3916003840191602090910190611a8b565b50608082015180516111bf916005840191602090910190611aec565b5060a0919091015160079091015560008681526002602052604090208681556111ec906001018686611939565b50600086815260026020818152604080842092830180546001600160a01b0319163317905560018054600790940184905583810181559093527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf690910188905581518482028181018301909352848152611284928992879187918291908501908490808284376000920191909152506117bb92505050565b60405186907f0e7857bfbd020070a2c8d2fe38c788de1e6adc88cb46f60cf6ec7cd385e81db190600090a2505050505050565b801561131757600082815260026020526040902054611312576040805162461bcd60e51b81526020600482015260126024820152711bdc99c8191bd95cc81b9bdd08195e1a5cdd60721b604482015290519081900360640190fd5b61136d565b6000828152600260205260409020541561136d576040805162461bcd60e51b81526020600482015260126024820152716f726720616c72656164792065786973747360701b604482015290519081900360640190fd5b5050565b600082815260026020819052604090912001546001600160a01b03163314806113bd57508080156113bd5750600082815260026020908152604080832033845260040190915290205415155b61136d576040805162461bcd60e51b815260206004820152601760248201527f756e617574686f72697a656420696e766f636174696f6e000000000000000000604482015290519081900360640190fd5b6000828152600260208181526040808420858552600681018352908420830154868552929091526005018054919291600019810190811061144b57fe5b906000526020600020015490508060026000868152602001908152602001600020600501838154811061147a57fe5b9060005260206000200154146114d95760008481526002602052604090206005018054829190849081106114aa57fe5b600091825260208083209091019290925585815260028083526040808320858452600601909352919020018290555b60008481526002602052604090206005018054806114f357fe5b60008281526020808220830160001990810183905590920190925585825260028152604080832086845260060190915281208181559061153660018301826118d4565b6002820160009055505050505050565b60008281526002602090815260408083206001600160a01b03851684526004019091529020541561136d5760008281526002602081815260408084206001600160a01b03861685526004810183529084205486855292909152600301805491929160001981019081106115b557fe5b60009182526020808320909101548683526002909152604090912060030180546001600160a01b039092169250829160001985019081106115f257fe5b6000918252602090912001546001600160a01b0316146116755760008481526002602052604090206003018054829190600019850190811061163057fe5b600091825260208083209190910180546001600160a01b0319166001600160a01b03948516179055868252600281526040808320938516835260049093019052208290555b600084815260026020526040902060030180548061168f57fe5b60008281526020808220830160001990810180546001600160a01b03191690559092019092558582526002815260408083206001600160a01b038716845260040190915281205550505050565b801561174e576000838152600260209081526040808320858452600601909152902054611749576040805162461bcd60e51b81526020600482015260166024820152751cd95c9d9a58d948191bd95cc81b9bdd08195e1a5cdd60521b604482015290519081900360640190fd5b6117b6565b6000838152600260209081526040808320858452600601909152902054156117b6576040805162461bcd60e51b81526020600482015260166024820152757365727669636520616c72656164792065786973747360501b604482015290519081900360640190fd5b505050565b60005b81518110156117b6576002600084815260200190815260200160002060040160008383815181106117eb57fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002054600014156118cc576002600084815260200190815260200160002060030182828151811061183e57fe5b602090810291909101810151825460018101845560009384528284200180546001600160a01b0319166001600160a01b0390921691909117905584825260029052604081206003810154845190926004909201919085908590811061189f57fe5b60200260200101516001600160a01b03166001600160a01b03168152602001908152602001600020819055505b6001016117be565b50805460018160011615610100020316600290046000825580601f106118fa5750611918565b601f0160209004906000526020600020908101906119189190611b26565b50565b50805460008255906000526020600020908101906119189190611b26565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061197a5782800160ff198235161785556119a7565b828001600101855582156119a7579182015b828111156119a757823582559160200191906001019061198c565b506119b3929150611b26565b5090565b60405180606001604052806000801916815260200160608152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611a1c57805160ff19168380011785556119a7565b828001600101855582156119a7579182015b828111156119a7578251825591602001919060010190611a2e565b6040518060c00160405280600080191681526020016060815260200160006001600160a01b031681526020016060815260200160608152602001600081525090565b828054828255906000526020600020908101928215611ae0579160200282015b82811115611ae057825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611aab565b506119b3929150611b40565b8280548282559060005260206000209081019282156119a757916020028201828111156119a7578251825591602001919060010190611a2e565b6107ec91905b808211156119b35760008155600101611b2c565b6107ec91905b808211156119b35780546001600160a01b0319168155600101611b4656fea2646970667358221220eadf149db10ada8694fba72f8777a07dbad3788ae9e2699cd350158aca4db09064736f6c63430006020033",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// RegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RegistryMetaData.Bin instead.
var RegistryBin = RegistryMetaData.Bin

// DeployRegistry deploys a new Ethereum contract, binding an instance of Registry to it.
func DeployRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Registry, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// GetOrganizationById is a free data retrieval call binding the contract method 0x9a017e0a.
//
// Solidity: function getOrganizationById(bytes32 orgId) view returns(bool found, bytes32 id, bytes orgMetadataURI, address owner, address[] members, bytes32[] serviceIds)
func (_Registry *RegistryCaller) GetOrganizationById(opts *bind.CallOpts, orgId [32]byte) (struct {
	Found          bool
	Id             [32]byte
	OrgMetadataURI []byte
	Owner          common.Address
	Members        []common.Address
	ServiceIds     [][32]byte
}, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getOrganizationById", orgId)

	outstruct := new(struct {
		Found          bool
		Id             [32]byte
		OrgMetadataURI []byte
		Owner          common.Address
		Members        []common.Address
		ServiceIds     [][32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Found = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Id = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.OrgMetadataURI = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Members = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)
	outstruct.ServiceIds = *abi.ConvertType(out[5], new([][32]byte)).(*[][32]byte)

	return *outstruct, err

}

// GetOrganizationById is a free data retrieval call binding the contract method 0x9a017e0a.
//
// Solidity: function getOrganizationById(bytes32 orgId) view returns(bool found, bytes32 id, bytes orgMetadataURI, address owner, address[] members, bytes32[] serviceIds)
func (_Registry *RegistrySession) GetOrganizationById(orgId [32]byte) (struct {
	Found          bool
	Id             [32]byte
	OrgMetadataURI []byte
	Owner          common.Address
	Members        []common.Address
	ServiceIds     [][32]byte
}, error) {
	return _Registry.Contract.GetOrganizationById(&_Registry.CallOpts, orgId)
}

// GetOrganizationById is a free data retrieval call binding the contract method 0x9a017e0a.
//
// Solidity: function getOrganizationById(bytes32 orgId) view returns(bool found, bytes32 id, bytes orgMetadataURI, address owner, address[] members, bytes32[] serviceIds)
func (_Registry *RegistryCallerSession) GetOrganizationById(orgId [32]byte) (struct {
	Found          bool
	Id             [32]byte
	OrgMetadataURI []byte
	Owner          common.Address
	Members        []common.Address
	ServiceIds     [][32]byte
}, error) {
	return _Registry.Contract.GetOrganizationById(&_Registry.CallOpts, orgId)
}

// GetServiceRegistrationById is a free data retrieval call binding the contract method 0x8895303f.
//
// Solidity: function getServiceRegistrationById(bytes32 orgId, bytes32 serviceId) view returns(bool found, bytes32 id, bytes metadataURI)
func (_Registry *RegistryCaller) GetServiceRegistrationById(opts *bind.CallOpts, orgId [32]byte, serviceId [32]byte) (struct {
	Found       bool
	Id          [32]byte
	MetadataURI []byte
}, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getServiceRegistrationById", orgId, serviceId)

	outstruct := new(struct {
		Found       bool
		Id          [32]byte
		MetadataURI []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Found = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Id = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.MetadataURI = *abi.ConvertType(out[2], new([]byte)).(*[]byte)

	return *outstruct, err

}

// GetServiceRegistrationById is a free data retrieval call binding the contract method 0x8895303f.
//
// Solidity: function getServiceRegistrationById(bytes32 orgId, bytes32 serviceId) view returns(bool found, bytes32 id, bytes metadataURI)
func (_Registry *RegistrySession) GetServiceRegistrationById(orgId [32]byte, serviceId [32]byte) (struct {
	Found       bool
	Id          [32]byte
	MetadataURI []byte
}, error) {
	return _Registry.Contract.GetServiceRegistrationById(&_Registry.CallOpts, orgId, serviceId)
}

// GetServiceRegistrationById is a free data retrieval call binding the contract method 0x8895303f.
//
// Solidity: function getServiceRegistrationById(bytes32 orgId, bytes32 serviceId) view returns(bool found, bytes32 id, bytes metadataURI)
func (_Registry *RegistryCallerSession) GetServiceRegistrationById(orgId [32]byte, serviceId [32]byte) (struct {
	Found       bool
	Id          [32]byte
	MetadataURI []byte
}, error) {
	return _Registry.Contract.GetServiceRegistrationById(&_Registry.CallOpts, orgId, serviceId)
}

// ListOrganizations is a free data retrieval call binding the contract method 0x340d84e5.
//
// Solidity: function listOrganizations() view returns(bytes32[] orgIds)
func (_Registry *RegistryCaller) ListOrganizations(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "listOrganizations")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// ListOrganizations is a free data retrieval call binding the contract method 0x340d84e5.
//
// Solidity: function listOrganizations() view returns(bytes32[] orgIds)
func (_Registry *RegistrySession) ListOrganizations() ([][32]byte, error) {
	return _Registry.Contract.ListOrganizations(&_Registry.CallOpts)
}

// ListOrganizations is a free data retrieval call binding the contract method 0x340d84e5.
//
// Solidity: function listOrganizations() view returns(bytes32[] orgIds)
func (_Registry *RegistryCallerSession) ListOrganizations() ([][32]byte, error) {
	return _Registry.Contract.ListOrganizations(&_Registry.CallOpts)
}

// ListServicesForOrganization is a free data retrieval call binding the contract method 0x3a2cb860.
//
// Solidity: function listServicesForOrganization(bytes32 orgId) view returns(bool found, bytes32[] serviceIds)
func (_Registry *RegistryCaller) ListServicesForOrganization(opts *bind.CallOpts, orgId [32]byte) (struct {
	Found      bool
	ServiceIds [][32]byte
}, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "listServicesForOrganization", orgId)

	outstruct := new(struct {
		Found      bool
		ServiceIds [][32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Found = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ServiceIds = *abi.ConvertType(out[1], new([][32]byte)).(*[][32]byte)

	return *outstruct, err

}

// ListServicesForOrganization is a free data retrieval call binding the contract method 0x3a2cb860.
//
// Solidity: function listServicesForOrganization(bytes32 orgId) view returns(bool found, bytes32[] serviceIds)
func (_Registry *RegistrySession) ListServicesForOrganization(orgId [32]byte) (struct {
	Found      bool
	ServiceIds [][32]byte
}, error) {
	return _Registry.Contract.ListServicesForOrganization(&_Registry.CallOpts, orgId)
}

// ListServicesForOrganization is a free data retrieval call binding the contract method 0x3a2cb860.
//
// Solidity: function listServicesForOrganization(bytes32 orgId) view returns(bool found, bytes32[] serviceIds)
func (_Registry *RegistryCallerSession) ListServicesForOrganization(orgId [32]byte) (struct {
	Found      bool
	ServiceIds [][32]byte
}, error) {
	return _Registry.Contract.ListServicesForOrganization(&_Registry.CallOpts, orgId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Registry *RegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Registry *RegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Registry.Contract.SupportsInterface(&_Registry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Registry *RegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Registry.Contract.SupportsInterface(&_Registry.CallOpts, interfaceId)
}

// AddOrganizationMembers is a paid mutator transaction binding the contract method 0xd9219d28.
//
// Solidity: function addOrganizationMembers(bytes32 orgId, address[] newMembers) returns()
func (_Registry *RegistryTransactor) AddOrganizationMembers(opts *bind.TransactOpts, orgId [32]byte, newMembers []common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "addOrganizationMembers", orgId, newMembers)
}

// AddOrganizationMembers is a paid mutator transaction binding the contract method 0xd9219d28.
//
// Solidity: function addOrganizationMembers(bytes32 orgId, address[] newMembers) returns()
func (_Registry *RegistrySession) AddOrganizationMembers(orgId [32]byte, newMembers []common.Address) (*types.Transaction, error) {
	return _Registry.Contract.AddOrganizationMembers(&_Registry.TransactOpts, orgId, newMembers)
}

// AddOrganizationMembers is a paid mutator transaction binding the contract method 0xd9219d28.
//
// Solidity: function addOrganizationMembers(bytes32 orgId, address[] newMembers) returns()
func (_Registry *RegistryTransactorSession) AddOrganizationMembers(orgId [32]byte, newMembers []common.Address) (*types.Transaction, error) {
	return _Registry.Contract.AddOrganizationMembers(&_Registry.TransactOpts, orgId, newMembers)
}

// ChangeOrganizationMetadataURI is a paid mutator transaction binding the contract method 0x987b3f4d.
//
// Solidity: function changeOrganizationMetadataURI(bytes32 orgId, bytes orgMetadataURI) returns()
func (_Registry *RegistryTransactor) ChangeOrganizationMetadataURI(opts *bind.TransactOpts, orgId [32]byte, orgMetadataURI []byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "changeOrganizationMetadataURI", orgId, orgMetadataURI)
}

// ChangeOrganizationMetadataURI is a paid mutator transaction binding the contract method 0x987b3f4d.
//
// Solidity: function changeOrganizationMetadataURI(bytes32 orgId, bytes orgMetadataURI) returns()
func (_Registry *RegistrySession) ChangeOrganizationMetadataURI(orgId [32]byte, orgMetadataURI []byte) (*types.Transaction, error) {
	return _Registry.Contract.ChangeOrganizationMetadataURI(&_Registry.TransactOpts, orgId, orgMetadataURI)
}

// ChangeOrganizationMetadataURI is a paid mutator transaction binding the contract method 0x987b3f4d.
//
// Solidity: function changeOrganizationMetadataURI(bytes32 orgId, bytes orgMetadataURI) returns()
func (_Registry *RegistryTransactorSession) ChangeOrganizationMetadataURI(orgId [32]byte, orgMetadataURI []byte) (*types.Transaction, error) {
	return _Registry.Contract.ChangeOrganizationMetadataURI(&_Registry.TransactOpts, orgId, orgMetadataURI)
}

// ChangeOrganizationOwner is a paid mutator transaction binding the contract method 0xe443852d.
//
// Solidity: function changeOrganizationOwner(bytes32 orgId, address newOwner) returns()
func (_Registry *RegistryTransactor) ChangeOrganizationOwner(opts *bind.TransactOpts, orgId [32]byte, newOwner common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "changeOrganizationOwner", orgId, newOwner)
}

// ChangeOrganizationOwner is a paid mutator transaction binding the contract method 0xe443852d.
//
// Solidity: function changeOrganizationOwner(bytes32 orgId, address newOwner) returns()
func (_Registry *RegistrySession) ChangeOrganizationOwner(orgId [32]byte, newOwner common.Address) (*types.Transaction, error) {
	return _Registry.Contract.ChangeOrganizationOwner(&_Registry.TransactOpts, orgId, newOwner)
}

// ChangeOrganizationOwner is a paid mutator transaction binding the contract method 0xe443852d.
//
// Solidity: function changeOrganizationOwner(bytes32 orgId, address newOwner) returns()
func (_Registry *RegistryTransactorSession) ChangeOrganizationOwner(orgId [32]byte, newOwner common.Address) (*types.Transaction, error) {
	return _Registry.Contract.ChangeOrganizationOwner(&_Registry.TransactOpts, orgId, newOwner)
}

// CreateOrganization is a paid mutator transaction binding the contract method 0xef72a9af.
//
// Solidity: function createOrganization(bytes32 orgId, bytes orgMetadataURI, address[] members) returns()
func (_Registry *RegistryTransactor) CreateOrganization(opts *bind.TransactOpts, orgId [32]byte, orgMetadataURI []byte, members []common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "createOrganization", orgId, orgMetadataURI, members)
}

// CreateOrganization is a paid mutator transaction binding the contract method 0xef72a9af.
//
// Solidity: function createOrganization(bytes32 orgId, bytes orgMetadataURI, address[] members) returns()
func (_Registry *RegistrySession) CreateOrganization(orgId [32]byte, orgMetadataURI []byte, members []common.Address) (*types.Transaction, error) {
	return _Registry.Contract.CreateOrganization(&_Registry.TransactOpts, orgId, orgMetadataURI, members)
}

// CreateOrganization is a paid mutator transaction binding the contract method 0xef72a9af.
//
// Solidity: function createOrganization(bytes32 orgId, bytes orgMetadataURI, address[] members) returns()
func (_Registry *RegistryTransactorSession) CreateOrganization(orgId [32]byte, orgMetadataURI []byte, members []common.Address) (*types.Transaction, error) {
	return _Registry.Contract.CreateOrganization(&_Registry.TransactOpts, orgId, orgMetadataURI, members)
}

// CreateServiceRegistration is a paid mutator transaction binding the contract method 0xa4123f0f.
//
// Solidity: function createServiceRegistration(bytes32 orgId, bytes32 serviceId, bytes metadataURI) returns()
func (_Registry *RegistryTransactor) CreateServiceRegistration(opts *bind.TransactOpts, orgId [32]byte, serviceId [32]byte, metadataURI []byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "createServiceRegistration", orgId, serviceId, metadataURI)
}

// CreateServiceRegistration is a paid mutator transaction binding the contract method 0xa4123f0f.
//
// Solidity: function createServiceRegistration(bytes32 orgId, bytes32 serviceId, bytes metadataURI) returns()
func (_Registry *RegistrySession) CreateServiceRegistration(orgId [32]byte, serviceId [32]byte, metadataURI []byte) (*types.Transaction, error) {
	return _Registry.Contract.CreateServiceRegistration(&_Registry.TransactOpts, orgId, serviceId, metadataURI)
}

// CreateServiceRegistration is a paid mutator transaction binding the contract method 0xa4123f0f.
//
// Solidity: function createServiceRegistration(bytes32 orgId, bytes32 serviceId, bytes metadataURI) returns()
func (_Registry *RegistryTransactorSession) CreateServiceRegistration(orgId [32]byte, serviceId [32]byte, metadataURI []byte) (*types.Transaction, error) {
	return _Registry.Contract.CreateServiceRegistration(&_Registry.TransactOpts, orgId, serviceId, metadataURI)
}

// DeleteOrganization is a paid mutator transaction binding the contract method 0x6928848d.
//
// Solidity: function deleteOrganization(bytes32 orgId) returns()
func (_Registry *RegistryTransactor) DeleteOrganization(opts *bind.TransactOpts, orgId [32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "deleteOrganization", orgId)
}

// DeleteOrganization is a paid mutator transaction binding the contract method 0x6928848d.
//
// Solidity: function deleteOrganization(bytes32 orgId) returns()
func (_Registry *RegistrySession) DeleteOrganization(orgId [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.DeleteOrganization(&_Registry.TransactOpts, orgId)
}

// DeleteOrganization is a paid mutator transaction binding the contract method 0x6928848d.
//
// Solidity: function deleteOrganization(bytes32 orgId) returns()
func (_Registry *RegistryTransactorSession) DeleteOrganization(orgId [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.DeleteOrganization(&_Registry.TransactOpts, orgId)
}

// DeleteServiceRegistration is a paid mutator transaction binding the contract method 0xca97758c.
//
// Solidity: function deleteServiceRegistration(bytes32 orgId, bytes32 serviceId) returns()
func (_Registry *RegistryTransactor) DeleteServiceRegistration(opts *bind.TransactOpts, orgId [32]byte, serviceId [32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "deleteServiceRegistration", orgId, serviceId)
}

// DeleteServiceRegistration is a paid mutator transaction binding the contract method 0xca97758c.
//
// Solidity: function deleteServiceRegistration(bytes32 orgId, bytes32 serviceId) returns()
func (_Registry *RegistrySession) DeleteServiceRegistration(orgId [32]byte, serviceId [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.DeleteServiceRegistration(&_Registry.TransactOpts, orgId, serviceId)
}

// DeleteServiceRegistration is a paid mutator transaction binding the contract method 0xca97758c.
//
// Solidity: function deleteServiceRegistration(bytes32 orgId, bytes32 serviceId) returns()
func (_Registry *RegistryTransactorSession) DeleteServiceRegistration(orgId [32]byte, serviceId [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.DeleteServiceRegistration(&_Registry.TransactOpts, orgId, serviceId)
}

// RemoveOrganizationMembers is a paid mutator transaction binding the contract method 0xd2e544f7.
//
// Solidity: function removeOrganizationMembers(bytes32 orgId, address[] existingMembers) returns()
func (_Registry *RegistryTransactor) RemoveOrganizationMembers(opts *bind.TransactOpts, orgId [32]byte, existingMembers []common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "removeOrganizationMembers", orgId, existingMembers)
}

// RemoveOrganizationMembers is a paid mutator transaction binding the contract method 0xd2e544f7.
//
// Solidity: function removeOrganizationMembers(bytes32 orgId, address[] existingMembers) returns()
func (_Registry *RegistrySession) RemoveOrganizationMembers(orgId [32]byte, existingMembers []common.Address) (*types.Transaction, error) {
	return _Registry.Contract.RemoveOrganizationMembers(&_Registry.TransactOpts, orgId, existingMembers)
}

// RemoveOrganizationMembers is a paid mutator transaction binding the contract method 0xd2e544f7.
//
// Solidity: function removeOrganizationMembers(bytes32 orgId, address[] existingMembers) returns()
func (_Registry *RegistryTransactorSession) RemoveOrganizationMembers(orgId [32]byte, existingMembers []common.Address) (*types.Transaction, error) {
	return _Registry.Contract.RemoveOrganizationMembers(&_Registry.TransactOpts, orgId, existingMembers)
}

// UpdateServiceRegistration is a paid mutator transaction binding the contract method 0xbcb43444.
//
// Solidity: function updateServiceRegistration(bytes32 orgId, bytes32 serviceId, bytes metadataURI) returns()
func (_Registry *RegistryTransactor) UpdateServiceRegistration(opts *bind.TransactOpts, orgId [32]byte, serviceId [32]byte, metadataURI []byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "updateServiceRegistration", orgId, serviceId, metadataURI)
}

// UpdateServiceRegistration is a paid mutator transaction binding the contract method 0xbcb43444.
//
// Solidity: function updateServiceRegistration(bytes32 orgId, bytes32 serviceId, bytes metadataURI) returns()
func (_Registry *RegistrySession) UpdateServiceRegistration(orgId [32]byte, serviceId [32]byte, metadataURI []byte) (*types.Transaction, error) {
	return _Registry.Contract.UpdateServiceRegistration(&_Registry.TransactOpts, orgId, serviceId, metadataURI)
}

// UpdateServiceRegistration is a paid mutator transaction binding the contract method 0xbcb43444.
//
// Solidity: function updateServiceRegistration(bytes32 orgId, bytes32 serviceId, bytes metadataURI) returns()
func (_Registry *RegistryTransactorSession) UpdateServiceRegistration(orgId [32]byte, serviceId [32]byte, metadataURI []byte) (*types.Transaction, error) {
	return _Registry.Contract.UpdateServiceRegistration(&_Registry.TransactOpts, orgId, serviceId, metadataURI)
}

// RegistryOrganizationCreatedIterator is returned from FilterOrganizationCreated and is used to iterate over the raw logs and unpacked data for OrganizationCreated events raised by the Registry contract.
type RegistryOrganizationCreatedIterator struct {
	Event *RegistryOrganizationCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryOrganizationCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOrganizationCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryOrganizationCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryOrganizationCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOrganizationCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOrganizationCreated represents a OrganizationCreated event raised by the Registry contract.
type RegistryOrganizationCreated struct {
	OrgId [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOrganizationCreated is a free log retrieval operation binding the contract event 0x0e7857bfbd020070a2c8d2fe38c788de1e6adc88cb46f60cf6ec7cd385e81db1.
//
// Solidity: event OrganizationCreated(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) FilterOrganizationCreated(opts *bind.FilterOpts, orgId [][32]byte) (*RegistryOrganizationCreatedIterator, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OrganizationCreated", orgIdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryOrganizationCreatedIterator{contract: _Registry.contract, event: "OrganizationCreated", logs: logs, sub: sub}, nil
}

// WatchOrganizationCreated is a free log subscription operation binding the contract event 0x0e7857bfbd020070a2c8d2fe38c788de1e6adc88cb46f60cf6ec7cd385e81db1.
//
// Solidity: event OrganizationCreated(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) WatchOrganizationCreated(opts *bind.WatchOpts, sink chan<- *RegistryOrganizationCreated, orgId [][32]byte) (event.Subscription, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OrganizationCreated", orgIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOrganizationCreated)
				if err := _Registry.contract.UnpackLog(event, "OrganizationCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrganizationCreated is a log parse operation binding the contract event 0x0e7857bfbd020070a2c8d2fe38c788de1e6adc88cb46f60cf6ec7cd385e81db1.
//
// Solidity: event OrganizationCreated(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) ParseOrganizationCreated(log types.Log) (*RegistryOrganizationCreated, error) {
	event := new(RegistryOrganizationCreated)
	if err := _Registry.contract.UnpackLog(event, "OrganizationCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryOrganizationDeletedIterator is returned from FilterOrganizationDeleted and is used to iterate over the raw logs and unpacked data for OrganizationDeleted events raised by the Registry contract.
type RegistryOrganizationDeletedIterator struct {
	Event *RegistryOrganizationDeleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryOrganizationDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOrganizationDeleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryOrganizationDeleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryOrganizationDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOrganizationDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOrganizationDeleted represents a OrganizationDeleted event raised by the Registry contract.
type RegistryOrganizationDeleted struct {
	OrgId [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOrganizationDeleted is a free log retrieval operation binding the contract event 0xb1dbc279d80967cb8113073651e5919753c82fad7b002581eb7a020eaf89e9a7.
//
// Solidity: event OrganizationDeleted(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) FilterOrganizationDeleted(opts *bind.FilterOpts, orgId [][32]byte) (*RegistryOrganizationDeletedIterator, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OrganizationDeleted", orgIdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryOrganizationDeletedIterator{contract: _Registry.contract, event: "OrganizationDeleted", logs: logs, sub: sub}, nil
}

// WatchOrganizationDeleted is a free log subscription operation binding the contract event 0xb1dbc279d80967cb8113073651e5919753c82fad7b002581eb7a020eaf89e9a7.
//
// Solidity: event OrganizationDeleted(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) WatchOrganizationDeleted(opts *bind.WatchOpts, sink chan<- *RegistryOrganizationDeleted, orgId [][32]byte) (event.Subscription, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OrganizationDeleted", orgIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOrganizationDeleted)
				if err := _Registry.contract.UnpackLog(event, "OrganizationDeleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrganizationDeleted is a log parse operation binding the contract event 0xb1dbc279d80967cb8113073651e5919753c82fad7b002581eb7a020eaf89e9a7.
//
// Solidity: event OrganizationDeleted(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) ParseOrganizationDeleted(log types.Log) (*RegistryOrganizationDeleted, error) {
	event := new(RegistryOrganizationDeleted)
	if err := _Registry.contract.UnpackLog(event, "OrganizationDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryOrganizationModifiedIterator is returned from FilterOrganizationModified and is used to iterate over the raw logs and unpacked data for OrganizationModified events raised by the Registry contract.
type RegistryOrganizationModifiedIterator struct {
	Event *RegistryOrganizationModified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryOrganizationModifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOrganizationModified)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryOrganizationModified)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryOrganizationModifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOrganizationModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOrganizationModified represents a OrganizationModified event raised by the Registry contract.
type RegistryOrganizationModified struct {
	OrgId [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOrganizationModified is a free log retrieval operation binding the contract event 0x06ccb920be65231f5c9d04dd4883d3c7648ebe5f5317cc7177ee4f4a7cc2d038.
//
// Solidity: event OrganizationModified(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) FilterOrganizationModified(opts *bind.FilterOpts, orgId [][32]byte) (*RegistryOrganizationModifiedIterator, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OrganizationModified", orgIdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryOrganizationModifiedIterator{contract: _Registry.contract, event: "OrganizationModified", logs: logs, sub: sub}, nil
}

// WatchOrganizationModified is a free log subscription operation binding the contract event 0x06ccb920be65231f5c9d04dd4883d3c7648ebe5f5317cc7177ee4f4a7cc2d038.
//
// Solidity: event OrganizationModified(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) WatchOrganizationModified(opts *bind.WatchOpts, sink chan<- *RegistryOrganizationModified, orgId [][32]byte) (event.Subscription, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OrganizationModified", orgIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOrganizationModified)
				if err := _Registry.contract.UnpackLog(event, "OrganizationModified", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrganizationModified is a log parse operation binding the contract event 0x06ccb920be65231f5c9d04dd4883d3c7648ebe5f5317cc7177ee4f4a7cc2d038.
//
// Solidity: event OrganizationModified(bytes32 indexed orgId)
func (_Registry *RegistryFilterer) ParseOrganizationModified(log types.Log) (*RegistryOrganizationModified, error) {
	event := new(RegistryOrganizationModified)
	if err := _Registry.contract.UnpackLog(event, "OrganizationModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryServiceCreatedIterator is returned from FilterServiceCreated and is used to iterate over the raw logs and unpacked data for ServiceCreated events raised by the Registry contract.
type RegistryServiceCreatedIterator struct {
	Event *RegistryServiceCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryServiceCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryServiceCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryServiceCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryServiceCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryServiceCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryServiceCreated represents a ServiceCreated event raised by the Registry contract.
type RegistryServiceCreated struct {
	OrgId       [32]byte
	ServiceId   [32]byte
	MetadataURI []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterServiceCreated is a free log retrieval operation binding the contract event 0x3229c7118d95880b1bb8abc6231f47f1a63d7b1e7e22fbd91a8ccffc9fa75df9.
//
// Solidity: event ServiceCreated(bytes32 indexed orgId, bytes32 indexed serviceId, bytes metadataURI)
func (_Registry *RegistryFilterer) FilterServiceCreated(opts *bind.FilterOpts, orgId [][32]byte, serviceId [][32]byte) (*RegistryServiceCreatedIterator, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "ServiceCreated", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryServiceCreatedIterator{contract: _Registry.contract, event: "ServiceCreated", logs: logs, sub: sub}, nil
}

// WatchServiceCreated is a free log subscription operation binding the contract event 0x3229c7118d95880b1bb8abc6231f47f1a63d7b1e7e22fbd91a8ccffc9fa75df9.
//
// Solidity: event ServiceCreated(bytes32 indexed orgId, bytes32 indexed serviceId, bytes metadataURI)
func (_Registry *RegistryFilterer) WatchServiceCreated(opts *bind.WatchOpts, sink chan<- *RegistryServiceCreated, orgId [][32]byte, serviceId [][32]byte) (event.Subscription, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "ServiceCreated", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryServiceCreated)
				if err := _Registry.contract.UnpackLog(event, "ServiceCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseServiceCreated is a log parse operation binding the contract event 0x3229c7118d95880b1bb8abc6231f47f1a63d7b1e7e22fbd91a8ccffc9fa75df9.
//
// Solidity: event ServiceCreated(bytes32 indexed orgId, bytes32 indexed serviceId, bytes metadataURI)
func (_Registry *RegistryFilterer) ParseServiceCreated(log types.Log) (*RegistryServiceCreated, error) {
	event := new(RegistryServiceCreated)
	if err := _Registry.contract.UnpackLog(event, "ServiceCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryServiceDeletedIterator is returned from FilterServiceDeleted and is used to iterate over the raw logs and unpacked data for ServiceDeleted events raised by the Registry contract.
type RegistryServiceDeletedIterator struct {
	Event *RegistryServiceDeleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryServiceDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryServiceDeleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryServiceDeleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryServiceDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryServiceDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryServiceDeleted represents a ServiceDeleted event raised by the Registry contract.
type RegistryServiceDeleted struct {
	OrgId     [32]byte
	ServiceId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterServiceDeleted is a free log retrieval operation binding the contract event 0x3caed2ddcd24bae20b0075a02d974ba3a229dace98271af93ce8c8f3ebd9d272.
//
// Solidity: event ServiceDeleted(bytes32 indexed orgId, bytes32 indexed serviceId)
func (_Registry *RegistryFilterer) FilterServiceDeleted(opts *bind.FilterOpts, orgId [][32]byte, serviceId [][32]byte) (*RegistryServiceDeletedIterator, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "ServiceDeleted", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryServiceDeletedIterator{contract: _Registry.contract, event: "ServiceDeleted", logs: logs, sub: sub}, nil
}

// WatchServiceDeleted is a free log subscription operation binding the contract event 0x3caed2ddcd24bae20b0075a02d974ba3a229dace98271af93ce8c8f3ebd9d272.
//
// Solidity: event ServiceDeleted(bytes32 indexed orgId, bytes32 indexed serviceId)
func (_Registry *RegistryFilterer) WatchServiceDeleted(opts *bind.WatchOpts, sink chan<- *RegistryServiceDeleted, orgId [][32]byte, serviceId [][32]byte) (event.Subscription, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "ServiceDeleted", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryServiceDeleted)
				if err := _Registry.contract.UnpackLog(event, "ServiceDeleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseServiceDeleted is a log parse operation binding the contract event 0x3caed2ddcd24bae20b0075a02d974ba3a229dace98271af93ce8c8f3ebd9d272.
//
// Solidity: event ServiceDeleted(bytes32 indexed orgId, bytes32 indexed serviceId)
func (_Registry *RegistryFilterer) ParseServiceDeleted(log types.Log) (*RegistryServiceDeleted, error) {
	event := new(RegistryServiceDeleted)
	if err := _Registry.contract.UnpackLog(event, "ServiceDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryServiceMetadataModifiedIterator is returned from FilterServiceMetadataModified and is used to iterate over the raw logs and unpacked data for ServiceMetadataModified events raised by the Registry contract.
type RegistryServiceMetadataModifiedIterator struct {
	Event *RegistryServiceMetadataModified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryServiceMetadataModifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryServiceMetadataModified)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryServiceMetadataModified)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryServiceMetadataModifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryServiceMetadataModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryServiceMetadataModified represents a ServiceMetadataModified event raised by the Registry contract.
type RegistryServiceMetadataModified struct {
	OrgId       [32]byte
	ServiceId   [32]byte
	MetadataURI []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterServiceMetadataModified is a free log retrieval operation binding the contract event 0xb7b13a2b2a9c0147b27815cbec2c7e5ed10588c9d5812211619614c379174c5a.
//
// Solidity: event ServiceMetadataModified(bytes32 indexed orgId, bytes32 indexed serviceId, bytes metadataURI)
func (_Registry *RegistryFilterer) FilterServiceMetadataModified(opts *bind.FilterOpts, orgId [][32]byte, serviceId [][32]byte) (*RegistryServiceMetadataModifiedIterator, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "ServiceMetadataModified", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryServiceMetadataModifiedIterator{contract: _Registry.contract, event: "ServiceMetadataModified", logs: logs, sub: sub}, nil
}

// WatchServiceMetadataModified is a free log subscription operation binding the contract event 0xb7b13a2b2a9c0147b27815cbec2c7e5ed10588c9d5812211619614c379174c5a.
//
// Solidity: event ServiceMetadataModified(bytes32 indexed orgId, bytes32 indexed serviceId, bytes metadataURI)
func (_Registry *RegistryFilterer) WatchServiceMetadataModified(opts *bind.WatchOpts, sink chan<- *RegistryServiceMetadataModified, orgId [][32]byte, serviceId [][32]byte) (event.Subscription, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "ServiceMetadataModified", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryServiceMetadataModified)
				if err := _Registry.contract.UnpackLog(event, "ServiceMetadataModified", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseServiceMetadataModified is a log parse operation binding the contract event 0xb7b13a2b2a9c0147b27815cbec2c7e5ed10588c9d5812211619614c379174c5a.
//
// Solidity: event ServiceMetadataModified(bytes32 indexed orgId, bytes32 indexed serviceId, bytes metadataURI)
func (_Registry *RegistryFilterer) ParseServiceMetadataModified(log types.Log) (*RegistryServiceMetadataModified, error) {
	event := new(RegistryServiceMetadataModified)
	if err := _Registry.contract.UnpackLog(event, "ServiceMetadataModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryServiceTagsModifiedIterator is returned from FilterServiceTagsModified and is used to iterate over the raw logs and unpacked data for ServiceTagsModified events raised by the Registry contract.
type RegistryServiceTagsModifiedIterator struct {
	Event *RegistryServiceTagsModified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryServiceTagsModifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryServiceTagsModified)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryServiceTagsModified)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryServiceTagsModifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryServiceTagsModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryServiceTagsModified represents a ServiceTagsModified event raised by the Registry contract.
type RegistryServiceTagsModified struct {
	OrgId     [32]byte
	ServiceId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterServiceTagsModified is a free log retrieval operation binding the contract event 0xd8b715d9bf49bca034a993b5b934475414e84c904dcbf90dbddd139808b97b05.
//
// Solidity: event ServiceTagsModified(bytes32 indexed orgId, bytes32 indexed serviceId)
func (_Registry *RegistryFilterer) FilterServiceTagsModified(opts *bind.FilterOpts, orgId [][32]byte, serviceId [][32]byte) (*RegistryServiceTagsModifiedIterator, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "ServiceTagsModified", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryServiceTagsModifiedIterator{contract: _Registry.contract, event: "ServiceTagsModified", logs: logs, sub: sub}, nil
}

// WatchServiceTagsModified is a free log subscription operation binding the contract event 0xd8b715d9bf49bca034a993b5b934475414e84c904dcbf90dbddd139808b97b05.
//
// Solidity: event ServiceTagsModified(bytes32 indexed orgId, bytes32 indexed serviceId)
func (_Registry *RegistryFilterer) WatchServiceTagsModified(opts *bind.WatchOpts, sink chan<- *RegistryServiceTagsModified, orgId [][32]byte, serviceId [][32]byte) (event.Subscription, error) {

	var orgIdRule []interface{}
	for _, orgIdItem := range orgId {
		orgIdRule = append(orgIdRule, orgIdItem)
	}
	var serviceIdRule []interface{}
	for _, serviceIdItem := range serviceId {
		serviceIdRule = append(serviceIdRule, serviceIdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "ServiceTagsModified", orgIdRule, serviceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryServiceTagsModified)
				if err := _Registry.contract.UnpackLog(event, "ServiceTagsModified", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseServiceTagsModified is a log parse operation binding the contract event 0xd8b715d9bf49bca034a993b5b934475414e84c904dcbf90dbddd139808b97b05.
//
// Solidity: event ServiceTagsModified(bytes32 indexed orgId, bytes32 indexed serviceId)
func (_Registry *RegistryFilterer) ParseServiceTagsModified(log types.Log) (*RegistryServiceTagsModified, error) {
	event := new(RegistryServiceTagsModified)
	if err := _Registry.contract.UnpackLog(event, "ServiceTagsModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
