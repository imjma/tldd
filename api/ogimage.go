package api

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func HandleOGImage(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("u")
	if url == "" {
		fmt.Fprintf(w, "u?")
		return
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			// Do something with n...
			matched := false
			var result string
			for _, attr := range n.Attr {
				if attr.Key == "property" && attr.Val == "og:image" {
					matched = true
				}
				if attr.Key == "content" {
					result = attr.Val
				}
			}
			if matched {
				fmt.Fprintf(w, result)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

}
