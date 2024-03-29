package types

import (
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"encoding/json"
	"errors"
	cosmwasm "github.com/CosmWasm/wasmvm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/modules/light-clients/08-wasm/internal/ibcwasm"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var (
	WasmVM        ibcwasm.WasmEngine
	VMGasRegister = NewDefaultWasmGasRegister()
	wasmvmAPI     = cosmwasm.GoAPI{
		HumanAddress:     humanAddress,
		CanonicalAddress: canonicalAddress,
	}
)

type queryResponse struct {
	Status          exported.Status               `json:"status,omitempty"`
	GenesisMetadata []clienttypes.GenesisMetadata `json:"genesis_metadata,omitempty"`
}

type ClientCreateRequest struct {
	ClientCreateRequest ClientState `json:"client_create_request,omitempty"`
}

type ContractResult interface {
	Validate() bool
	Error() string
}

type contractResult struct {
	IsValid           bool   `json:"is_valid,omitempty"`
	ErrorMsg          string `json:"error_msg,omitempty"`
	Data              []byte `json:"data,omitempty"`
	FoundMisbehaviour bool   `json:"found_misbehaviour"`
}

func (r contractResult) Validate() bool {
	return r.IsValid
}

func (r contractResult) Error() string {
	return r.ErrorMsg
}

// Calls vm.Init with appropriate arguments
func initContract(codeID []byte, ctx sdk.Context, store storetypes.KVStore, payload InstantiateMessage) (*wasmvmtypes.Response, error) {
	sdkGasMeter := ctx.GasMeter()
	multipliedGasMeter := NewMultipliedGasMeter(sdkGasMeter, VMGasRegister)
	gasLimit := VMGasRegister.runtimeGasForContract(ctx)
	chainID := ctx.BlockHeader().ChainID
	height := ctx.BlockHeader().Height
	// safety checks before casting below
	if height < 0 {
		panic("Block height must never be negative")
	}
	nsec := ctx.BlockTime().UnixNano()
	if nsec < 0 {
		panic("Block (unix) time must never be negative ")
	}
	env := wasmvmtypes.Env{
		Block: wasmvmtypes.BlockInfo{
			Height:  uint64(height),
			Time:    uint64(nsec),
			ChainID: chainID,
		},
		Contract: wasmvmtypes.ContractInfo{
			Address: "",
		},
	}

	encodedData, err := json.Marshal(payload)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to marshal payload for wasm contract instantiation")
	}

	msgInfo := wasmvmtypes.MessageInfo{
		Sender: "",
		Funds:  nil,
	}

	ctx.GasMeter().ConsumeGas(VMGasRegister.NewContractInstanceCosts(len(encodedData)), "Loading CosmWasm module: instantiate")
	response, gasUsed, err := WasmVM.Instantiate(codeID, env, msgInfo, encodedData, NewStoreAdapter(store), cosmwasm.GoAPI{}, nil, multipliedGasMeter, gasLimit, costJSONDeserialization)
	VMGasRegister.consumeRuntimeGas(ctx, gasUsed)
	return response, err
}

// Calls vm.Execute with internally constructed Gas meter and environment
func callContract(codeID []byte, ctx sdk.Context, store storetypes.KVStore, msg []byte) (*wasmvmtypes.Response, error) {
	sdkGasMeter := ctx.GasMeter()
	multipliedGasMeter := NewMultipliedGasMeter(sdkGasMeter, VMGasRegister)
	gasLimit := VMGasRegister.runtimeGasForContract(ctx)
	chainID := ctx.BlockHeader().ChainID
	height := ctx.BlockHeader().Height
	// safety checks before casting below
	if height < 0 {
		panic("Block height must never be negative")
	}
	nsec := ctx.BlockTime().UnixNano()
	if nsec < 0 {
		panic("Block (unix) time must never be negative ")
	}
	env := wasmvmtypes.Env{
		Block: wasmvmtypes.BlockInfo{
			Height:  uint64(height),
			Time:    uint64(nsec),
			ChainID: chainID,
		},
		Contract: wasmvmtypes.ContractInfo{
			Address: "",
		},
	}

	ctx.GasMeter().ConsumeGas(VMGasRegister.InstantiateContractCosts(len(msg)), "Loading CosmWasm module: execute")
	resp, gasUsed, err := WasmVM.Sudo(codeID, env, msg, NewStoreAdapter(store), wasmvmAPI,
		nil, multipliedGasMeter, gasLimit, costJSONDeserialization)
	VMGasRegister.consumeRuntimeGas(ctx, gasUsed)
	return resp, err
}

// Call vm.Query
func queryContractWithStore(codeID cosmwasm.Checksum, ctx sdk.Context, store storetypes.KVStore, msg []byte) ([]byte, error) {
	sdkGasMeter := ctx.GasMeter()
	multipliedGasMeter := NewMultipliedGasMeter(sdkGasMeter, VMGasRegister)
	gasLimit := VMGasRegister.runtimeGasForContract(ctx)
	chainID := ctx.BlockHeader().ChainID
	height := ctx.BlockHeader().Height
	// safety checks before casting below
	if height < 0 {
		panic("Block height must never be negative")
	}
	nsec := ctx.BlockTime().UnixNano()
	if nsec < 0 {
		panic("Block (unix) time must never be negative ")
	}
	env := wasmvmtypes.Env{
		Block: wasmvmtypes.BlockInfo{
			Height:  uint64(height),
			Time:    uint64(nsec),
			ChainID: chainID,
		},
		Contract: wasmvmtypes.ContractInfo{
			Address: "",
		},
	}
	ctx.GasMeter().ConsumeGas(VMGasRegister.InstantiateContractCosts(len(msg)), "Loading CosmWasm module: query")
	resp, gasUsed, err := WasmVM.Query(codeID, env, msg, NewStoreAdapter(store), cosmwasm.GoAPI{}, nil, multipliedGasMeter, gasLimit, costJSONDeserialization)
	VMGasRegister.consumeRuntimeGas(ctx, gasUsed)
	return resp, err
}

func humanAddress(canon []byte) (string, uint64, error) {
	return "", 0, errors.New("humanAddress not implemented")
}

func canonicalAddress(human string) ([]byte, uint64, error) {
	return nil, 0, errors.New("canonicalAddress not implemented")
}
