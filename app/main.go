package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zRrrGet/blockstore/blockstore"
)

func main() {
	bsChaincode, err := contractapi.NewChaincode(&blockstore.SmartContract{})
	if err != nil {
		log.Panicf("Error creating blockstore chaincode: %v", err)
	}

	if err := bsChaincode.Start(); err != nil {
		log.Panicf("Error starting blockstore chaincode: %v", err)
	}
}