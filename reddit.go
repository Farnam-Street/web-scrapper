package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"os"
)

type item struct {
	NumberOfLikes    string
	Title            string
	SubmittedAt      string
	SubmittedBy      string
	NumberOfComments string
	ImageUrl         string
	linkLabel        string
	Comments         []string
}

func main() {
	var source = "https://old.reddit.com/r/wallstreetbets/top/?sort=top&t=month"
	stories := []item{}
	c := colly.NewCollector(
		colly.AllowedDomains("old.reddit.com"),
	)

	//err := c.Post("http://old.reddit.com/login", map[string]string{"username": "seaborn07", "password": "Csk16209870"})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//c.OnResponse(func(r *colly.Response) {
	//	log.Println("response received", r.StatusCode)
	//})

	c.OnHTML(".thing", func(e *colly.HTMLElement) {
		c.Visit(e.ChildAttr("a.comments", "href") + "/?limit=500")
	})

	c.OnHTML(".comments-page", func(e *colly.HTMLElement) {
		temp := item{}
		temp.Title = e.ChildText("div.top-matter p.title a.title")
		temp.NumberOfLikes = e.ChildAttr("div.likes", "title")
		temp.SubmittedAt = e.ChildAttr("p.tagline time", "datetime")
		temp.SubmittedBy = e.ChildAttr("div.top-matter p.tagline a.author", "href")
		temp.NumberOfComments = e.ChildText("li.first a.comments")
		temp.ImageUrl = e.ChildAttr("img.preview", "src")
		temp.linkLabel = e.ChildAttr("span.linkflairlabel", "title")
		temp.Comments = e.ChildTexts("form div.usertext-body div.md p")
		stories = append(stories, temp)
	})

	c.OnHTML("span.next-button", func(h *colly.HTMLElement) {

		t := h.ChildAttr("a", "href")
		c.Visit(t)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Crawl all reddits the user passes in
	c.Visit(source)
	//fmt.Println(stories)
	jsonString, _ := json.MarshalIndent(stories, "", "\t")

	ioutil.WriteFile("redditall.json", jsonString, os.ModePerm)

}
