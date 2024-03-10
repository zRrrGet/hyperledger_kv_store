package blockstore_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/zRrrGet/blockstore/blockstore"
	"github.com/zRrrGet/blockstore/blockstore/mocks"
	"github.com/stretchr/testify/require"
)


//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestPutEntry(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	blockstoreSC := blockstore.SmartContract{}
	err := blockstoreSC.PutEntry(transactionContext, "123", "456")
	require.NoError(t, err)

	err = blockstoreSC.PutEntry(transactionContext, "123", "456")
	require.NoError(t, err)
}

func TestReadEntry(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedEntry := &blockstore.Entry{Key: "entry"}
	bytes, err := json.Marshal(expectedEntry)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	entrySC := blockstore.SmartContract{}
	entry, err := entrySC.ReadEntry(transactionContext, "entry")
	require.NoError(t, err)
	require.Equal(t, expectedEntry, entry)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve entry"))
	_, err = entrySC.ReadEntry(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve entry")

	chaincodeStub.GetStateReturns(nil, nil)
	entry, err = entrySC.ReadEntry(transactionContext, "entry1")
	require.EqualError(t, err, "the entry entry1 does not exist")
	require.Nil(t, entry)
}

func TestDeleteEntry(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	entry := &blockstore.Entry{Key: "entry1"}
	bytes, err := json.Marshal(entry)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	entrySC := blockstore.SmartContract{}
	err = entrySC.DeleteEntry(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = entrySC.DeleteEntry(transactionContext, "entry1")
	require.EqualError(t, err, "the entry entry1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve entry"))
	err = entrySC.DeleteEntry(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve entry")
}

func TestReadAllEntries(t *testing.T) {
	entry := &blockstore.Entry{Key: "entry1"}
	bytes, err := json.Marshal(entry)
	require.NoError(t, err)

	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	entrySC := &blockstore.SmartContract{}
	entrys, err := entrySC.ReadAllEntries(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*blockstore.Entry{entry}, entrys)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	entrys, err = entrySC.ReadAllEntries(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, entrys)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all entrys"))
	entrys, err = entrySC.ReadAllEntries(transactionContext)
	require.EqualError(t, err, "failed retrieving all entrys")
	require.Nil(t, entrys)
}
