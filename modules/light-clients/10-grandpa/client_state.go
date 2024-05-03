package grandpa

import (
	"cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

// The implementation of ClientState is needed to register the protobuf type, but none of the functions is actually
// used anywhere

func (m *ClientState) ClientType() string {
	panic("unimplemented")
}

func (m *ClientState) GetLatestHeight() exported.Height {
	panic("unimplemented")
}

func (m *ClientState) Validate() error {
	panic("unimplemented")
}

func (m *ClientState) Status(_ types.Context, _ types.KVStore, _ codec.BinaryCodec) exported.Status {
	panic("unimplemented")
}

func (m *ClientState) ExportMetadata(_ types.KVStore) []exported.GenesisMetadata {
	panic("unimplemented")
}

func (m *ClientState) ZeroCustomFields() exported.ClientState {
	panic("unimplemented")
}

func (m *ClientState) GetTimestampAtHeight(_ types.Context, _ types.KVStore, _ codec.BinaryCodec, _ exported.Height) (uint64, error) {
	panic("unimplemented")
}

func (m *ClientState) Initialize(_ types.Context, _ codec.BinaryCodec, _ types.KVStore, _ exported.ConsensusState) error {
	panic("unimplemented")
}

func (m *ClientState) VerifyMembership(_ types.Context, _ types.KVStore, _ codec.BinaryCodec, _ exported.Height, _ uint64, _ uint64, _ []byte, _ exported.Path, _ []byte) error {
	panic("unimplemented")
}

func (m *ClientState) VerifyNonMembership(_ types.Context, _ types.KVStore, _ codec.BinaryCodec, _ exported.Height, _ uint64, _ uint64, _ []byte, _ exported.Path) error {
	panic("unimplemented")
}

func (m *ClientState) VerifyClientMessage(_ types.Context, _ codec.BinaryCodec, _ types.KVStore, _ exported.ClientMessage) error {
	panic("unimplemented")
}

func (m *ClientState) CheckForMisbehaviour(_ types.Context, _ codec.BinaryCodec, _ types.KVStore, _ exported.ClientMessage) bool {
	panic("unimplemented")
}

func (m *ClientState) UpdateStateOnMisbehaviour(_ types.Context, _ codec.BinaryCodec, _ types.KVStore, _ exported.ClientMessage) {
	panic("unimplemented")
}

func (m *ClientState) UpdateState(_ types.Context, _ codec.BinaryCodec, _ types.KVStore, _ exported.ClientMessage) []exported.Height {
	panic("unimplemented")
}

func (m *ClientState) CheckSubstituteAndUpdateState(_ types.Context, _ codec.BinaryCodec, _, _ types.KVStore, _ exported.ClientState) error {
	panic("unimplemented")
}

func (m *ClientState) VerifyUpgradeAndUpdateState(_ types.Context, _ codec.BinaryCodec, _ types.KVStore, _ exported.ClientState, _ exported.ConsensusState, _, _ []byte) error {
	panic("unimplemented")
}
