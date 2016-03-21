package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var serviceDir string
var pagePath string
var currentStream string

const port = ":8041"

type streamInfo struct {
	Current string
	Host    string
	Playing bool
}

func handler(w http.ResponseWriter, r *http.Request) {
	currentStream = r.FormValue("url")

	t, _ := template.ParseFiles(pagePath)
	formInfo := streamInfo{Current: currentStream, Host: strings.TrimSuffix(r.Host, port), Playing: currentStream != ""}
	t.Execute(w, formInfo)
}

func main() {
	serviceDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	pagePath = path.Join(serviceDir, "mainpage.tmpl")
	// assuming it's running from cwd if go run is used
	if _, err := os.Stat(pagePath); os.IsNotExist(err) {
		serviceDir = "."
		pagePath = path.Join(serviceDir, "mainpage.tmpl")
	}

	http.HandleFunc("/", handler)
	http.Handle("/static/", http.FileServer(http.Dir(path.Join(serviceDir, "static"))))
	http.ListenAndServe(port, nil)
}
