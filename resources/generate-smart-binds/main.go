package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	contracts "github.com/singnet/snet-ecosystem-contracts"
	"log"
	"os"
)

func main() {
	bindContent, err := bind.Bind(
		[]string{"MultiPartyEscrow", "Registry"},
		[]string{
			string(contracts.GetABIClean(contracts.MultiPartyEscrow)),
			string(contracts.GetABIClean(contracts.Registry))},
		[]string{
			string(contracts.GetBytecodeClean(contracts.MultiPartyEscrow)),
			string(contracts.GetBytecodeClean(contracts.Registry))},
		nil, "blockchain", bind.LangGo, nil, nil)
	if err != nil {
		log.Fatalf("failed to generate binding: %v", err)
	}

	if err = os.WriteFile("registry_and_mpe.go", []byte(bindContent), 0600); err != nil {
		log.Fatalf("failed to write ABI binding: %v", err)
	}
}
