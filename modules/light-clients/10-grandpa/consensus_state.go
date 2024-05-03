package grandpa

import (
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

// The implementation of ConsensusState is needed to register the protobuf type, but none of the functions is actually
// used anywhere

func (m ConsensusState) ClientType() string {
	panic("unimplemented")
}

func (m ConsensusState) GetTimestamp() uint64 {
	panic("unimplemented")
}

func (m ConsensusState) ValidateBasic() error {
	panic("unimplemented")
}
