package main

import (
	"context"
	"github.com/tensved/snet-matrix-framework/pkg/lib"
)

func main() {

	snetEngine := lib.DefaultSNETEngine()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start engine
	if err := snetEngine.Run(ctx); err != nil {
		panic(err)
	}
}
