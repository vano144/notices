package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
)

const templ = ` <!DOCTYPE HTML>
<html>
  <head>
   	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<style>

		</style>
	</head>
		<title>
			Notices
		</title>
	<body>
		<h3>Notices</h3>
		<p>You may create notice or delete all</p>
		<form action="/message" method="POST">
		<input type="textarea" name="Notice" >
		<input type="submit" name="sendButton" value="send">
		<input type="submit" name="deleteButton" value="delete">
		</form>
		<ol>
			{{ range .Slice }}
			<li>{{.}}</li>
		{{end}}
		</ol>
		<ol>
			{{ range .Errors }}
			<li>ERROR:{{.}}</li>
		{{end}}
		</ol>
	</body>
</html>
`

type notices struct {
	Slice  []string
	Errors []error
}

var mu = &sync.Mutex{}
var E notices

func main() {
	E.Slice = make([]string, 0)
	E.Errors = make([]error, 0)
	http.HandleFunc("/message", homePage)
	handlerCMDArgs()
}

func handlerCMDArgs() {
	port := flag.String("port", ":9111", "port in server")
	flag.Parse()
	if err4 := http.ListenAndServeTLS(*port, "cert.pem", "key.pem", nil); err4 != nil {
		log.Fatal("failed to start server", err4)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	var err2 error
	E.Errors = make([]error, 0)
	err := request.ParseForm()
	t := template.New("Person template")
	t, err1 := t.Parse(templ)
	a := request.FormValue("sendButton")
	if a != "" {
		if slice, found := request.Form["Notice"]; found && len(slice) > 0 {
			s := ""
			for i := 0; i < len(slice); i++ {
				s += slice[i]
			}
			l := strings.Fields(s)
			if len(l) > 0 {
				mu.Lock()
				E.Slice = append(E.Slice, s)
				mu.Unlock()
			} else {
				err2 = errors.New("clear string in input form, 204 No Content")
			}
		}
	}
	d := request.FormValue("deleteButton")
	if d != "" {
		mu.Lock()
		E.Slice = make([]string, 0)
		mu.Unlock()
	}
	switch true {
	case err != nil:
		mu.Lock()
		E.Errors = append(E.Errors, err)
		mu.Unlock()
		fallthrough
	case err1 != nil:
		mu.Lock()
		E.Errors = append(E.Errors, err1)
		mu.Unlock()
		fallthrough
	case err2 != nil:
		mu.Lock()
		E.Errors = append(E.Errors, err2)
		mu.Unlock()
	default:
	}
	t.Execute(writer, E)
}
