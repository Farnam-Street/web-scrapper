package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type ticker struct {
	Equity      string
	Symbol      string
	CrawledTime time.Time
}

func main() {

	tickers := []ticker{}
	c := colly.NewCollector(
		colly.AllowedDomains("www.advfn.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/nasdaq/") {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML(".ts0", func(element *colly.HTMLElement) {
		//var name = element.ChildAttr("td","align")
		var child = element.ChildTexts("td")
		tick := ticker{}
		tick.CrawledTime = time.Now()
		for i := range child {
			if i == 0 {
				tick.Equity = child[i]
			} else if i == 1 {
				tick.Symbol = child[i]
			}
		}
		tickers = append(tickers, tick)
	})

	c.OnHTML(".ts1", func(element *colly.HTMLElement) {
		var child = element.ChildTexts("td")
		tick := ticker{}
		tick.CrawledTime = time.Now()
		for i := range child {
			if i == 0 {
				tick.Equity = child[i]
			} else if i == 1 {
				tick.Symbol = child[i]
			}
		}
		tickers = append(tickers, tick)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.advfn.com/nasdaq/nasdaq.asp")

	fmt.Println(tickers)

	f, err := os.Create("./tickers.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, obj := range tickers {
		var record []string
		record = append(record, obj.Equity)
		record = append(record, obj.Symbol)
		record = append(record, obj.CrawledTime.String())
		w.Write(record)
	}
	w.Flush()

}
