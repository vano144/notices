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

type Notice struct {
	Note  string
	Owner string
}

type notices struct {
	Store []Notice
	sync.Mutex
}

var StoreNotices notices
var t *template.Template

func main() {
	var err1 error
	StoreNotices.Store = make([]Notice, 0)
	file := path.Join("html", "disignFile.html")
	t, err1 = template.ParseFiles(file)
	if err1 != nil {
		log.Fatal("problem with parsing file", err1)
	}
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/message/", homePage)
	port := flag.String("port", ":9111", "port in server")
	flag.Parse()
	if err4 := http.ListenAndServeTLS(*port, "cert.pem", "key.pem", nil); err4 != nil {
		log.Fatal("failed to start server", err4)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	name, _, flg := request.BasicAuth()
	if flg {
		writer.Header().Set("Content-type", "text/html")
		err := request.ParseForm()
		if err != nil {
			log.Fatal("Problem with parsing form", err)
		}
		reqSend := request.PostFormValue("sendButton")
		if reqSend != "" {
			if slice, found := request.Form["Notice"]; found && len(slice) > 0 {
				s := ""
				s = strings.Join(slice, "")
				var k Notice
				k.Note = s
				k.Owner = name
				StoreNotices.Lock()
				StoreNotices.Store = append(StoreNotices.Store, k)
				StoreNotices.Unlock()
			}
		} else {
			reqDel := request.PostFormValue("deleteButton")
			if reqDel != "" {
				StoreNotices.Lock()
				StoreNotices.Store = make([]Notice, 0)
				StoreNotices.Unlock()
			}
		}
		t.Execute(writer, StoreNotices)
	}
	writer.Header().Set("WWW-Authenticate", `Basic realm="protectedpage"`)
	writer.WriteHeader(401)
}
