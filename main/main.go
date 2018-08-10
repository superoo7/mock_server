package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type TestFunc struct {
	Hello string
}

type DataFunc struct {
	Body   string
	Url    string
	Status string
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/test", testFunc).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(getFunc)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	if err := http.ListenAndServeTLS(":8000", "server.crt", "server.key", handler); err != nil {
		log.Fatal(err)
	}
}

func getFunc(w http.ResponseWriter, r *http.Request) {
	// enableCors(&w)
	path := r.URL.Path

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}
	bodyString := string(body)

	fmt.Println("GET", path, bodyString)
	data := DataFunc{bodyString, path, "OK"}
	respondWithJSON(w, http.StatusOK, &data)
}

// testFunc a sample route for responding to API
func testFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test route")
	respondWithJSON(w, http.StatusOK, &TestFunc{"test"})
}

// responsWithJson will reply to client in JSON format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, error := json.Marshal(payload)
	if error != nil {
		panic(error)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError will reply to client with error message with {"error": err}
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
