package main

import (
	"net/http"
	"log"
	"bytes"
	"net/url"
	//"fmt"
	"os"
	"fmt"
	"encoding/json"
	"strings"
	"io"
)

//Response structure serialisable to JSON
type Response struct {
	Text string `json:"text"`
}

type Synonymes struct {
	Word      string
	Synonymes []string
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	var data bytes.Buffer // []byte with IO

	_, err := data.ReadFrom(r.Body)
	if err != nil {
		return
	}

	q, err := url.ParseQuery(data.String())
	if err != nil {
		log.Fatal(err)
	}

	if q.Get("token") != os.Getenv("BOT_TOKEN") {
		log.Println("invalid token")
		return
	}

	text := q.Get("text")

	getSynonymes(strings.Split(text, " "))

	// Send personalised response
	json.NewEncoder(w).Encode(&Response{
		Text: fmt.Sprintf("Hello %s, how are you? %s", q.Get("user_name"), "Ja jestem SlackBotem"),
	})
}
func getSynonymes(words []string) {
	//GET http://workshop.x7467.com:1080/code
	// Fetch the URL
	for i, word := range words {
		getSynonyme(word)
	}

}
func getSynonyme(word string) string {
	resp, getErr := http.Get("http://workshop.x7467.com:1080/" + word)
	if getErr != nil {
		log.Fatal(getErr) // problem with connection or protocol
	}

	// Close the underlying connection when done
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("status: ", resp.Status) // unexpected status (not 200 OK)
	}


	var synonymes []Synonymes
	jsonErr := json.Unmarshal(resp.Body, &synonymes)
	if jsonErr != nil {
		fmt.Println("error:", jsonErr)
	}
	fmt.Printf("%+v", synonymes)

	// Print the response
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatal(err)
	}
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusNotFound {
		return ""
	}

	if r.StatusCode != http.StatusOK {
		log.Fatal("status: ", r.Status) // unexpected status (not 200 OK)
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	// Listen on address specified in $BOT_ADDR, or :8080 if empty
	addr := os.Getenv("BOT_ADDR")
	if addr == "" {
		addr = ":14605"
	}
	log.Println("Using adress:", addr)

	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}