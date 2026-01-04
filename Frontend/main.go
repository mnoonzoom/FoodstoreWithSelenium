package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	const staticDir = "."

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := staticDir + r.URL.Path
		if r.URL.Path == "/" {
			path = staticDir + "/index.html"
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, path)
	})

	log.Println("Serving frontend at http://localhost:8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
