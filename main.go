package main

import (
	"context"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var log = logging.Logger("lily-shed")

func main() {
	ctx := context.Background()

	app := &cli.App{
		Name:  "lily-shed",
		Usage: "smol tools to make working with lily easier",
		Commands: []*cli.Command{
			ConvertCmd,
			SnapshotCmd,
		},
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}
