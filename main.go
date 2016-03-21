package main

import (
	"html/template"
	"net/http"
	"strings"
)

var currentStream string

const port = ":8041"

type streamInfo struct {
	Current string
	Host    string
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("mainpage.tmpl")
	formInfo := streamInfo{Current: currentStream, Host: strings.TrimSuffix(r.Host, port)}
	t.Execute(w, formInfo)
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(port, nil)
}
