package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}

func hello(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello, 世界")
}
