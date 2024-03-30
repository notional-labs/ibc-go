package types

import errorsmod "cosmossdk.io/errors"

const MaxWasmSize = 3 * 1024 * 1024

func ValidateWasmCode(code []byte) error {
	if len(code) == 0 {
		return ErrWasmEmptyCode
	}
	if len(code) > MaxWasmSize {
		return ErrWasmCodeTooLarge
	}

	return nil
}

// ValidateWasmChecksum validates that the checksum is of the correct length
func ValidateWasmChecksum(checksum Checksum) error {
	lenChecksum := len(checksum)
	if lenChecksum == 0 {
		return errorsmod.Wrap(ErrInvalidChecksum, "checksum cannot be empty")
	}
	if lenChecksum != 32 { // sha256 output is 256 bits long
		return errorsmod.Wrapf(ErrInvalidChecksum, "expected length of 32 bytes, got %d", lenChecksum)
	}

	return nil
}
