package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"urlq/globals"
	"urlq/links"
	"urlq/qrcodes"
	"urlq/settings"

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
	r.HandleFunc("/image/{hash:-?[0-9a-zA-Z]+}", qrcodes.CreateQRCode).Methods("GET", "OPTIONS")

	link := r.PathPrefix("/link").Subrouter().StrictSlash(true)
	link.Use(settings.Middleware)
	link.HandleFunc("/", links.CreateLink).Methods("POST", "OPTIONS")
	link.HandleFunc("/{hash:-?[0-9a-zA-Z]+}", links.RetrieveLink).Methods("GET", "OPTIONS")

	log.Fatal(http.ListenAndServe(":9999", handlers.CORS(handlers.AllowedHeaders(
		[]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(r)))
}

// main : main function, where program launches.
func main() {

	if !globals.Prod {
		fmt.Println("Dev API starting to run.")
	}

	settings.ConnectDB()

	if globals.DBerr != nil {
		fmt.Println(globals.DBerr)
		return
	}

	if !globals.MigrationDone || (len(os.Args) > 1 && os.Args[1] == "migrate") {
		err := settings.Migrate()
		if err != nil {
			panic("Migration did not work.")
		} else {
			globals.MigrationDone = true
			fmt.Println("Migration done!")
		}
		return
	}

	myRouter()

}
