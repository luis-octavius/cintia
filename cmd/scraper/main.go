package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Not enough arguments provided")
	}

	pageURL := args[1]
	fmt.Println("Page URL: ", pageURL)

	html, err := getHTML(pageURL)
	if err != nil {
		log.Fatal("getHTML error: ", err)
	}

	fmt.Println("HTML: ", html)
}
