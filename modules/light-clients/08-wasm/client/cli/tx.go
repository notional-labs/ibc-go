package cli

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	types "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// newStoreCodeCmd returns a message to store wasm bytes code
func newStoreCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "store-code [path/to/wasm-file]",
		Short:   "Reads wasm code from the file and push wasm code to chain",
		Long:    "Reads wasm code from the file and push wasm code to chain",
		Example: fmt.Sprintf("%s tx %s-wasm store-code [path/to/wasm_file]", version.AppName, ibcexported.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			code, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgStoreCode{
				Signer:       clientCtx.GetFromAddress().String(),
				WasmByteCode: code,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Create a transaction from the message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func newMigrateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate-contract [client-id] [checksum] [migrate-msg]",
		Short:   "Migrates a contract to a new byte code",
		Long:    "Migrates the contract for the specified client ID to the byte code corresponding to checksum, passing the JSON-encoded migrate message to the contract",
		Example: fmt.Sprintf("%s tx %s-wasm migrate-contract 08-wasm-0 b3a49b2914f5e6a673215e74325c1d153bb6776e079774e52c5b7e674d9ad3ab {}", version.AppName, ibcexported.ModuleName),
		Args:    cobra.ExactArgs(3), // Ensure exactly three arguments are passed
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			clientID := args[0]
			checksumHex := args[1]
			migrateMsg := args[2]

			checksum, err := hex.DecodeString(checksumHex)
			if err != nil {
				return fmt.Errorf("invalid checksum format: %w", err)
			}

			// Construct the message
			msg := &types.MsgMigrateContract{
				Signer:   clientCtx.GetFromAddress().String(),
				ClientId: clientID,
				Checksum: []byte(checksum),
				Msg:      []byte(migrateMsg),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// Generate or broadcast the transaction
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
