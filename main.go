// +build linux darwin

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
	"strings"
)

var (
	domain string
	port   string
	m      sync.RWMutex
	urls   = make(map[string]string)
)

func main() {
	readArgs()
	readURLs()
	go writeURLs()
	startServer()
}

func readArgs() {
	flag.StringVar(&domain, "domain", "localhost", "")
	flag.StringVar(&port, "port", "8080", "")
	flag.Parse()
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

func startServer() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
	}
}

func handler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		postMethod(resp, req)
	case "GET":
		getMethod(resp, req)
	default:
		http.Error(resp, "Method is not supported.", http.StatusInternalServerError)
	}
}

func postMethod(resp http.ResponseWriter, req *http.Request) {
	typedURL := req.FormValue("url")
	longURL := extendURL(typedURL)
	shortURL := "http://" + domain + ":" + port + "/" + hashID(8)
	m.Lock()
	urls[shortURL] = longURL
	m.Unlock()
	fmt.Fprintln(resp, shortURL)
}

func extendURL(url string) string {
	prefixes := []string{"http://", "https://", "ftp://"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(url, prefix) {
			return url
		}
	}
	return "http://" + url
}

func getMethod(resp http.ResponseWriter, req *http.Request) {
	shortURL := "http://" + domain + ":" + port + req.URL.Path
	m.RLock()
	longURL, ok := urls[shortURL]
	m.RUnlock()
	if !ok {
		http.Error(resp, "Page not found.", http.StatusNotFound)
		return
	}
	http.Redirect(resp, req, longURL, http.StatusFound)
}

func hashID(len int) string {
	b := make([]byte, len/2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
