package mockservice2

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

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
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

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
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Mock 2 is alive"))
	})
	http.Handle("/headers", http.HandlerFunc(headers))

	port := flag.String("port", "9004", "port to listen to")
	flag.Parse()

	fmt.Printf("MOCKSERVICE 2 is listening on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}

func RunMultiple(port string, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", streamingHandler)
	mux.Handle("/headers", http.HandlerFunc(headers))

	fmt.Printf("MOCKSERVICE 2 is listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server on port %s failed: %v\n", port, err)
	}
}
