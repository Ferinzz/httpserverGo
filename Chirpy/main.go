package main

import (
	//"fmt"
	"log"
	"net/http"
)

//type apiHandler struct{}

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()


	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/assets/logo.png", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/healthz", handlerReadiness)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
		
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		
		//buf := "Ok"
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}