package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

//type apiHandler struct{}


type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	
	mux := http.NewServeMux()

	hits:= apiConfig{
		fileserverHits: atomic.Int32{},
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("GET /app/", hits.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("GET /assets/logo.png", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", hits.handlerMetrics)
	mux.HandleFunc("POST /reset", hits.reset)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
		
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		
		//buf := "Ok"
		w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	
		cfg.fileserverHits.Store(0)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits reset: %d", cfg.fileserverHits.Load())))
}