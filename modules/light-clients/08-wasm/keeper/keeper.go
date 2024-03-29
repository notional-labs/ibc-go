package keeper

import (
	"bytes"
	"context"
	"cosmossdk.io/collections"
	"crypto/sha256"
	"encoding/hex"
	"github.com/cosmos/ibc-go/modules/light-clients/08-wasm/internal/ibcwasm"
	clientkeeper "github.com/cosmos/ibc-go/v8/modules/core/02-client/keeper"
	"math"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	storetypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	cosmwasm "github.com/CosmWasm/wasmvm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkquery "github.com/cosmos/cosmos-sdk/types/query"

	"github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
)

type Keeper struct {
	storeKey     storetypes.KVStoreService
	cdc          codec.BinaryCodec
	wasmVM       ibcwasm.WasmEngine
	authority    string
	clientKeeper *clientkeeper.Keeper
}

func NewKeeper(cdc codec.BinaryCodec,
	key storetypes.KVStoreService, authority string,
	homeDir string, clientKeeper *clientkeeper.Keeper, queryRouter ibcwasm.QueryRouter,
) Keeper {
	// Wasm VM
	wasmDataDir := filepath.Join(homeDir, "wasm_client_data")
	wasmSupportedFeatures := strings.Join([]string{"storage", "iterator"}, ",")
	wasmMemoryLimitMb := uint32(math.Pow(2, 12))
	wasmPrintDebug := true
	wasmCacheSizeMb := uint32(math.Pow(2, 8))

	vm, err := cosmwasm.NewVM(wasmDataDir, wasmSupportedFeatures, wasmMemoryLimitMb, wasmPrintDebug, wasmCacheSizeMb)
	if err != nil {
		panic(err)
	}
	types.WasmVM = vm
	ibcwasm.SetQueryPlugins(types.NewDefaultQueryPlugins())
	ibcwasm.SetQueryRouter(queryRouter)
	ibcwasm.SetupWasmStoreService(key)

	// governance authority

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		wasmVM:       vm,
		authority:    authority,
		clientKeeper: clientKeeper,
	}
}

func NewKeeperWithVm(cdc codec.BinaryCodec,
	key storetypes.KVStoreService, authority string,
	homeDir string, clientKeeper *clientkeeper.Keeper, queryRouter ibcwasm.QueryRouter,
	vm ibcwasm.WasmEngine,
) Keeper {
	types.WasmVM = vm
	ibcwasm.SetQueryPlugins(types.NewDefaultQueryPlugins())
	ibcwasm.SetQueryRouter(queryRouter)
	ibcwasm.SetupWasmStoreService(key)

	// governance authority

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		wasmVM:       vm,
		authority:    authority,
		clientKeeper: clientKeeper,
	}
}

func (k Keeper) storeWasmCode(ctx sdk.Context, code []byte, storeFn func(code cosmwasm.WasmCode) (cosmwasm.Checksum, error)) ([]byte, error) {
	var err error
	if IsGzip(code) {
		ctx.GasMeter().ConsumeGas(types.VMGasRegister.UncompressCosts(len(code)), "Uncompress gzip bytecode")
		code, err = Uncompress(code, uint64(types.MaxWasmSize))
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrCreateFailed, err.Error())
		}
	}

	// Check to see if the store has a code with the same code it
	codeHash := generateWasmCodeHash(code)
	codeIDKey := types.CodeID(codeHash)
	if types.HasChecksum(ctx, codeIDKey) {
		return nil, types.ErrWasmCodeExists
	}

	// run the code through the wasm light client validation process
	if err := types.ValidateWasmCode(code); err != nil {
		return nil, errorsmod.Wrap(err, "wasm bytecode validation failed")
	}

	// create the code in the vm
	ctx.GasMeter().ConsumeGas(types.VMGasRegister.CompileCosts(len(code)), "Compiling wasm bytecode")
	codeID, err := storeFn(code)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrWasmInvalidCode, "unable to compile wasm code: %s", err)
	}

	// safety check to assert that code id returned by WasmVM equals to code hash
	if !bytes.Equal(codeID, codeHash) {
		return nil, types.ErrWasmInvalidCodeID
	}

	// pin the code to the vm in-memory cache
	if err := types.WasmVM.Pin(codeID); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to pin contract with checksum (%s) to vm cache", hex.EncodeToString(codeID))
	}

	err = ibcwasm.Checksums.Set(ctx, codeHash)

	return codeID, nil
}

func generateWasmCodeHash(code []byte) []byte {
	hash := sha256.Sum256(code)
	return hash[:]
}

func (k Keeper) getWasmCode(c context.Context, query *types.WasmCodeQuery) (*types.WasmCodeResponse, error) {
	if query == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	codeID, err := hex.DecodeString(query.CodeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid code id")
	}

	// Only return checksums we previously stored, not arbitrary checksums that might be stored via e.g Wasmd.
	if !types.HasChecksum(ctx, codeID) {
		return nil, status.Error(codes.NotFound, errorsmod.Wrap(types.ErrWasmChecksumNotFound, query.CodeId).Error())
	}

	codeKey := types.CodeID(codeID)
	code, err := types.WasmVM.GetCode(codeKey)
	if err != nil {
		return nil, status.Error(
			codes.NotFound,
			errorsmod.Wrap(types.ErrWasmCodeIDNotFound, query.CodeId).Error(),
		)
	}

	return &types.WasmCodeResponse{
		Code: code,
	}, nil
}

func (k Keeper) getAllWasmCodeID(c context.Context, query *types.AllWasmCodeIDQuery) (*types.AllWasmCodeIDResponse, error) {
	checksums, pageRes, err := sdkquery.CollectionPaginate(
		c,
		ibcwasm.Checksums,
		query.Pagination,
		func(key []byte, value collections.NoValue) (string, error) {
			return hex.EncodeToString(key), nil
		})
	if err != nil {
		return nil, err
	}

	return &types.AllWasmCodeIDResponse{
		CodeIds:    checksums,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) error {
	for _, contract := range gs.Contracts {
		_, err := k.storeWasmCode(ctx, contract.ContractCode, types.WasmVM.StoreCodeUnchecked)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	checksums, err := types.GetAllChecksums(ctx)
	if err != nil {
		panic(err)
	}

	// Grab code from wasmVM and add to genesis state.
	var genesisState types.GenesisState
	for _, checksum := range checksums {
		code, err := types.WasmVM.GetCode(checksum)
		if err != nil {
			panic(err)
		}
		genesisState.Contracts = append(genesisState.Contracts, types.GenesisContract{
			CodeHash:     checksum,
			ContractCode: code,
		})
	}

	return genesisState
}

// InitializePinnedCodes updates wasmvm to pin to cache all contracts marked as pinned
func InitializePinnedCodes(ctx sdk.Context) error {
	checksums, err := types.GetAllChecksums(ctx)
	if err != nil {
		return err
	}

	for _, checksum := range checksums {
		if err := types.WasmVM.Pin(checksum); err != nil {
			return err
		}
	}
	return nil
}
