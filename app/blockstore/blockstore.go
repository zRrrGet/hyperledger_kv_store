package blockstore

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Entry struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (s *SmartContract) ReadAllEntries(ctx contractapi.TransactionContextInterface) ([]*Entry, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var entries []*Entry
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var entry Entry
		err = json.Unmarshal(queryResponse.Value, &entry)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (s *SmartContract) ReadEntry(ctx contractapi.TransactionContextInterface, key string) (*Entry, error) {
	entryJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if entryJSON == nil {
		return nil, fmt.Errorf("the entry %s does not exist", key)
	}

	var entry Entry
	err = json.Unmarshal(entryJSON, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *SmartContract) PutEntry(ctx contractapi.TransactionContextInterface, key string, value string) error {
	entry := Entry{
		Key:             key,
		Value:           value,
	}
	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(key, entryJSON)
}

func (s *SmartContract) DeleteEntry(ctx contractapi.TransactionContextInterface, key string) error {
	exists, err := s.EntryExists(ctx, key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the entry %s does not exist", key)
	}

	return ctx.GetStub().DelState(key)
}

func (s *SmartContract) EntryExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	entryJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return entryJSON != nil, nil
}