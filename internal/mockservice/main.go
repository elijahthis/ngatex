package mockservice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User represents a user with a name and age.
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
	http.Handle("/headers", http.HandlerFunc(headers))

	http.ListenAndServe(":9000", nil)
}
