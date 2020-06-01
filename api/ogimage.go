package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func HandleOGImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	p := struct {
		URL string `json:"url"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return
	}

	if p.URL == "" {
		return
	}

	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
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
				fmt.Fprintf(w, result)
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
