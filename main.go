package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	domain string
	port   string
	m      sync.RWMutex
	urls   = make(map[string]string)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	f, err := os.Open("urls.json")
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
	signal.Notify(sig, os.Interrupt)
	<-sig
	f, err := os.Create("urls.json")
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
	protocols := [...]string{"http://", "https://", "ftp://"}
	for _, protocol := range protocols {
		if strings.HasPrefix(url, protocol) {
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

func hashID(length int) string {
	emojis := ""
	for i := 0; i < length; i++ {
		n := rand.Intn(len(emojiCodes))
		emojis += emojiCodes[n]
	}
	return fmt.Sprintf("%s", emojis)
}
