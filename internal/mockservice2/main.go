package mockservice2

import (
	"fmt"
	"net/http"
	"time"
)

// User represents a user with a name and age.
type Res struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func streamingHandler(w http.ResponseWriter, req *http.Request) {
	// check if flusher is supported
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// set headers
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	for i := range 5 {
		fmt.Fprintf(w, "Sending chunk %d...\n", i+1)
		flusher.Flush()
		time.Sleep(2 * time.Second)
	}

	fmt.Fprintf(w, "Streaming complete.\n")

}

func Run() {
	http.HandleFunc("/", streamingHandler)
	http.Handle("/headers", http.HandlerFunc(headers))

	http.ListenAndServe(":9001", nil)
}
