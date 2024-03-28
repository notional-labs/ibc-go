package types

import (
	"fmt"

	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ClientMessage = &Header{}

func (m Header) ClientType() string {
	return Wasm
}

func (m Header) ValidateBasic() error {
	if m.Data == nil || len(m.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	return nil
}
