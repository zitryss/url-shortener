package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"net/http"
	"sync"
)

var domain string
var port string
var m sync.RWMutex
var urls map[string]string

func init() {
	flag.StringVar(&domain, "domain", "localhost", "")
	flag.StringVar(&port, "port", "8080", "")
	urls = make(map[string]string)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		longURL := req.FormValue("url")
		shortURL := "http://" + domain + ":" + port + "/" + hashID(8)
		m.Lock()
		urls[shortURL] = longURL
		m.Unlock()
		fmt.Fprintln(resp, shortURL)
	case "GET":
		shortURL := "http://" + domain + ":" + port + req.URL.Path
		m.RLock()
		longURL, ok := urls[shortURL]
		m.RUnlock()
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
