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

		//w.WriteHeader(http.StatusOK)
		//w.Write([]byte(text))

		n := emoji.ReplaceEmoticonsWithCodes(text)
		log.Printf("n: %s", n)
		w.Write([]byte(n))

		return true
	case "/invoke":
		fn := strings.Split(text, " ")
		domains := strings.Join(fn[1:], " ")
		r := strings.NewReader(domains)

		log.Printf("calling function: %s with [%s]", fn[0], domains)

		resp, err := http.Post(fmt.Sprintf("https://waterdrips.heyal.uk/%s", fn[0]), "application/x-www-form-urlencoded", r)
		if err != nil {
			log.Printf("Error calling gateway %v", err)
		}

		log.Printf("Response code: %d, content-length:%d", resp.StatusCode, resp.ContentLength)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Error reading body %v", err)
		}
		wrapped := fmt.Sprintf("you asked:%s\noutput:\n```%s```", text, string(body))
		w.Write([]byte(wrapped))
		return true
	}

	return false
}
