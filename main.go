package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	readURLs()
	go writeURLs()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func readURLs() {
	f, err := os.Open("./urls.json")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.NewDecoder(f).Decode(&urls)
	if err != nil {
		log.Println(err)
	}
}

func writeURLs() {
	defer os.Exit(0)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	f, err := os.Create("./urls.json")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.NewEncoder(f).Encode(&urls)
	if err != nil {
		log.Println(err)
	}
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
