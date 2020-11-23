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
	headers := w.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Methods", "POST")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, X-User-ID")
	headers.Add("Access-Control-Max-Age", "86400")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
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

	userID := r.Header.Get("X-User-ID")
	response, err := action.Execute(r.Context(), handler.Firestore, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("err: %+v", err))
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}
