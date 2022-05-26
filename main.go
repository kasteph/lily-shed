package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var log = logging.Logger("lily-shed")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-interrupt:
			cancel()
		case <-ctx.Done():
		}
	}()

	app := &cli.App{
		Name:  "lily-shed",
		Usage: "smol tools to make working with lily easier",
		Commands: []*cli.Command{
			ConvertCmd,
			SnapshotCmd,
		},
	}

	app.Setup()

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}
