package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/spf13/cobra"
)

// newPushNewWasmCodeCmd returns the command to create a PushNewWasmCode transaction
func newPushNewWasmCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push-wasm [wasm-file]",
		Short: "Reads wasm code from the file and creates push transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fileName := args[0]

			code, err := ioutil.ReadFile(fileName)
			if err != nil {
				return err
			}

			msg := &types.MsgPushNewWasmCode{
				Code:   code,
				Signer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// newUpdateWasmCodeId returns the command to create a UpdateWasmCodeId transaction
func newUpdateWasmCodeId() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-wasm-code-id [client-id] [code-id]",
		Short: "Updates wasm code id for a client",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			clientId := args[0]
			codeId, err := transfertypes.ParseHexHash(args[1])

			if err != nil {
				return err
			}

			msg := &types.MsgUpdateWasmCodeId{
				ClientId: clientId,
				CodeId:   codeId,
				Signer:   clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
