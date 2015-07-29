package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	pageTop = `<!DOCTYPE HTML>
    <html>
    <head>
<style>
.error{color:#FF0000;}
</style></head>
<title>Notices</title>
<body><h3>Notices</h3>
<p>You may create notice or delete all</p>`
	form       = `<form action="/" method="POST"><label for="Notices">Input text of notice:</label><br /><input type="textarea" name="Notice" ><br /><input type="submit" name="sendButton" value="send"><input type="submit" name="deleteButton" value="Delete All"></form>`
	pageBottom = `</body></html>`
	anError    = `<p class="error">%s</p>`
)

var mu = &sync.Mutex{}
var text = " "
var E []string

func main() {
	E = make([]string, 0)
	http.HandleFunc("/", homePage)
	handlerCMDArgs()
}

func handlerCMDArgs() {
	port := flag.String("port", ":9111", "port in server")
	flag.Parse()
	if err4 := http.ListenAndServeTLS(*port, "cert.pem", "key.pem", nil); err3 != nil {
		log.Fatal("failed to start server", err4)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		log.Println("problem with reflection of page", anError)
	} else {
		if message, ok := processRequest(request); ok && message != nil {
			formatStats(message)
			fmt.Fprint(writer, text)
		} else if ok == true && message == nil {
			fmt.Fprintf(writer, anError, "clear string in input form, 204 No Content")
		}
	}
	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]string, bool) {
	s := request.FormValue("sendButton")
	if s == "send" {
		if slice, found := request.Form["Notice"]; found && len(slice) > 0 {
			s := ""
			for i := 0; i < len(slice); i++ {
				s += slice[i]
			}
			s = html.EscapeString(s)
			l := strings.Fields(s)
			if len(l) > 0 {
				mu.Lock()
				E = append(E, s)
				mu.Unlock()
				return E, true
			} else {
				return nil, true
			}
		}
	}
	d := request.FormValue("deleteButton")
	if d != "" {
		mu.Lock()
		text = " "
		E = make([]string, 0)
		mu.Unlock()
		return nil, false
	}
	return nil, false
}

func formatStats(stats []string) {
	s := " "
	for i := 0; i < len(stats); i++ {
		s += `<textarea>` + stats[i] + `</textarea>` + " "
	}
	mu.Lock()
	text = " "
	text = text + " " + s + " "
	mu.Unlock()
}
