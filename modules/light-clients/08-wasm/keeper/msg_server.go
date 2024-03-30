package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
)

var _ types.MsgServer = Keeper{}

// PushNewWasmCode defines a rpc handler method for MsgPushNewWasmCode
func (k Keeper) PushNewWasmCode(goCtx context.Context, msg *types.MsgPushNewWasmCode) (*types.MsgPushNewWasmCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Signer {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority: expected %s, got %s", k.authority, msg.Signer)
	}

	codeID, err := k.storeWasmCode(ctx, msg.Code, types.WasmVM.StoreCode)
	if err != nil {
		return nil, errorsmod.Wrap(err, "pushing new wasm code failed")
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePushWasmCode,
			sdk.NewAttribute(types.AttributeKeyWasmCodeID, hex.EncodeToString(codeID)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, clienttypes.AttributeValueCategory),
		),
	})

	return &types.MsgPushNewWasmCodeResponse{
		CodeId: codeID,
	}, nil
}

// UpdateWasmCodeId defines a rpc handler method for MsgUpdateWasmCodeId
func (k Keeper) UpdateWasmCodeId(goCtx context.Context, msg *types.MsgUpdateWasmCodeId) (*types.MsgUpdateWasmCodeIdResponse, error) {
	if k.authority != msg.Signer {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority: expected %s, got %s", k.authority, msg.Signer)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	codeId := msg.CodeId

	codeIDKey := types.CodeID(codeId)
	if !types.HasChecksum(ctx, codeIDKey) {
		return nil, errorsmod.Wrapf(types.ErrInvalidCodeId, "code id %s does not exist", hex.EncodeToString(codeId))
	}

	clientId := msg.ClientId
	unknownClientState, found := k.clientKeeper.GetClientState(ctx, clientId)
	if !found {
		return nil, errorsmod.Wrapf(clienttypes.ErrClientNotFound, "cannot update client with ID %s", clientId)
	}

	clientState, ok := unknownClientState.(*types.ClientState)
	if !ok {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "client state type %T, expected %T", unknownClientState, (*types.ClientState)(nil))
	}

	clientState.CodeId = codeId

	k.clientKeeper.SetClientState(ctx, clientId, clientState)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateWasmCodeId,
			sdk.NewAttribute(clienttypes.AttributeKeyClientID, clientId),
			sdk.NewAttribute(types.AttributeKeyWasmCodeID, hex.EncodeToString(codeId)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, clienttypes.AttributeValueCategory),
		),
	})

	return &types.MsgUpdateWasmCodeIdResponse{
		ClientId: clientId,
		CodeId:   codeId,
	}, nil
}
