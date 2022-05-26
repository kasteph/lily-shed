package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSnapshot(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("snapshot"))
	}))
	defer ts.Close()

	c := newClient(ts.URL, 4)
	s, _ := c.getSnapshot("2022-02-28_08-00-00", 0)
	f, _ := os.OpenFile(s.Name(), os.O_RDONLY, 0755)
	o, _ := ioutil.ReadAll(f)
	assert.Equal(t, []byte("snapshot"), o)
}
