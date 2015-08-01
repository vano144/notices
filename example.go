package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
)

type notices struct {
	Slice []string
}

var mu sync.Mutex
var E notices
var t *template.Template

func main() {
	var err1 error
	E.Slice = make([]string, 0)
	file := path.Join("html", "disignFile.html")
	t, err1 = template.ParseFiles(file)
	if err1 != nil {
		log.Fatal("problem with parsing file")
	}
	http.HandleFunc("/message", homePage)
	port := flag.String("port", ":9111", "port in server")
	flag.Parse()
	if err4 := http.ListenAndServeTLS(*port, "cert.pem", "key.pem", nil); err4 != nil {
		log.Fatal("failed to start server", err4)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Fatal("Problem with parsing form")
	}
	a := request.FormValue("sendButton")
	if a != "" {
		if slice, found := request.Form["Notice"]; found && len(slice) > 0 {
			s := ""
			s = strings.Join(slice, "")
			mu.Lock()
			E.Slice = append(E.Slice, s)
			mu.Unlock()
		}
	} else {
		d := request.FormValue("deleteButton")
		if d != "" {
			mu.Lock()
			E.Slice = make([]string, 0)
			mu.Unlock()
		}
	}
	t.Execute(writer, E)
}
