package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (m ConsensusState) ClientType() string {
	return Wasm
}

func (m ConsensusState) GetTimestamp() uint64 {
	return m.Timestamp
}

func (m ConsensusState) ValidateBasic() error {
	if m.Timestamp == 0 {
		return errorsmod.Wrap(clienttypes.ErrInvalidConsensus, "timestamp cannot be zero Unix time")
	}

	if m.Data == nil || len(m.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	return nil
}

// NewConsensusState creates a new ConsensusState instance.
func NewConsensusState(data []byte) *ConsensusState {
	return &ConsensusState{
		Data: data,
	}
}
