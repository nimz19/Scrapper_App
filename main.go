package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// Review represents a single review structure
type Review struct {
	Title   string `json:"title"`
	Rating  string `json:"rating"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func main() {
	// Create a collector
	c := colly.NewCollector(
		colly.AllowedDomains("apps.apple.com"),
	)

	// Slice to hold reviews
	var reviews []Review

	// Scrape review data
	c.OnHTML(".we-customer-review", func(e *colly.HTMLElement) {
		review := Review{
			Title:   e.ChildText(".we-customer-review__title"),
			Rating:  e.ChildAttr(".we-star-rating", "aria-label"),
			Content: e.ChildText(".we-customer-review__body .we-truncate__text"),
			Author:  e.ChildText(".we-customer-review__user"),
		}
		reviews = append(reviews, review)
	})

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error:", err)
	})

	// Visit the app review page
	c.Visit("https://apps.apple.com/gb/app/lloyds-mobile-banking/id469964520?see-all=reviews")

	// Save data to JSON
	file, err := os.Create("reviews.json")
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(reviews); err != nil {
		log.Fatalf("Could not encode reviews: %v", err)
	}

	fmt.Println("Scraping complete. Data saved to reviews.json")
}
