package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Stories struct {
	Para []string
}

type SubSites struct {
	StoryTitle string
	Author     string
	Date       string
	Story      Stories
}

var TextFileTitle string

func main() {
	Date := time.Now().Format("2006-01-02")

	fmt.Println(Date)
	stringLinks := make([]string, 0)

	allStories := make([]Stories, 0)
	allSubSites := make([]SubSites, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("text.npr.org"),
	)

	baseUrl := "https://text.npr.org/"

	c.OnHTML("div[class=topic-container]", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, h *colly.HTMLElement) {
			l := baseUrl + h.Attr("href")
			stringLinks = append(stringLinks, l)
		})

	})

	c.OnHTML("article", func(e *colly.HTMLElement) {

		headers := make([]string, 0)
		paras := make([]string, 0)

		title := e.ChildText("h1[class=story-title]") + "\n"

		TextFileTitle = e.ChildText("h1[class=story-title]") + "_" + Date + ".txt"
		TextFileTitle = strings.Replace(TextFileTitle, ",", "", -1)
		TextFileTitle = strings.Replace(TextFileTitle, "'", "", -1)
		TextFileTitle = strings.Replace(TextFileTitle, "/", "", -1)
		TextFileTitle = strings.Replace(TextFileTitle, "?", "", -1)

		WriteFile(title)

		e.ForEach("div[class=story-head] p", func(_ int, h *colly.HTMLElement) {
			header := h.Text + "\n"
			headers = append(headers, header)
			WriteFile(header)
		})

		e.ForEach("div[class=paragraphs-container] p", func(_ int, h *colly.HTMLElement) {
			para := " " + h.Text + "\n"
			paras = append(paras, para)
			WriteFile(para)
		})

		u := Stories{
			Para: paras,
		}

		allStories = append(allStories, u)

		t := SubSites{
			StoryTitle: title,
			Author:     headers[0],
			Date:       headers[1],
			Story:      u,
		}

		allSubSites = append(allSubSites, t)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://text.npr.org")
	if err != nil {
		log.Printf("Unable to visit site, Err:%s", err)
		return
	}

	for i := 0; i < len(stringLinks); i++ {
		err := c.Visit(stringLinks[i])
		if err != nil {
			log.Printf("Unable to visit %s, Err: %s", stringLinks[i], err)
			return
		}
	}
}

func WriteFile(x string) {

	f, err := os.OpenFile(TextFileTitle, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Unable to write to file", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(x + "\n"); err != nil {
		log.Println("Unable to write string to file", err)
	}
}
