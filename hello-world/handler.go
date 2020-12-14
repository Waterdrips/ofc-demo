package function

import (
	"fmt"
	"github.com/openfaas/openfaas-cloud/sdk"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	token, err := sdk.ReadSecret("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	if token != query.Get("token") {
		http.Error(w, fmt.Sprintf("Token: %s, invalid", query.Get("token")), http.StatusUnauthorized)
		return
	}

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
		case "/hello-bot":
			if len(text) == 0 {
				w.Write([]byte("Please give a function name with this slash command"))
				w.WriteHeader(http.StatusOK)
				return true
			}



			w.WriteHeader(http.StatusOK)
			w.Write([]byte(text))
			return true
		}

	return false
}

