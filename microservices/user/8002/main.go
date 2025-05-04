package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/user/hello", helloHandler)
	mux.HandleFunc("/user/health", healthHandler)

	port := ":8002"
	log.Printf("Server starting on http://localhost%s ...", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, User!"})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
