package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPHandler is a handler for POSTs to / with action payloads
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// TODO 405 Method Not Allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("err: %+v", err))
		return
	}

	action, err := decodeAction(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("err: %+v", err))
		return
	}

	err = action.Execute()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("err: %+v", err))
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, "OK")
}