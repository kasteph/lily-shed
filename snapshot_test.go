package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetSnapshot(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=test")
		w.Header().Set("Accept-Ranges", "bytes")
		//w.Write([]byte("snapshot"))

		contents := []byte("snapshot")
		http.ServeContent(w, r, "test", time.Now(), bytes.NewReader(contents))
	}))
	defer ts.Close()

	c := newClient(withHost(ts.URL), withMaxAttempts(4))
	s, _ := c.getSnapshot("2022-02-28_08-00-00", 0)
	f, _ := os.OpenFile(s.Name(), os.O_RDONLY, 0755)
	o, _ := io.ReadAll(f)
	assert.Equal(t, []byte("snapshot"), o)
	os.Remove("minimal_finality_stateroots_1590960_2022-02-28_08-00-00.car")
}

func TestGetSnapshotPrintToStdout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=test")
		w.Header().Set("Accept-Ranges", "bytes")
		//w.Write([]byte("snapshot"))

		contents := []byte("snapshot")
		http.ServeContent(w, r, "test", time.Now(), bytes.NewReader(contents))
	}))
	defer ts.Close()

	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := newClient(withHost(ts.URL), withMaxAttempts(4), withPrint(true))
	s, _ := c.getSnapshot("2022-02-28_08-00-00", 0)
	w.Close()
	os.Stdout = stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	assert.Nil(t, s)
	assert.Equal(
		t,
		ts.URL+"/minimal_finality_stateroots_1590960_2022-02-28_08-00-00.car",
		buf.String(),
	)
}
