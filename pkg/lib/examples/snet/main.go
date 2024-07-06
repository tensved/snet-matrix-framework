package main

import (
	"context"
	"matrix-ai-framework/pkg/lib"
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
