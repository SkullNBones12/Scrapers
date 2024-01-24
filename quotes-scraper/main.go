package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// This is the struct where data goes
type QuoteScrape struct {
	Quote  string `json:"text"`
	Author string `json:"author"`
	Tags   string `json:"tags"`
}

func main() {
	// Makes an slice from struct to put data
	allQuotes := make([]QuoteScrape, 0)

	// New collector to grab websites
	collector := colly.NewCollector(
		colly.AllowedDomains("quotes.toscrape.com", "www.quotes.toscrape.com"),
	)

	// The HTML that is being parsed from "quote" class; the ChildText returns concatenated text (was a real problem for qsTags)
	collector.OnHTML(".quote", func(element *colly.HTMLElement) {
		qsDesc := element.ChildText(".text")
		qsAuth := element.ChildText(".author")
		qsTags := element.ChildText(".tags")

		// Create new QuoteScrape struct
		qS := QuoteScrape{
			Quote:  qsDesc,
			Author: qsAuth,
			Tags:   cleanTags(qsTags),
		}
		// Appends each new instance of qS to allQuotes slice
		allQuotes = append(allQuotes, qS)
	})

	// Simple, shows scraper is working and what page it is on
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	// Simple for-loop to iterate through websites
	for i := 1; i < 11; i++ {
		a := strconv.Itoa(i)
		b := fmt.Sprintf("http://quotes.toscrape.com/page/%s/", a)
		err := collector.Visit(b)

		if err != nil {
			log.Printf("Unable to visit: %s", b)
		}
	}

	writeJson(allQuotes)
}

// This function was copied from a tutorial, not entirely sure how it works
func writeJson(data []QuoteScrape) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = os.WriteFile("quotes.json", file, 0644)

}

// Cleans the tags, which came with all sorts of whitespace and line break drama
func cleanTags(x string) string {
	x = strings.Replace(x, "Tags:\n", "", -1)
	x = strings.Replace(x, " ", "", -1)
	x = strings.Trim(x, "\n")
	x = strings.Replace(x, "\n", ",", -1)
	x = strings.Replace(x, ",,", ", ", -1)

	return x
}
