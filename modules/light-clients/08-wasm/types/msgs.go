package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"

	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
)

var TypeMsgPushNewWasmCode = "push_wasm_code"
var _ sdk.Msg = &MsgPushNewWasmCode{}

var TypeMsgUpdateWasmCodeId = "update_wasm_code_id"
var _ sdk.Msg = &MsgUpdateWasmCodeId{}

// NewMsgPushNewWasmCode creates a new MsgPushNewWasmCode instance
//
//nolint:interfacer

// Route Implements Msg.
func (msg MsgPushNewWasmCode) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgPushNewWasmCode) Type() string { return TypeMsgPushNewWasmCode }

func NewMsgPushNewWasmCode(signer string, code []byte) *MsgPushNewWasmCode {
	return &MsgPushNewWasmCode{
		Signer: signer,
		Code:   code,
	}
}

func (m MsgPushNewWasmCode) ValidateBasic() error {
	if len(m.Code) == 0 {
		return errorsmod.Wrapf(ErrWasmEmptyCode,
			"empty wasm code",
		)
	}

	return nil
}

func (m MsgPushNewWasmCode) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgPushNewWasmCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// NewMsgUpdateWasmCodeId creates a new MsgUpdateWasmCodeId instance
//
//nolint:interfacer

// Route Implements Msg.
func (msg MsgUpdateWasmCodeId) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgUpdateWasmCodeId) Type() string { return TypeMsgUpdateWasmCodeId }

func NewMsgUpdateWasmCodeId(signer string, codeId []byte, clientId string) *MsgUpdateWasmCodeId {
	return &MsgUpdateWasmCodeId{
		Signer:   signer,
		CodeId:   codeId,
		ClientId: clientId,
	}
}

func (m MsgUpdateWasmCodeId) ValidateBasic() error {
	if len(m.CodeId) != tmhash.Size {
		return errorsmod.Wrapf(ErrWasmEmptyCode,
			"invalid code id length (expected 32, got %d)", len(m.CodeId),
		)
	}

	err := host.ClientIdentifierValidator(m.ClientId)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgUpdateWasmCodeId) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateWasmCodeId) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
