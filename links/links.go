package links

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"urlq/globals"

	"github.com/gorilla/mux"
)

// Link represents the model for database where links are matched to an hash.
type Link struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	URL       string    `gorm:"unique" json:"url"`
	Hash      string    `gorm:"unique" json:"hash"`
}

// CreateLink : create a Link object and save it in the database
// if it does not already exists, and returns object as json.
func CreateLink(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var obj Link
	json.Unmarshal(reqBody, &obj)
	obj.Hash = hashLink(obj.URL)
	res := globals.Database.FirstOrCreate(&obj, obj)
	if res.Error != nil {
		fmt.Println(res.Error)
		panic("Creation did not work.")
	}
	json.NewEncoder(w).Encode(obj)
	if !globals.Prod {
		fmt.Println("request: creation link")
	}
	return
}

// RetrieveLink : retrieves a Link object given an hash.
func RetrieveLink(w http.ResponseWriter, r *http.Request) {
	var obj Link = Link{Hash: mux.Vars(r)["hash"]}
	globals.Database.Take(&obj, "hash = ?", obj.Hash)
	json.NewEncoder(w).Encode(obj)
	if !globals.Prod {
		fmt.Println("request: retrieve link")
	}
	return
}

// hashLink: takes an url as input and returns its hash.
// Uses CRC64 (ISO) on url and a salt. Returns a string format base36 of the hash.
func hashLink(url string) (res string) {
	res = strconv.FormatInt(int64((crc64.Checksum([]byte(url+"321_SeCReT_KeyFoR_SaLT_123"), crc64.MakeTable(crc64.ISO)))), 36)
	if matched, _ := regexp.Match("-[1-9].*", []byte(res)); matched {
		transform, _ := strconv.Atoi(string(res[1]))
		res = string(byte(65+transform)) + res[2:]
	}
	return
}
