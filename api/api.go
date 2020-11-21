package api

import (
	"io"
	"net/http"
)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, `Hello, SSL-API-world!`)
}
