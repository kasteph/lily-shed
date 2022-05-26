package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			fmt.Printf("must provide a date with the format of YYYY-MM-DD_HH-SS-MM\n")
			return nil
		}

		date := c.Args().Get(0)

		cl := newClient(s3Host, c.Int("max-attempts"))

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
}

func newClient(host string, maxAttempts int) client {
	return client{host: host, maxAttempts: maxAttempts}
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

	file, err := os.Create(car)
	if err != nil {
		return nil, fmt.Errorf("could not create file: %s", err)
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		return nil, fmt.Errorf("could not write to file: %s", err)
	}

	return file, nil
}
