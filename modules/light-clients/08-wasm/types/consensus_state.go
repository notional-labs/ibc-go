package types

import (
	errorsmod "cosmossdk.io/errors"
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
	if len(m.Data) == 0 {
		return errorsmod.Wrap(ErrInvalidData, "data cannot be empty")
	}

	return nil
}

// NewConsensusState creates a new ConsensusState instance.
func NewConsensusState(data []byte) *ConsensusState {
	return &ConsensusState{
		Data: data,
	}
}
