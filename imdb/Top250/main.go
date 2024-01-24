package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Info struct {
	Title  string
	Year   string
	Length string
	Rating string
}

func main() {

	allInfo := make([]Info, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com"),
	)

	c.OnHTML("div.sc-1e00898e-0", func(e *colly.HTMLElement) {

		infos := make([]string, 0)

		title := e.ChildText("a h3")

		e.ForEach("div span", func(_ int, e *colly.HTMLElement) {
			info := e.Text
			infos = append(infos, info)
		})

		t := Info{
			Title:  title,
			Year:   infos[0],
			Length: infos[1],
			Rating: infos[2],
		}

		allInfo = append(allInfo, t)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://www.imdb.com/chart/top/")
	if err != nil {
		log.Println("Unable to visit site", err)
		return
	}

	writeJsonLinks(allInfo)

}

func writeJsonLinks(data []Info) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file", err)
		return
	}

	_ = os.WriteFile("IMDBtop250.json", file, 0644)
}
