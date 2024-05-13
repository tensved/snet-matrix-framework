package main

import (
	"context"
	"matrix-ai-framework/pkg/lib"
)

func main() {

	snetEngine := lib.DefaultSNETEngine()

	ctx, _ := context.WithCancel(context.Background())

	// start engine
	if err := snetEngine.Run(ctx); err != nil {
		panic(err)
	}
}
