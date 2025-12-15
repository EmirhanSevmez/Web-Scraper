package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	siteURL := flag.String("url", "", "The URL of the site to scrape")
	flag.Parse()

	if *siteURL == "" {
		log.Fatal("Please provide a URL using the -url flag")
		return
	}
	targetUrl := *siteURL

	if !strings.HasPrefix(targetUrl, "http://") && !strings.HasPrefix(targetUrl, "https://") {
		targetUrl = "http://" + targetUrl
	}

	domainName := cleanup(targetUrl)

	outputDir := filepath.Join("results", domainName)

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	browserPath, found := launcher.LookPath()
	if !found {
		log.Fatal("Chrome/Chromium browser not found")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.ExecPath(browserPath),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Create timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	fmt.Printf("Visiting site: %s\n", targetUrl)

	var htmlContent string
	var imageByte []byte

	err = chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.FullScreenshot(&imageByte, 90),
	)

	if err != nil {
		log.Fatal(err)
	}
	htmlPath := filepath.Join(outputDir, "source.html")
	screenshotPath := filepath.Join(outputDir, "screenshot.png")

	// Save screenshot to file
	err = os.WriteFile(screenshotPath, imageByte, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//save HTML content to file
	err = os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTML content saved as ", htmlPath)
	fmt.Println("Screenshot saved as ", screenshotPath)
}

func cleanup(url string) string {
	cleaned := strings.ReplaceAll(url, "http://", "")
	cleaned = strings.ReplaceAll(cleaned, "https://", "")
	cleaned = strings.ReplaceAll(cleaned, "www.", "")
	cleaned = strings.Split(cleaned, "/")[0]
	cleaned = strings.Split(cleaned, "?")[0]
	return cleaned
}
