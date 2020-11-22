package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/firestore"
)

// HTTPHandler is a handler for POSTs to / with action payloads
type HTTPHandler struct {
	Firestore *firestore.Client
}

func (handler *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	err = action.Execute(r.Context(), handler.Firestore)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("err: %+v", err))
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, "OK")
}
