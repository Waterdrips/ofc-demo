package function

import (
	"fmt"
	"github.com/kenshaw/emoji"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	var query *url.Values
	if len(input) > 0 {
		q, err := url.ParseQuery(string(input))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		query = &q
	}

	//if token != query.Get("token") {
	//	http.Error(w, fmt.Sprintf("Token: %s, invalid", query.Get("token")), http.StatusUnauthorized)
	//	return
	//}
	var body []byte

	r.Body.Read(body)

	os.Stderr.Write(body)

	command := query.Get("command")
	text := query.Get("text")

	os.Stderr.Write([]byte(fmt.Sprintf("debug - command: %q, text: %q\n", command, text)))

	headerWritten := processCommand(w, command, text)

	if !headerWritten {
		http.Error(w, "Nothing to do", http.StatusBadRequest)
	}
}

func processCommand(w http.ResponseWriter, command, text string) bool {

	switch command {
	case "/echo":
		if len(text) == 0 {
			w.Write([]byte("Please give a message!"))
			w.WriteHeader(http.StatusOK)
			return true
		}

		w.WriteHeader(http.StatusOK)
		n := emoji.ReplaceEmoticonsWithCodes(text)
		log.Printf("n: %s", n)
		w.Write([]byte(n))
		return true
	case "/func":
		r := strings.NewReader(text)
		fn := strings.Split(text, " ")
		strings.Join(fn[1:], " ")
		resp, err := http.Post(fmt.Sprintf("http://gateway.openfaas:8080/function/%s", text), "text/plain", r)
		if err != nil {
			log.Printf("Error calling gateway %v", err)
		}

		var output []byte
		resp.Body.Read(output)
		w.Write(output)
	}

	return false
}
