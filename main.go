package main

import (
	"context"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var glog = logging.Logger("lily-shed")

func main() {
	err := logging.SetLogLevel("lily-shed:snapshot", "info")
	if err != nil {
		glog.Errorf(err.Error())
	}

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
		glog.Fatal(err.Error())
		os.Exit(1)
	}
}
