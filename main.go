package main

import {
	"fmt"
	"net/http"
	"context"
	"time"
	"os"
	"log"

	"github.com/chromedp/chromedp"
}

func main() {
	siteURL := "https://example.com"

	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	fmt.Printf("Visiting site: %s\n", siteURL)

	var htmlContent string
	var imageByte []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(siteURL),
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.FullScreenshot(&imageByte, 90),

	)

	if err != nil {
		log.Fatal(err)
	}

	// Save screenshot to file
	err = os.WriteFile("screenshot.png", imageByte, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Print HTML content length
	fmt.Printf("HTML content length: %d\n", len(htmlContent))
	fmt.Println("Screenshot saved as screenshot.png")
}