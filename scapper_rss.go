package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Review represents a single review structure
type Review struct {
	Title   string `json:"title"`
	Rating  string `json:"rating"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// RSSFeed represents the structure of Apple’s RSS feed for app reviews
type RSSFeed struct {
	Feed struct {
		Entry []struct {
			Title struct {
				Label string `json:"label"`
			} `json:"title"`
			Content struct {
				Label string `json:"label"`
			} `json:"content"`
			Author struct {
				Name struct {
					Label string `json:"label"`
				} `json:"name"`
			} `json:"author"`
			ImRating struct {
				Label string `json:"label"`
			} `json:"im:rating"`
		} `json:"entry"`
	} `json:"feed"`
}

func main() {
	// URL for Apple’s RSS feed of recent reviews for the Lloyds Mobile Banking app (GB store)
	url := "https://itunes.apple.com/gb/rss/customerreviews/id=469964520/json"

	// Make the HTTP request to fetch the RSS feed
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch RSS feed: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the JSON response
	var rss RSSFeed
	if err := json.Unmarshal(body, &rss); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Convert RSS feed entries to our Review struct
	var reviews []Review
	for _, entry := range rss.Feed.Entry {
		review := Review{
			Title:   entry.Title.Label,
			Rating:  entry.ImRating.Label,
			Content: entry.Content.Label,
			Author:  entry.Author.Name.Label,
		}
		reviews = append(reviews, review)
	}

	// Save reviews to a JSON file
	file, err := os.Create("reviews_rss.json")
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(reviews); err != nil {
		log.Fatalf("Could not encode reviews: %v", err)
	}

	fmt.Println("RSS feed scraping complete. Data saved to reviews.json")
}
