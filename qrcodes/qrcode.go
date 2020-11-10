package qrcodes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"urlq/globals"
	"urlq/links"

	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

// CreateQRCode : creates a QRCode for a link.
func CreateQRCode(w http.ResponseWriter, r *http.Request) {
	pathImage := "./tmp/"
	hash := mux.Vars(r)["hash"]
	lnk := links.Link{Hash: hash}
	globals.Database.First(&lnk)
	_, err := os.Stat(pathImage + lnk.Hash + ".png")
	if os.IsNotExist(err) {
		qrcode.WriteFile(lnk.URL, qrcode.Medium, 256, pathImage+lnk.Hash+".png")
	} else if err != nil {
		w.WriteHeader(500)
		return
	}
	img, _ := os.Open(pathImage + lnk.Hash + ".png")
	defer img.Close()

	w.Header().Set("Content-Type", "image/png") // <-- set the content-type header
	_, err = io.Copy(w, img)
	if err != nil && !globals.Prod {
		fmt.Println(err)
	}

}
