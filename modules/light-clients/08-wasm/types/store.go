package types

import (
	"bytes"
	"context"
	errorsmod "cosmossdk.io/errors"
	wasmvm "github.com/CosmWasm/wasmvm"
	"github.com/cosmos/ibc-go/modules/light-clients/08-wasm/internal/ibcwasm"
	"io"

	"cosmossdk.io/store/cachekv"
	"cosmossdk.io/store/tracekv"
	storetypes "cosmossdk.io/store/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// WrappedStore combines two KVStores into one while transparently routing the calls based on key prefix
type WrappedStore struct {
	first  storetypes.KVStore
	second storetypes.KVStore

	firstPrefix  []byte
	secondPrefix []byte
}

func NewWrappedStore(first, second storetypes.KVStore, firstPrefix, secondPrefix []byte) WrappedStore {
	return WrappedStore{
		first:        first,
		second:       second,
		firstPrefix:  firstPrefix,
		secondPrefix: secondPrefix,
	}
}

func (ws WrappedStore) Get(key []byte) []byte {
	return ws.getStore(key).Get(ws.trimPrefix(key))
}

func (ws WrappedStore) Has(key []byte) bool {
	return ws.getStore(key).Has(ws.trimPrefix(key))
}

func (ws WrappedStore) Set(key, value []byte) {
	ws.getStore(key).Set(ws.trimPrefix(key), value)
}

func (ws WrappedStore) Delete(key []byte) {
	ws.getStore(key).Delete(ws.trimPrefix(key))
}

func (ws WrappedStore) GetStoreType() storetypes.StoreType {
	return ws.first.GetStoreType()
}

func (ws WrappedStore) Iterator(start, end []byte) storetypes.Iterator {
	return ws.getStore(start).Iterator(ws.trimPrefix(start), ws.trimPrefix(end))
}

func (ws WrappedStore) ReverseIterator(start, end []byte) storetypes.Iterator {
	return ws.getStore(start).ReverseIterator(ws.trimPrefix(start), ws.trimPrefix(end))
}

func (ws WrappedStore) CacheWrap() storetypes.CacheWrap {
	return cachekv.NewStore(ws)
}

func (ws WrappedStore) CacheWrapWithTrace(w io.Writer, tc storetypes.TraceContext) storetypes.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(ws, w, tc))
}

func (ws WrappedStore) trimPrefix(key []byte) []byte {
	if bytes.HasPrefix(key, ws.firstPrefix) {
		key = bytes.TrimPrefix(key, ws.firstPrefix)
	} else {
		key = bytes.TrimPrefix(key, ws.secondPrefix)
	}

	return key
}

func (ws WrappedStore) getStore(key []byte) storetypes.KVStore {
	if bytes.HasPrefix(key, ws.firstPrefix) {
		return ws.first
	}

	return ws.second
}

// setClientState stores the client state
func setClientState(clientStore storetypes.KVStore, cdc codec.BinaryCodec, clientState *ClientState) {
	key := host.ClientStateKey()
	val := clienttypes.MustMarshalClientState(cdc, clientState)
	clientStore.Set(key, val)
}

// setConsensusState stores the consensus state at the given height.
func setConsensusState(clientStore storetypes.KVStore, cdc codec.BinaryCodec, consensusState *ConsensusState, height exported.Height) {
	key := host.ConsensusStateKey(height)
	val := clienttypes.MustMarshalConsensusState(cdc, consensusState)
	clientStore.Set(key, val)
}

// GetConsensusState retrieves the consensus state from the client prefixed
// store. An error is returned if the consensus state does not exist.
func GetConsensusState(store storetypes.KVStore, cdc codec.BinaryCodec, height exported.Height) (*ConsensusState, error) {
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, errorsmod.Wrapf(
			clienttypes.ErrConsensusStateNotFound,
			"consensus state does not exist for height %s", height,
		)
	}

	consensusStateI, err := clienttypes.UnmarshalConsensusState(cdc, bz)
	if err != nil {
		return nil, errorsmod.Wrapf(clienttypes.ErrInvalidConsensus, "unmarshal error: %v", err)
	}

	consensusState, ok := consensusStateI.(*ConsensusState)
	if !ok {
		return nil, errorsmod.Wrapf(
			clienttypes.ErrInvalidConsensus,
			"invalid consensus type %T, expected %T", consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}

var _ wasmvmtypes.KVStore = &StoreAdapter{}

// StoreAdapter adapter to bridge SDK store impl to wasmvm
type StoreAdapter struct {
	parent storetypes.KVStore
}

// NewStoreAdapter constructor
func NewStoreAdapter(s storetypes.KVStore) *StoreAdapter {
	if s == nil {
		panic("store must not be nil")
	}
	return &StoreAdapter{parent: s}
}

func (s StoreAdapter) Get(key []byte) []byte {
	return s.parent.Get(key)
}

func (s StoreAdapter) Set(key, value []byte) {
	s.parent.Set(key, value)
}

func (s StoreAdapter) Delete(key []byte) {
	s.parent.Delete(key)
}

func (s StoreAdapter) Iterator(start, end []byte) wasmvmtypes.Iterator {
	return s.parent.Iterator(start, end)
}

func (s StoreAdapter) ReverseIterator(start, end []byte) wasmvmtypes.Iterator {
	return s.parent.ReverseIterator(start, end)
}

// Checksum is a type alias used for wasm byte code checksums.
type Checksum = wasmvmtypes.Checksum

// CreateChecksum creates a sha256 checksum from the given wasm code, it forwards the
// call to the wasmvm package. The code is checked for the following conditions:
// - code length is zero.
// - code length is less than 4 bytes (magic number length).
// - code does not start with the wasm magic number.
func CreateChecksum(code []byte) (Checksum, error) {
	return wasmvm.CreateChecksum(code)
}

// GetAllChecksums is a helper to get all checksums from the store.
// It returns an empty slice if no checksums are found
func GetAllChecksums(ctx context.Context) ([]Checksum, error) {
	iterator, err := ibcwasm.Checksums.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	keys, err := iterator.Keys()
	if err != nil {
		return nil, err
	}

	checksums := []Checksum{}
	for _, key := range keys {
		checksums = append(checksums, key)
	}

	return checksums, nil
}

// HasChecksum returns true if the given checksum exists in the store and
// false otherwise.
func HasChecksum(ctx context.Context, checksum Checksum) bool {
	found, err := ibcwasm.Checksums.Has(ctx, checksum)
	if err != nil {
		return false
	}

	return found
}
