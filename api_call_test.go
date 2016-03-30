package main

import (
	"net/http/httptest"
	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func TestApp(t *testing.T) {
	response := &Response {
		path:			"",
		contenttype:	"",
		body: ``,
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		if g, w := r.URL.Path, response.path; g != w {
			t.Errorf("request got path %s, want %s", g, w)
		}

		w.Header().Set("Content-Type", response.contenttype)
		io.WriteString(w, response.body)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
}
