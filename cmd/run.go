package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/craiggwilson/go-typeproviders/pkg/generate"
)

func run(p generate.StructProvider) {

	ctx := signalContext(context.Background())
	pkg := rootCmd.PersistentFlags().Lookup("pkg").Value.String()
	embedStructs, err := strconv.ParseBool(rootCmd.PersistentFlags().Lookup("embedStructs").Value.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = generate.Generate(ctx, p, "", pkg, embedStructs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func signalContext(ctx context.Context) context.Context {
	signalCtx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-c:
			signal.Stop(c)
			cancel()
		case <-ctx.Done():
			signal.Stop(c)
			cancel()
		}
	}()
	return signalCtx
}
