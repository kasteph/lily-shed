package main

import (
	"bytes"
	"io/ioutil"
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
	o, _ := ioutil.ReadAll(f)
	assert.Equal(t, []byte("snapshot"), o)
	os.Remove("minimal_finality_stateroots_1590960_2022-02-28_08-00-00.car")
}
