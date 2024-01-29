package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/gocolly/colly"
)

type Prep struct {
	Author    string
	PrepTime  string
	CookTime  string
	TotalTime string
	Yield     string
}

type Ingredients struct {
	Ingredients []string
}

type Instructions struct {
	Instructions []string
}

type Notes struct {
	Notes []string
}

type Recipe struct {
	Link         string
	Title        string
	Synopsis     string
	Prep         Prep
	Description  string
	Ingredients  Ingredients
	Instructions Instructions
	Notes        Notes
}

var siteMaps []string = []string{"https://sallysbakingaddiction.com/post-sitemap.xml", "https://sallysbakingaddiction.com/post-sitemap2.xml"}
var CurrentUrl string

func main() {
	knownUrls := make([]string, 0)

	recipe := make([]Recipe, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("sallysbakingaddiction.com"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*http.*",
		RandomDelay: 5 * time.Second,
	})

	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		f := e.Text
		if !slices.Contains(badSites, f) {
			knownUrls = append(knownUrls, f)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		CurrentUrl = r.URL.String()
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML("div[class=site-inner]", func(e *colly.HTMLElement) {
		url := CurrentUrl

		ingredients := make([]string, 0)
		instructions := make([]string, 0)
		notes := make([]string, 0)

		f := e.ChildText("h1[class=entry-title]")

		g := e.ChildText("p em")

		p1 := e.ChildText("li[class=author]")
		p2 := e.ChildText("li[class=prep-time]")
		p3 := e.ChildText("li[class=cook-time]")
		p4 := e.ChildText("li[class=total-time]")
		p5 := e.ChildText("li[class=yield]")

		h := e.ChildText("div[class=tasty-recipes-desciption-body]")

		e.ForEach("li[data-tr-ingredient-checkbox]", func(_ int, h *colly.HTMLElement) {
			f := h.Text
			ingredients = append(ingredients, f)

		})

		e.ForEach("div[class=tasty-recipes-instructions-body]", func(_ int, h *colly.HTMLElement) {
			h.ForEach("li", func(_ int, j *colly.HTMLElement) {
				f := j.Text
				instructions = append(instructions, f)
			})

		})

		e.ForEach("div[class=tasty-recipes-notes-body]", func(_ int, h *colly.HTMLElement) {
			h.ForEach("li", func(_ int, j *colly.HTMLElement) {
				f := j.Text
				notes = append(notes, f)
			})
		})

		t := Prep{
			Author:    p1,
			PrepTime:  p2,
			CookTime:  p3,
			TotalTime: p4,
			Yield:     p5,
		}

		v := Ingredients{
			Ingredients: ingredients,
		}

		w := Instructions{
			Instructions: instructions,
		}

		x := Notes{
			Notes: notes,
		}

		u := Recipe{
			Link:         url,
			Title:        f,
			Synopsis:     g,
			Prep:         t,
			Description:  h,
			Ingredients:  v,
			Instructions: w,
			Notes:        x,
		}

		recipe = append(recipe, u)
	})

	for _, v := range siteMaps {
		err := c.Visit(v)
		if err != nil {
			log.Println("Unable to visit site", err)
		}
	}

	for _, v := range knownUrls {
		err := c.Visit(v)
		if err != nil {
			log.Println("Unable to visit recipe site", err)
		}
	}

	writeJsonRecipes(recipe)
}

func writeJsonRecipes(data []Recipe) {
	file, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		log.Println("Unable to create json file", err)
		return
	}

	_ = os.WriteFile("SallyBakingAddiction.json", file, 0644)
}
