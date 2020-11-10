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

// CreateQRCode : creates a QRCode for a link and send it in request.
func CreateQRCode(w http.ResponseWriter, r *http.Request) {
	lnk := links.Link{Hash: mux.Vars(r)["hash"]}
	globals.Database.First(&lnk)
	pathImage := "./tmp/" + lnk.Hash + ".png"
	_, err := os.Stat(pathImage)
	if os.IsNotExist(err) {
		qrcode.WriteFile(lnk.URL, qrcode.Medium, 256, pathImage)
	} else if err != nil {
		w.WriteHeader(500)
		return
	}
	img, _ := os.Open(pathImage)
	defer img.Close()

	w.Header().Set("Content-Type", "image/png") // <-- set the content-type header
	_, err = io.Copy(w, img)
	if err != nil && !globals.Prod {
		fmt.Println(err)
	}

}
