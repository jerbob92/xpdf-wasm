package main

import (
	"context"
	"crypto/rand"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/experimental/logging"
	"github.com/tetratelabs/wazero/imports/emscripten"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed wasm/*
var wasmBinaries embed.FS

func main() {
	ctx := context.WithValue(context.Background(), experimental.FunctionListenerFactoryKey{}, logging.NewLoggingListenerFactory(os.Stdout))
	ctx = context.Background() // Comment this line to get debug information.

	runtimeConfig := wazero.NewRuntimeConfig()
	cache, err := wazero.NewCompilationCacheWithDir(".wazero-cache")
	if err == nil {
		runtimeConfig = runtimeConfig.WithCompilationCache(cache)
	}
	wazeroRuntime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	defer wazeroRuntime.Close(ctx)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, wazeroRuntime); err != nil {
		log.Fatal(err)
	}

	if _, err := emscripten.Instantiate(ctx, wazeroRuntime); err != nil {
		log.Fatal(err)
	}

	availableWASMFiles, err := wasmBinaries.ReadDir("wasm")
	if err != nil {
		log.Fatal(err)
	}

	availableBinaries := make([]string, len(availableWASMFiles))
	for i := range availableWASMFiles {
		availableBinaries[i] = strings.TrimSuffix(availableWASMFiles[i].Name(), ".wasm")
	}

	incorrectStartArgument := func() {
		log.Fatalf("You should minimally start the program with one of the following arguments: %s", strings.Join(availableBinaries, ", "))
	}

	if len(os.Args) < 2 {
		incorrectStartArgument()
	}

	var wasmData []byte
	program := os.Args[1]
	for i := range availableBinaries {
		if availableBinaries[i] == program {
			var err error
			wasmData, err = wasmBinaries.ReadFile(fmt.Sprintf("wasm/%s.wasm", program))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if wasmData == nil {
		incorrectStartArgument()
	}

	compiledModule, err := wazeroRuntime.CompileModule(ctx, wasmData)
	if err != nil {
		log.Fatal(err)
	}

	fsConfig := wazero.NewFSConfig()

	// On Windows we mount the volume of the current working directory as
	// root. On Linux we mount / as root.
	if runtime.GOOS == "windows" {
		cwdDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		volumeName := filepath.VolumeName(cwdDir)
		if volumeName != "" {
			fsConfig = fsConfig.WithDirMount(fmt.Sprintf("%s\\", volumeName), "/")
		}
	} else {
		fsConfig = fsConfig.WithDirMount("/", "/")
	}

	moduleConfig := wazero.NewModuleConfig().
		WithStartFunctions("_start").
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithRandSource(rand.Reader).
		WithFSConfig(fsConfig).
		WithName("").
		WithArgs(os.Args[1:]...)

	_, err = wazeroRuntime.InstantiateModule(ctx, compiledModule, moduleConfig)
	if err != nil {
		log.Fatal(err)
	}
}
