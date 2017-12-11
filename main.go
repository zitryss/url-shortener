package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
)

var urls = make(map[string]string)
var m = sync.Mutex{}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		longURL := req.FormValue("url")
		shortURL := "http://localhost:8080/" + hashID(8)
		m.Lock()
		urls[shortURL] = longURL
		m.Unlock()
		fmt.Fprintln(resp, shortURL)
	case "GET":
		shortURL := "http://localhost:8080" + req.URL.Path
		m.Lock()
		longURL, ok := urls[shortURL]
		m.Unlock()
		if !ok {
			http.Error(resp, "Page not found.", http.StatusNotFound)
			return
		}
		http.Redirect(resp, req, longURL, http.StatusFound)
	default:
		http.Error(resp, "Method is not supported.", http.StatusInternalServerError)
	}
}

func hashID(len int) string {
	b := make([]byte, len/2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
