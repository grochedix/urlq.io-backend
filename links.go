package main

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Link represents the model for database where links are matched to an hash.
type Link struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	URL       string    `gorm:"unique" json:"url"`
	Hash      string    `gorm:"unique" json:"hash"`
}

// createLink: create a Link object and save it in the database
// if it does not already exists, and returns object as json.
func createLink(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var obj Link
	json.Unmarshal(reqBody, &obj)
	obj.Hash = hashLink(obj.URL)
	res := database.FirstOrCreate(&obj, obj)
	if res.Error != nil {
		fmt.Println(res.Error)
		panic("Creation did not work.")
	}
	json.NewEncoder(w).Encode(obj)
	if !PROD {
		fmt.Println("request: creation link")
	}
	return
}

// retrieveLink: retrieves a Link object given an hash.
func retrieveLink(w http.ResponseWriter, r *http.Request) {
	var obj Link = Link{Hash: mux.Vars(r)["hash"]}
	database.Take(&obj, "hash = ?", obj.Hash)
	json.NewEncoder(w).Encode(obj)
	if !PROD {
		fmt.Println("request: retrieve link")
	}
	return
}

// hashLink: takes an url as input and returns its hash.
// Uses CRC64 (ISO) on url and a salt. Returns a string format base36 of the hash.
func hashLink(url string) string {
	return strconv.FormatInt(int64((crc64.Checksum([]byte(url+"321_SeCReT_KeyFoR_SaLT_123"), crc64.MakeTable(crc64.ISO)))), 36)
}
