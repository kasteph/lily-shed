package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	layoutISO = "2006-01-02_15-04-05"
	genesis   = 1598306400
	epochUnit = 30
)

var (
	ErrDateParse = errors.New("date cannot be parsed")
	ErrEpoch     = errors.New("epoch given must be >= 0")
)

var ConvertCmd = &cli.Command{
	Name:  "convert",
	Usage: "Convert between epochs and dates",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "short",
			Aliases: []string{"s"},
			Usage:   "return epoch or date only",
			Value:   false,
		},
	},
	Subcommands: []*cli.Command{
		DateToEpochCmd,
		EpochToDateCmd,
	},
}

var DateToEpochCmd = &cli.Command{
	Name:  "date",
	Usage: "converts a list of dates to their respective epochs",
	Action: func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			fmt.Printf("must provide a date with the format of YYYY-MM-DD_HH-SS-MM\n")
			return nil
		}

		for i := 0; i < c.Args().Len(); i++ {
			date := c.Args().Get(i)
			e, err := dateToEpoch(date)
			if err != nil {
				return err
			}

			if c.Bool("short") {
				fmt.Printf("%v\n", e)
			} else {
				fmt.Printf("epoch of %s is %v\n", date, e)
			}

		}

		return nil
	},
}

var EpochToDateCmd = &cli.Command{
	Name:  "epoch",
	Usage: "converts a given epoch to a date",
	Action: func(c *cli.Context) error {

		epoch := c.Args().Get(0)
		p, err := strconv.ParseInt(epoch, 10, 64)
		if err != nil {
			return err
		}
		d, err := epochToDateString(p)
		if err != nil {
			return err
		}

		if c.Bool("short") {
			fmt.Printf("%v\n", d)
		} else {
			fmt.Printf("date of %v is %s\n", epoch, d)
		}

		return nil
	},
}

func dateToEpoch(date string) (int64, error) {
	log.Infof("getting epoch for date: %s", date)

	t, err := time.Parse(layoutISO, date)
	if err != nil {
		log.Warnf("could not parse date: %v", err)
		return 0, ErrDateParse
	}

	return (t.Unix() - genesis) / epochUnit, nil
}

func epochToDateString(epoch int64) (string, error) {
	log.Infof("getting date for epoch: #{epoch}")

	if epoch < 0 {
		return "", ErrEpoch
	}

	ut := genesis + (epochUnit * int64(epoch))
	t := time.Unix(ut, 0).UTC().Format(layoutISO)

	return t, nil
}
