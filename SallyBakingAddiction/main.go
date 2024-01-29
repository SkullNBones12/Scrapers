package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

type Recipe struct {
	Title        string
	Synopsis     string
	Prep         Prep
	Description  string
	Ingredients  Ingredients
	Instructions Instructions
}

var siteMaps []string = []string{
	"https://sallysbakingaddiction.com/post-sitemap.xml",
	"https://sallysbakingaddiction.com/post-sitemap2.xml",
}

func main() {
	knownUrls := make([]string, 0)
	recipe := make([]Recipe, 0)
	ingredients := make([]string, 0)
	instructions := make([]string, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("sallysbakingaddiction.com"),
	)

	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		f := e.Text
		knownUrls = append(knownUrls, f)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "1 Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML("div[class=site-inner]", func(e *colly.HTMLElement) {
		f := e.ChildText("h1[class=entry-title]")

		g := e.ChildText("p em")

		p1 := e.ChildText("li[class=author]")
		p2 := e.ChildText("li[class=prep-time]")
		p3 := e.ChildText("li[class=cook-time]")
		p4 := e.ChildText("li[class=total-time]")
		p5 := e.ChildText("li[class=yield]")

		t := Prep{
			Author:    p1,
			PrepTime:  p2,
			CookTime:  p3,
			TotalTime: p4,
			Yield:     p5,
		}

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

		v := Ingredients{
			Ingredients: ingredients,
		}

		w := Instructions{
			Instructions: instructions,
		}

		u := Recipe{
			Title:        f,
			Synopsis:     g,
			Prep:         t,
			Description:  h,
			Ingredients:  v,
			Instructions: w,
		}

		recipe = append(recipe, u)
	})

	//for _, v := range siteMaps {
	//	err := c.Visit(v)
	//	if err != nil {
	//		log.Println("Unable to visit site", err)
	//	}
	//	return
	//}

	c.Visit("https://sallysbakingaddiction.com/double-crust-chicken-pot-pie/")

	writeJsonRecipes(recipe)

	c.Wait()
}

func writeJsonRecipes(data []Recipe) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file", err)
		return
	}

	_ = os.WriteFile("Recipes.json", file, 0644)
}
