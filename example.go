package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

type Notice struct {
	Note  string
	Owner string
	Time  string
}

type notices struct {
	Store []Notice
	sync.Mutex
}

var StoreNotices notices
var templt *template.Template

func main() {
	var err1 error
	StoreNotices.Store = make([]Notice, 0)
	file := path.Join("html", "disignFile.html")
	if templt, err1 = template.ParseFiles(file); err1 != nil {
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
	name, _, successAuth := request.BasicAuth()
	stri := fmt.Sprint(request.Header.Get("Accept"))
	controlQuerry := strings.Contains(stri, "application/json`")
	if controlQuerry && successAuth {
		notesJson, erro := json.Marshal(StoreNotices)
		if erro != nil {
			log.Println("Error Json")
			return
		}
		writer.Header().Set("Content-type", "application/json")
		writer.Write(notesJson)
		return
	} else if controlQuerry && !successAuth {
		writer.Write([]byte("authorization is necessary"))
		return
	}
	if successAuth {
		writer.Header().Set("Content-type", "text/html")
		if err := request.ParseForm(); err != nil {
			log.Fatal("Problem with parsing form", err)
		}
		if reqSend := request.PostFormValue("sendButton"); reqSend != "" {
			if slice, found := request.Form["Notice"]; found && len(slice) > 0 {
				var k = Notice{
					Note:  strings.Join(slice, ""),
					Owner: name,
					Time:  fmt.Sprint(time.Now().Local().Format("15:04")),
				}
				StoreNotices.Lock()
				StoreNotices.Store = append(StoreNotices.Store, k)
				StoreNotices.Unlock()
			}
		} else {
			if reqDel := request.PostFormValue("deleteButton"); reqDel != "" {
				StoreNotices.Lock()
				StoreNotices.Store = make([]Notice, 0)
				StoreNotices.Unlock()
			}
		}
		templt.Execute(writer, StoreNotices)
	}
	writer.Header().Set("WWW-Authenticate", `Basic realm="protectedpage"`)
	writer.WriteHeader(401)
}
