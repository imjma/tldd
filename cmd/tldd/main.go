package main

import (
	"log"
	"net/http"

	"github.com/imjma/tldd/api"
)

func main() {
	http.HandleFunc("/ogimage", api.HandleOGImage)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
