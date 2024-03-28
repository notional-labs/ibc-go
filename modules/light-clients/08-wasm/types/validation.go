package types

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
