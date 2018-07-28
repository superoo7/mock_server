package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type TestFunc struct {
	Hello string
	Test  string
}

type DataFunc struct {
	Body   string
	Url    string
	Status string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getFunc).Methods("GET", "POST")
	r.HandleFunc("/test", testFunc).Methods("GET")
	r.HandleFunc("/{any}", getFunc).Methods("GET", "POST")
	r.NotFoundHandler = http.HandlerFunc(getFunc)

	if err := http.ListenAndServe(":3001", r); err != nil {
		log.Fatal(err)
	}
}

func getFunc(w http.ResponseWriter, r *http.Request) {
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
	respondWithJSON(w, http.StatusOK, &TestFunc{"test", "b"})
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
