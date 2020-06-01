package api

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type jsonResponse struct {
	Data string `json:"data,omitempty"`
}

func HandleOGImage(w http.ResponseWriter, r *http.Request) {
	response := processOGImage(r)

	w.Header().Set("Content-Type", "application/json")
	if response.Data == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

func processOGImage(r *http.Request) jsonResponse {
	var response jsonResponse
	if r.Method != http.MethodPost {
		return response
	}

	p := struct {
		URL string `json:"url"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err)
		return response
	}

	if p.URL == "" {
		log.Println("no url")
		return response
	}
	log.Println("url: ", p.URL)

	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		log.Println(err)
		return response
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return response
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			found := false
			var result string
			for _, attr := range n.Attr {
				if attr.Key == "property" && attr.Val == "og:image" {
					found = true
				}
				if attr.Key == "content" {
					result = attr.Val
				}
			}
			if found {
				log.Println(" og:image: ", result)
				response.Data = result
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return response
}
