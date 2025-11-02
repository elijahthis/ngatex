package mockservice

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Res struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello Service 1")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(Res{
		Message: "Service 1 is running successfully",
		Success: true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func Run() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Mock 1 is alive"))
	})
	http.Handle("/headers", http.HandlerFunc(headers))

	port := flag.String("port", "9001", "port to listen to")
	flag.Parse()

	fmt.Printf("MOCKSERVICE 1 is listening on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}

func RunMultiple(port string, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.Handle("/headers", http.HandlerFunc(headers))

	fmt.Printf("MOCKSERVICE 1 is listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server on port %s failed: %v\n", port, err)
	}
}
