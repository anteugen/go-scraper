package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type BookProduct struct {
	url, image, name, price string
}

func scratch(c *colly.Collector, bookProducts *[]BookProduct, stopSignal *bool) {
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 404 {
			fmt.Println("Page not found, stopping: ", r.Request.URL)
			*stopSignal = true
		} else {
			fmt.Println("Something went wrong: ", err)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
		bookProduct := BookProduct{
			url:   e.ChildAttr(".image_container a", "href"),
			image: e.ChildAttr(".image_container img", "src"),
			name:  e.ChildAttr("h3 a", "title"),
			price: e.ChildText("p.price_color"),
		}

		*bookProducts = append(*bookProducts, bookProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
}

func storeCSV(bookProducts []BookProduct) {
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}

	writer.Write(headers)

	for _, bookProduct := range bookProducts {
		record := []string{
			bookProduct.url,
			bookProduct.image,
			bookProduct.name,
			bookProduct.price,
		}

		writer.Write(record)
	}

	defer writer.Flush()
}

func main() {
	fmt.Println("Hello world!")

	var bookProducts []BookProduct
	stopSignal := false

	c := colly.NewCollector()

	scratch(c, &bookProducts, &stopSignal)

	i := 1
	for !stopSignal {
		url := fmt.Sprintf("https://books.toscrape.com/catalogue/page-%d.html", i)
		err := c.Visit(url)
		if err != nil {
			log.Printf("Error visiting URL %s: %v", url, err)
			break
		}
		if stopSignal {
			break
		}
		i++
	}

	storeCSV(bookProducts)
}
