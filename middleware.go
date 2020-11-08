package main

import (
	"net/http"
	"strings"
	"time"
)

// RequestIP : matching an IP and its last request.
type RequestIP struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	LastRequest time.Time `gorm:"autoUpdateTime"`
	IP          string    `gorm:"uniqueIndex"`
}

// middleware : implementing middleware for the REST API.
func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !checkIPaddr(r) {
			w.WriteHeader(429)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

// checkIPaddr: allow a request every 0.2 second for a same user.
func checkIPaddr(r *http.Request) (access bool) {
	ip := r.RemoteAddr
	ip = ip[1:strings.Index(ip, "]")]
	obj := RequestIP{IP: ip}
	res := database.First(&obj)
	if res.RowsAffected == 0 {
		database.Create(&obj)
		access = true
	} else {
		if time.Now().Sub(obj.LastRequest) > 200*time.Millisecond {
			access = true
		}
		obj.LastRequest = time.Now()
		database.Save(&obj)
	}
	return access
}
