package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// homepage : test function for connecting to the API.
func homepage(w http.ResponseWriter, r *http.Request) {
	q := make(map[string]string)
	q["message"] = "You reached the API! :)"
	json.NewEncoder(w).Encode(q)
}

// myRouter : matching routes to function.
func myRouter() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homepage).Methods("GET", "OPTIONS")
	r.Use(middleware)
	r.HandleFunc("/link", createLink).Methods("POST", "OPTIONS")
	r.HandleFunc("/link/{hash:-?[0-9a-z]+}", retrieveLink).Methods("GET", "OPTIONS")
	log.Fatal(http.ListenAndServe(":10000", handlers.CORS(handlers.AllowedHeaders(
		[]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(r)))
}

// main : main function, where program launches.
func main() {
	
	if !PROD {
		fmt.Println("Dev API starting to run.")
	}

	connectDB()

	if DBerr != nil {
		fmt.Println(DBerr)
		return
	}

	if !MigrationDone || (len(os.Args) > 1 && os.Args[1] == "migrate") {
		err := migrate()
		if err != nil {
			panic("Migration did not work.")
		} else {
			MigrationDone = true
			fmt.Println("Migration done!")
		}
		return
	}

	myRouter()

}
