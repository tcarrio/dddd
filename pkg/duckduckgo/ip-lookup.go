package duckduckgo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	answerAPI     = "https://api.duckduckgo.com/?q=what%27s+my+ip&format=json"
	exampleAnswer = "Your IP address is 35.29.211.13 in <a href=\"https://duckduckgo.com/\">Detroit, Michigan, United States (48243)</a>"
)

var re = regexp.MustCompile(`^Your IP address is (?P<IP>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) in .+$`)

// Answer is the JSON response for the DuckDuckGo Answer API
type Answer struct {
	AbstractSource   string        `json:"AbstractSource"`
	DefinitionURL    string        `json:"DefinitionURL"`
	ImageIsLogo      string        `json:"ImageIsLogo"`
	RelatedTopics    []interface{} `json:"RelatedTopics"`
	Meta             interface{}   `json:"meta"`
	DefinitionSource string        `json:"DefinitionSource"`
	Redirect         string        `json:"Redirect"`
	Infobox          string        `json:"Infobox"`
	Abstract         string        `json:"Abstract"`
	AbstractURL      string        `json:"AbstractURL"`
	Type             string        `json:"Type"`
	ImageHeight      string        `json:"ImageHeight"`
	ImageWidth       string        `json:"ImageWidth"`
	Image            string        `json:"Image"`
	AbstractText     string        `json:"AbstractText"`
	AnswerType       string        `json:"AnswerType"`
	Answer           string        `json:"Answer"`
	Heading          string        `json:"Heading"`
	Definition       string        `json:"Definition"`
	Entity           string        `json:"Entity"`
	Results          []interface{} `json:"Results"`
}

// GetIPAddress returns the public IP address of this host
func GetIPAddress() string {
	response, err := http.Get(answerAPI)
	if err != nil {
		return ""
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	body := Answer{}
	json.Unmarshal(bytes, &body)

	ip := re.FindStringSubmatch(body.Answer)
	if len(ip) <= 1 {
		log.Fatal("No IP matches found for DuckDuckGo Answer")
	}

	return ip[1]
}
