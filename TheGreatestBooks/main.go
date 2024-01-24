package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Item struct {
	BookNum    int
	BookURL    string
	BookName   string
	AuthorURL  string
	AuthorName string
	Synopsis   string
}

func main() {

	count := 1

	allItems := make([]Item, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("thegreatestbooks.org"),
	)

	baseUrl := "https://thegreatestbooks.org"

	c.OnHTML("li.list-group-item", func(e *colly.HTMLElement) {
		text := make([]string, 0)

		bUrl := e.ChildAttr("a", "href")

		e.ForEach("a", func(_ int, e *colly.HTMLElement) {
			f2 := e.Text
			text = append(text, f2)

		})

		aUrl := e.ChildAttr("a[data-turbo-frame]", "href")

		synop := e.ChildText("p")

		t := Item{
			BookNum:    count,
			BookURL:    baseUrl + bUrl,
			BookName:   text[0],
			AuthorURL:  baseUrl + aUrl,
			AuthorName: text[1],
			Synopsis:   synop,
		}
		allItems = append(allItems, t)
		count++
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for i := 1; i < 274; i++ {
		link := fmt.Sprintf("https://thegreatestbooks.org/page/%d", i)
		err := c.Visit(link)
		if err != nil {
			log.Println("Unable to visit site", err)
			return
		}
	}

	writeJson(allItems)
}

func writeJson(data []Item) {
	file, err := json.MarshalIndent(data, " ", "")
	if err != nil {
		log.Println("Unable to marshal json", err)
		return
	}

	_ = os.WriteFile("books.json", file, 0644)
}
