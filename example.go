package main

import (
	"fmt"
	"log"
	"net/http"
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
	form = `<form action="/" method="POST">
<label for="Notices">Input text of notice:</label><br />
<input type="textarea" name="Notice" ><br />
<input type="submit" value="send">
<input type="submit" value="Delete All">
</form>`
	pageBottom = `</body></html>`
	anError    = `<p class="error">%s</p>`
)

type notices struct {
	notice string
}

func main() {
	http.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":9007", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	} else {
		if message, ok := processRequest(request); ok {
			fmt.Fprint(writer, formatStats(message))
		} else if message != "" {
			fmt.Fprintf(writer, anError, message)
		}
	}
	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) (string, bool) {
	if slice, found := request.Form["Notice"]; found && len(slice) > 0 {
		var k notices
		for i := 0; i < len(slice); i++ {
			k.notice += slice[i]
		}
		return k.notice, true
	} else {
		return "", false
	}
	return "", false
}

func formatStats(stats string) string {
	return fmt.Sprintf(`<textarea>
%v</textarea>`, stats)
}
