package ibcwasm

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"errors"
)

var (
	// state management
	Checksums    collections.KeySet[[]byte]
	Schema       collections.Schema
	ChecksumsKey = collections.NewPrefix(0)
	queryPlugins QueryPluginsI
	queryRouter  QueryRouter
)

// SetupWasmStoreService sets up the 08-wasm module's collections.
func SetupWasmStoreService(storeService storetypes.KVStoreService) {
	sb := collections.NewSchemaBuilder(storeService)

	Checksums = collections.NewKeySet(sb, ChecksumsKey, "checksums", collections.BytesKey)

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	Schema = schema
}

// SetQueryRouter sets the custom wasm query router for the 08-wasm module.
// Panics if the queryRouter is nil.
func SetQueryRouter(router QueryRouter) {
	if router == nil {
		panic(errors.New("query router must be not nil"))
	}
	queryRouter = router
}

// GetQueryRouter returns the custom wasm query router for the 08-wasm module.
func GetQueryRouter() QueryRouter {
	return queryRouter
}

// SetQueryPlugins sets the current query plugins
func SetQueryPlugins(plugins QueryPluginsI) {
	if plugins == nil {
		panic(errors.New("query plugins must be not nil"))
	}
	queryPlugins = plugins
}

// GetQueryPlugins returns the current query plugins
func GetQueryPlugins() QueryPluginsI {
	return queryPlugins
}
