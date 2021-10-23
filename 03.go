package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const LAUNCH_FAILED = 1
const FILE_READ = 2

const ADDRESS = ":3000"
const CERT = "cert.pem"
const KEY = "key.pem"

var MESSAGE string

func main() {
	s := strings.Split(os.Args[0], "/")
	html, e := ioutil.ReadFile(s[len(s)-1] + ".html")
	Abort(FILE_READ, e)
	MESSAGE = string(html)

	http.HandleFunc("/", billboard)
	Abort(LAUNCH_FAILED, http.ListenAndServeTLS(ADDRESS, CERT, KEY, nil))
}

func billboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, MESSAGE)
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}
