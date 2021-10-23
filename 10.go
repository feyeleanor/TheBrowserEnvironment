package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const ADDRESS = ":3000"

type Commands map[string]func(http.ResponseWriter, *http.Request)
type PageConfiguration struct{ Commands }

func main() {
	p := PageConfiguration{
		Commands{
			"A": AJAX_handler("A"), "B": AJAX_handler("B"), "C": AJAX_handler("C")}}

	html, e := template.ParseFiles(BaseName() + ".html")
	Abort(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, p))
	})

	for c, f := range p.Commands {
		http.HandleFunc("/"+c, f)
	}
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

func BaseName() string {
	s := strings.Split(os.Args[0], "/")
	return s[len(s)-1]
}

func AJAX_handler(c string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, c)
	}
}
