package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var serviceDir string
var templateFile *template.Template
var currentStream string

const port = ":8041"

// those 2 files are in pwd
const streamFileName = "input_video"
const pidFile = "vlc.pid"

type streamInfo struct {
	Current string
	Host    string
	Playing bool
}

func refreshStream() {
	// Write new stream file
	ioutil.WriteFile(streamFileName, []byte("STREAM_URL="+currentStream), 0644)

	// Kill whole process group (hence no use of signal API) referenced in vlc.pid
	content, err := ioutil.ReadFile(pidFile)
	if err == nil {
		gpid := strings.TrimSuffix(string(content), "\n")
		cmd := exec.Command("kill", "-9", "--", "-"+gpid)
		err := cmd.Run()
		if err != nil {
			log.Println("Couldn't kill "+gpid+". Error:", err)
		}
	}
}

/* Serve current handler form and refresh the stream file */
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		currentStream = r.FormValue("url")
		go refreshStream()
	}

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
