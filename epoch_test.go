package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDateToEpoch(t *testing.T) {
	date := "2022-02-28_08-00-00"
	expectedEpoch := int64(1590960)

	epoch, err := dateToEpoch(date)

	assert.Equal(t, expectedEpoch, epoch, err)
}

func TestDateToEpochWithEmptyString(t *testing.T) {
	epoch, err := dateToEpoch("")
	assert.Equal(t, int64(0), epoch, err)
}

func TestDateToEpochWithInvalidLayout(t *testing.T) {
	epoch, err := dateToEpoch("02-28-2022_09-00-00")
	assert.Equal(t, int64(0), epoch, err)
}

func TestEpochToDate(t *testing.T) {
	epoch := int64(1590960)
	expectedDate := "2022-02-28_08-00-00"

	e, err := epochToDateString(epoch)

	assert.Equal(t, expectedDate, e, err)
}

func TestEpochToDateWithMinusOne(t *testing.T) {
	epoch := int64(-1)

	e, err := epochToDateString(epoch)

	assert.Equal(t, "", e, err)
}
