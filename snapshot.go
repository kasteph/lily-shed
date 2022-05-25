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
	bucket         = "fil-chain-snapshots-fallback"
	snapshotPrefix = "minimal_finality_stateroots"
	maxAttempts    = 4
)

var SnapshotCmd = &cli.Command{
	Name:  "snapshot",
	Usage: fmt.Sprintf("Get a minimal Filecoin snapshot from the %s bucket", bucket),
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "attempts",
			Aliases: []string{"a"},
			Usage:   "the max number of attempts to make when finding a snapshot in the bucket",
			Value:   maxAttempts,
		},
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			fmt.Printf("must provide a date with the format of YYYY-MM-DD_HH-SS-MM\n")
			return nil
		}

		date := c.Args().Get(0)

		cl := newClient("https://%s.s3.amazonaws.com/mainnet/%s")

		s, err := cl.getSnapshot(date, c.Int("attempts"))
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
	host string
}

func newClient(host string) client {
	return client{host: host}
}

func (c *client) getSnapshot(date string, attempt int) (fp *os.File, err error) {
	if attempt >= maxAttempts {
		return nil, fmt.Errorf("reached max attempts of %d time(s)", maxAttempts)
	}

	epoch, err := dateToEpoch(date)
	if err != nil {
		return nil, fmt.Errorf("could not convert date: %s", err)
	}

	car := fmt.Sprintf("%s_%d_%s.car", snapshotPrefix, epoch, date)

	url := fmt.Sprintf(
		"%s/%s/%s",
		c.host,
		bucket,
		car,
	)

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
