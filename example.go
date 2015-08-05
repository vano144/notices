package main

import (
	"encoding/base64"
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
	http.HandleFunc("/message/", homePage)
	port := flag.String("port", ":9111", "port in server")
	flag.Parse()
	if err4 := http.ListenAndServeTLS(*port, "cert.pem", "key.pem", nil); err4 != nil {
		log.Fatal("failed to start server", err4)
	}
}

func BasicAuth(w http.ResponseWriter, r *http.Request) (bool, string) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false, ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false, ""
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false, ""
	}
	if pair[0] == "invtest" {
		return false, ""
	}
	return true, pair[0]
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	flg, name := BasicAuth(writer, request)
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
	writer.Header().Set("WWW-Authenticate", `Basic realm="protectedpage, invalid name for test- invtest"`)
	writer.WriteHeader(401)
}
