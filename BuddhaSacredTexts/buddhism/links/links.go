package links

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Links struct {
	Link string
}

func GospelLinks() {

	allLinks := make([]Links, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("archive.sacred-texts.com"),
	)

	c.OnHTML("a", func(e *colly.HTMLElement) {
		links := "https://archive.sacred-texts.com/bud/btg/" + e.Attr("href")
		fmt.Println(links)
		link := Links{
			Link: links,
		}

		allLinks = append(allLinks, link)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://archive.sacred-texts.com/bud/btg/index.htm")
	if err != nil {
		log.Println("Unable to visit site", err)
		return
	}

	writeJsonLinks(allLinks)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		t := e.Text
		f, err := os.OpenFile("GospelLinks.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Unable to open file", err)
		}
		defer f.Close()

		if _, err := f.WriteString(t + "\n"); err != nil {
			log.Println("Unable to write to file", err)
		}

	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for i := 0; i < len(L); i++ {
		err = c.Visit(L[i])
		if err != nil {
			log.Println("Unable to visit site", err)
			return
		}

	}

	CleanText()

}

// Writes links to json for storage
func writeJsonLinks(data []Links) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file", err)
		return
	}

	_ = os.WriteFile("GospelLinks.json", file, 0644)
}

// Removes those pesky headers
func CleanText() {
	input, err := os.ReadFile("GospelLinks.txt")
	if err != nil {
		log.Println("Unable to read file", err)
	}
	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, S[0]) || strings.Contains(line, S[1]) || strings.Contains(line, S[2]) || strings.Contains(line, S[3]) || strings.Contains(line, S[4]) {
			lines[i] = ""
			lines[i] = strings.TrimSpace(lines[i])
			lines[i] = strings.TrimRight(lines[i], "\n")

		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile("GospelLinks.txt", []byte(output), 0644)
	if err != nil {
		log.Println("Unable to write to file", err)
	}
}
