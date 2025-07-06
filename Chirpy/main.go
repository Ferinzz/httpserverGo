package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"encoding/json"
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
	mux.HandleFunc("GET /admin/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", hits.reset)
	mux.HandleFunc("GET /admin/metrics", hits.handlerAdminMetrics)

	mux.HandleFunc("POST /api/validate_chirp", ValidateJSON)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
		
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		
		w.Write([]byte(http.StatusText(http.StatusOK)))
}

//**********\\
//***APIs***\\
//**********\\

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



func (cfg *apiConfig) handlerAdminMetrics(w http.ResponseWriter, req *http.Request) {
	
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load())))
}

func ValidateJSON(w http.ResponseWriter, req *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}

	type jsonError struct {
		Error string `json:"error"`
	}
	
	decoder:= json.NewDecoder(req.Body)
	params:= parameter{}
	err:= decoder.Decode(&params)

	if err != nil {
		myError:= jsonError{
			Error: "Something went wrong",
		}
		log.Printf("Error decoding parameters: %s", err)
		w.Header().Add("Content-type", "text/json")
		w.WriteHeader(500)
		dat, err:= json.Marshal(myError)
		if err != nil {
			log.Printf("Failed to marhsal error message")
			return
		}
		w.Write(dat)
		return
	}

	log.Printf("%v",len(params.Body))
	if len(params.Body) > 140 {
		
		myError:= jsonError{
			Error: "Chirp too long",
		}
		log.Printf("Chirp is too long")
		w.Header().Add("Content-type", "text/json")
		w.WriteHeader(400)
		dat, err:= json.Marshal(myError)
		if err != nil {
			log.Printf("Failed to marhsal error message")
			return
		}
		w.Write(dat)
		return
	}

	type success struct {
		Valid bool `json:"valid"`
	}
	aSuccess:= success{
		Valid: true,
	}
	log.Printf("No errors here")
	w.Header().Add("Content-type", "text/json")
	w.WriteHeader(200)
	dat, err:= json.Marshal(aSuccess)
	if err != nil {
		log.Printf("Failed to marhsal success message")
		return
	}
	w.Write(dat)
	return
}