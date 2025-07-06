package main

import (
	//"fmt"
	"log"
	"net/http"
)

//type apiHandler struct{}

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	/*
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
	fmt.Fprintf(w, "Welcome to the home page!")
	})*/

	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
