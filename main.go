package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/experimental/logging"
	"github.com/tetratelabs/wazero/imports/emscripten"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"log"
	"os"
)

//go:embed pdftops.wasm
var wasmFile []byte

func main() {
	ctx := context.WithValue(context.Background(), experimental.FunctionListenerFactoryKey{}, logging.NewLoggingListenerFactory(os.Stdout))
	ctx = context.Background()
	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		runtime.Close(ctx)
		log.Fatal(err)
	}

	if _, err := emscripten.Instantiate(ctx, runtime); err != nil {
		runtime.Close(ctx)
		log.Fatal(err)
	}

	compiledModule, err := runtime.CompileModule(ctx, wasmFile)
	if err != nil {
		runtime.Close(ctx)
		log.Fatal(err)
	}

	moduleConfig := wazero.NewModuleConfig().
		WithStartFunctions("_start").
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithRandSource(rand.Reader).
		WithFSConfig(wazero.NewFSConfig().WithDirMount("/", "/")).
		WithName("").
		WithArgs(os.Args...)

	_, err = runtime.InstantiateModule(ctx, compiledModule, moduleConfig)
	if err != nil {
		log.Fatal(err)
	}
}
