package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var serviceDir string
var templateFile *template.Template
var currentStream string

const port = ":8041"

type streamInfo struct {
	Current string
	Host    string
	Playing bool
}

func handler(w http.ResponseWriter, r *http.Request) {
	currentStream = r.FormValue("url")

	formInfo := streamInfo{Current: currentStream, Host: strings.TrimSuffix(r.Host, port), Playing: currentStream != ""}
	templateFile.Execute(w, formInfo)
}

func main() {
	// Check for assets dir
	serviceDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	// Assuming it's running from cwd if go run is used
	if _, err := os.Stat(path.Join(serviceDir, "mainpage.tmpl")); os.IsNotExist(err) {
		serviceDir = "."
	}
	templateFile, err = template.ParseFiles(path.Join(serviceDir, "mainpage.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	// Load current stream if any
	content, err := ioutil.ReadFile(streamFileName)
	if err == nil {
		currentStream = strings.TrimPrefix(string(content), "STREAM_URL=")
	}

	// Starts server
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
}
