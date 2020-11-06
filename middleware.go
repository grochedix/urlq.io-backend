package main

import (
	"net/http"
	"strings"
	"time"
)

// RequestIP : for requests rate limitation
type RequestIP struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	LastRequest time.Time `gorm:"autoUpdateTime"`
	IP          string    `gorm:"uniqueIndex"`
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !checkIPaddr(r) {
			w.WriteHeader(429)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func checkIPaddr(r *http.Request) (access bool) {
	ip := r.RemoteAddr
	ip = ip[1:strings.Index(ip, "]")]
	obj := RequestIP{IP: ip}
	res := DB.First(&obj)
	if res.RowsAffected == 0 {
		DB.Create(&obj)
		access = true
	} else {
		if time.Now().Sub(obj.LastRequest) > 2*time.Second {
			access = true
		}
		obj.LastRequest = time.Now()
		DB.Save(&obj)
	}
	return access
}
