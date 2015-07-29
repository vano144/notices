package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	if err3 := http.ListenAndServeTLS(*port, "cert.pem", "key.pem", nil); err3 != nil {
		log.Fatal("Failed to start server", err3)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, "problem with reflection of page, 500 Internal Server Error")
	} else {
		if message, ok := processRequest(request); ok {
			formatStats(message)
			fmt.Fprint(writer, text)
		}
	}
	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]string, bool) {
	s, d := " ", " "
	s = request.FormValue("sendButton")
	if s != " " {
		if slice, found := request.Form["Notice"]; found && len(slice) > 0 {
			s := ""
			for i := 0; i < len(slice); i++ {
				s += slice[i]
			}
			mu.Lock()
			E = append(E, s)
			mu.Unlock()
			return E, true
		} else {
			fmt.Fprintf(writer, anError, "clear string in input form, 204 No Content")
			return nil, false
		}
	}
	d = request.FormValue("deleteButton")
	if d != " " {
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
	mu.Lock()
	text = " "
	mu.Unlock()
	for i := 0; i < len(stats); i++ {
		s += `<textarea>` + stats[i] + `</textarea>` + " "
	}
	mu.Lock()
	text = text + " " + s + " "
	mu.Unlock()
}
