package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

const (
	s3Host         = "https://fil-chain-snapshots-fallback.s3.amazonaws.com/mainnet"
	snapshotPrefix = "minimal_finality_stateroots"
)

var SnapshotCmd = &cli.Command{
	Name:  "snapshot",
	Usage: fmt.Sprintf("Get a minimal Filecoin snapshot from the %s bucket", s3Host),
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "max-attempts",
			Aliases: []string{"m"},
			Usage:   "the max number of attempts to make when finding a snapshot in the bucket",
			Value:   4,
		},
		&cli.StringFlag{
			Name:    "output-path",
			Aliases: []string{"o"},
			Usage:   "the file path to store the snapshot to",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			fmt.Printf("must provide a date with the format of YYYY-MM-DD_HH-SS-MM\n")
			return nil
		}

		date := c.Args().Get(0)

		opts := []func(*client){
			withHost(s3Host),
			withMaxAttempts(c.Int("max-attempts")),
		}

		o := c.String("output-path")
		if o != "" {
			p, err := os.Create(o)
			if err != nil {
				return fmt.Errorf("could not create file: %s", err)
			}
			opts = append(opts, withOutputPath(p))
		}

		cl := newClient(opts...)

		s, err := cl.getSnapshot(date, 0)
		if err != nil {
			log.Error(err)
			s.Close()
			os.Exit(1)
		}
		s.Close()

		return nil
	},
}

type client struct {
	host        string
	maxAttempts int
	outputPath  *os.File
}

func newClient(opts ...func(*client)) *client {
	c := &client{}
	for _, o := range opts {
		o(c)
	}
	return c
}

func withHost(host string) func(*client) {
	return func(c *client) {
		c.host = host
	}
}

func withMaxAttempts(m int) func(*client) {
	return func(c *client) {
		c.maxAttempts = m
	}
}

func withOutputPath(outputPath *os.File) func(*client) {
	return func(c *client) {
		c.outputPath = outputPath
	}
}

func (c *client) getSnapshot(date string, attempt int) (*os.File, error) {
	fmt.Printf("getSnapshot(date: %s, attempt: %d)\n", date, attempt)

	if attempt >= c.maxAttempts {
		return nil, fmt.Errorf("reached max attempts of %d time(s)", c.maxAttempts)
	}

	epoch, err := dateToEpoch(date)
	if err != nil {
		return nil, fmt.Errorf("could not convert date: %s", err)
	}

	car := fmt.Sprintf("%s_%d_%s.car", snapshotPrefix, epoch, date)

	url := fmt.Sprintf(
		"%s/%s",
		c.host,
		car,
	)

	fmt.Printf("getSnapshot(...): url: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Error("error in get request: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warnf("couldn't download: %s_%d_%s", snapshotPrefix, epoch, date)
		attempt++

		t, _ := time.Parse(layoutISO, date)
		return c.getSnapshot(t.Add(time.Hour*1).Format(layoutISO), attempt)
	}

	if err != nil {
		return nil, fmt.Errorf("could not download car file: %s", err)
	}

	var file *os.File
	if c.outputPath == nil {
		var err error
		file, err = os.Create(car)
		if err != nil {
			return nil, fmt.Errorf("could not create file: %s", err)
		}
	} else {
		file = c.outputPath
	}

	bar := progressbar.DefaultBytes(resp.ContentLength, "downloading")
	if _, err := io.Copy(io.MultiWriter(file, bar), resp.Body); err != nil {
		return nil, fmt.Errorf("could not write to file: %s", err)
	}

	return file, nil
}
