package backend

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartServer(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	// Register other routes
	// r.HandleFunc("/another-route", anotherHandler).Methods("POST")

	log.Printf("Backend server started at http://localhost%s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from root!"))
}

// func anotherHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello from another route!"))
// }
