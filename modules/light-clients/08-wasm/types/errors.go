package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidData              = errorsmod.Register(ModuleName, 1, "invalid data")
	ErrInvalidCodeId            = errorsmod.Register(ModuleName, 2, "invalid code ID")
	ErrInvalidHeader            = errorsmod.Register(ModuleName, 3, "invalid header")
	ErrUnableToUnmarshalPayload = errorsmod.Register(ModuleName, 4, "unable to unmarshal wasm contract return value")
	ErrUnableToInit             = errorsmod.Register(ModuleName, 5, "unable to initialize wasm contract")
	ErrUnableToCall             = errorsmod.Register(ModuleName, 6, "unable to call wasm contract")
	ErrUnableToQuery            = errorsmod.Register(ModuleName, 7, "unable to query wasm contract")
	ErrUnableToMarshalPayload   = errorsmod.Register(ModuleName, 8, "unable to marshal wasm contract payload")
	// Wasm specific
	ErrWasmEmptyCode        = errorsmod.Register(ModuleName, 9, "empty wasm code")
	ErrWasmChecksumNotFound = errorsmod.Register(ModuleName, 10, "wasm checksum not found")
	ErrWasmCodeTooLarge     = errorsmod.Register(ModuleName, 11, "wasm code too large")
	ErrWasmCodeExists       = errorsmod.Register(ModuleName, 12, "wasm code already exists")
	ErrWasmCodeValidation   = errorsmod.Register(ModuleName, 13, "unable to validate wasm code")
	ErrWasmInvalidCode      = errorsmod.Register(ModuleName, 14, "invalid wasm code")
	ErrWasmInvalidCodeID    = errorsmod.Register(ModuleName, 15, "invalid wasm code id")
	ErrWasmCodeIDNotFound   = errorsmod.Register(ModuleName, 16, "wasm code id not found")
	ErrInvalid              = errorsmod.Register(ModuleName, 17, "invalid")
	ErrCreateFailed         = errorsmod.Register(ModuleName, 18, "create wasm contract failed")
)
