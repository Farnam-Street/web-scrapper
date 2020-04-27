package main

import (
	"fmt"
	"io"
	"net/http"
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

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

func main() {

	var tickers []ticker
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

	//fmt.Println(tickers)
	//
	//f, err := os.Create("./tickers.csv")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer f.Close()
	//
	//w := csv.NewWriter(f)
	//for _, obj := range tickers {
	//	var record []string
	//	record = append(record, obj.Equity)
	//	record = append(record, obj.Symbol)
	//	record = append(record, obj.CrawledTime.String())
	//	w.Write(record)
	//}
	//w.Flush()

	for i2 := range tickers {
		//fmt.Println(tickers[i2].Symbol)
		fmt.Println("http://download.macrotrends.net/assets/php/stock_data_export.php?t=" + tickers[i2].Symbol)
		err := DownloadFile("http://download.macrotrends.net/assets/php/stock_data_export.php?t="+tickers[i2].Symbol, tickers[i2].Equity)
		if err != nil {
			panic(err)
		}
	}

}

func DownloadFile(url string, filepath string) error {
	var ret error
	Block{
		Try: func() {
			out, err := os.Create(filepath)
			if err != nil {
				ret = err
			}
			defer out.Close()

			// Get the data
			resp, err := http.Get(url)
			if err != nil {
				ret = err
			}
			defer resp.Body.Close()

			// Write the body to file
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				ret = err
			}

			ret = nil
		},
		Catch: func(e Exception) {
			fmt.Printf("Caught %v\n", e)
		},
		Finally: func() {

		},
	}.Do()
	return ret
	// Create the file

}
