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
	file := flag.String("file", "", "For file scan (example: sites.txt)")
	flag.Parse()

	if *siteURL == "" && *file == "" {
		fmt.Println("Error: Please provide -url or -file")
		return
	}
	var targetUrls []string
	if *siteURL != "" {
		targetUrls = append(targetUrls, *siteURL)
	}
	if *file != "" {
		fileContent, err := os.ReadFile(*file)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		targetUrls = strings.Split(string(fileContent), "\n")

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

	for _, targetUrl := range targetUrls {
		targetUrl = strings.TrimSpace(targetUrl)
		if targetUrl == "" {
			continue
		}
		if !strings.HasPrefix(targetUrl, "http://") && !strings.HasPrefix(targetUrl, "https://") {
			targetUrl = "https://" + targetUrl
		}
		fmt.Printf("Scanning : %s\n", targetUrl)
		scanSite(ctx, targetUrl)
	}
}

func cleanup(url string) string {
	cleaned := strings.ReplaceAll(url, "http://", "")
	cleaned = strings.ReplaceAll(cleaned, "https://", "")
	cleaned = strings.ReplaceAll(cleaned, "www.", "")
	cleaned = strings.Split(cleaned, "/")[0]
	cleaned = strings.Split(cleaned, "?")[0]
	return cleaned
}

func scanSite(ctx context.Context, url string) {
	cleanedURL := cleanup(url)
	outputDir := filepath.Join("output", cleanedURL)
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	var htmlContent string
	var imageByte []byte
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.FullScreenshot(&imageByte, 90),
	)
	if err != nil {
		log.Fatal(err)
	}
	htmlPath := filepath.Join(outputDir, "source.html")
	screenshotPath := filepath.Join(outputDir, "screenshot.png")

	err = os.WriteFile(screenshotPath, imageByte, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTML content saved as %s\n", htmlPath)
	fmt.Printf("Screenshot saved as %s\n", screenshotPath)
}
