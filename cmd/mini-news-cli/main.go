package main

import (
	"log"

	"github.com/nobe4/mini-news-cli/internal/ui"
)

func main() {
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
