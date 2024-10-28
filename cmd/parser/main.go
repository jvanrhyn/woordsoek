package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

var urls = []string{
	"http://gesellig.co.za/t/28745",
	"http://gesellig.co.za/t/28719",
	"http://gesellig.co.za/t/28720",
	"http://gesellig.co.za/t/28721",
	"http://gesellig.co.za/t/28722",
	"http://gesellig.co.za/t/28723",
	"http://gesellig.co.za/t/28724",
	"http://gesellig.co.za/t/28725",
	"http://gesellig.co.za/t/28726",
	"http://gesellig.co.za/t/28727",
	"http://gesellig.co.za/t/28728",
	"http://gesellig.co.za/t/28730",
	"http://gesellig.co.za/t/28731",
	"http://gesellig.co.za/t/28732",
	"http://gesellig.co.za/t/28733",
	"http://gesellig.co.za/t/28734",
	"http://gesellig.co.za/t/28735",
	"http://gesellig.co.za/t/28736",
	"http://gesellig.co.za/t/28737",
	"http://gesellig.co.za/t/28738",
	"http://gesellig.co.za/t/28739",
	"http://gesellig.co.za/t/28740",
	"http://gesellig.co.za/t/28741",
	"http://gesellig.co.za/t/28742",
	"http://gesellig.co.za/t/28743",
	"http://gesellig.co.za/t/28744",
}

func main() {
	var words []string

	for _, url := range urls {
		body, err := fetchHTML(url)
		if err != nil {
			fmt.Println("Error fetching HTML:", err)
			continue
		}

		doc, err := html.Parse(strings.NewReader(body))
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
			continue
		}

		extractWords(doc, &words)
		sort.Strings(words)
	}

	fmt.Println(words)
	file, err := os.OpenFile("af-za.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	for _, word := range words {
		if _, err := file.WriteString(word + "\n"); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func extractWords(n *html.Node, words *[]string) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "lys-met-kolomme" {
				text := extractText(n)
				*words = append(*words, strings.Fields(text)...)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractWords(c, words)
	}
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.Trim(n.Data, "\"")
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c) + " "
	}
	return strings.TrimSpace(text)
}
