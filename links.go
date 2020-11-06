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

// Link represents the model for DB where links are matched to an hash.
type Link struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	URL       string    `gorm:"unique" json:"url"`
	Hash      string    `gorm:"unique" json:"hash"`
}

func createLink(w http.ResponseWriter, r *http.Request) {
	if !PROD {
		fmt.Println("request: creation link")
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var obj Link
	json.Unmarshal(reqBody, &obj)
	obj.Hash = hashLink(obj.URL)
	res := DB.FirstOrCreate(&obj, obj)
	if res.Error != nil {
		fmt.Println(res.Error)
		panic("Creation did not work.")
	}
	json.NewEncoder(w).Encode(obj)
	return
}

func retrieveLink(w http.ResponseWriter, r *http.Request) {
	if !PROD {
		fmt.Println("request: retrieve link")
	}
	var obj Link = Link{Hash: mux.Vars(r)["hash"]}
	DB.Take(&obj, "hash = ?", obj.Hash)
	json.NewEncoder(w).Encode(obj)
	return
}

func hashLink(url string) string {
	return strconv.FormatInt(int64((crc64.Checksum([]byte(url+"321_SeCReT_KeyFoR_SaLT_123"), crc64.MakeTable(crc64.ISO)))), 36)
}
