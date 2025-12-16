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

	for _, targetUrl := range targetUrls {
		targetUrl = strings.TrimSpace(targetUrl)
		if targetUrl == "" {
			continue
		}
		if !strings.HasPrefix(targetUrl, "http://") && !strings.HasPrefix(targetUrl, "https://") {
			targetUrl = "https://" + targetUrl
		}
		fmt.Printf("Scanning : %s\n", targetUrl)
		scanSite(targetUrl)
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

func scanSite(url string) {
	browserPath, found := launcher.LookPath()
	if !found {
		fmt.Println("Error while founding browser path")
		return
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cleanedURL := cleanup(url)
	outputDir := filepath.Join("output", cleanedURL)
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	linkfinder := `Array.from(document.querySelectorAll('a')).map(a => a.href);`

	var links []string
	var htmlContent string
	var imageByte []byte

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.FullScreenshot(&imageByte, 90),
		chromedp.Evaluate(linkfinder, &links),
	)
	if err != nil {
		fmt.Println("Error while opening chromedp")
		return
	}
	htmlPath := filepath.Join(outputDir, "source.html")
	screenshotPath := filepath.Join(outputDir, "screenshot.png")
	linksPath := filepath.Join(outputDir, "links.txt")

	err = os.WriteFile(screenshotPath, imageByte, 0644)
	if err != nil {
		fmt.Printf("Error while writing screenshot: %s\n", url)
	}
	err = os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error while writing html: %s\n", url)
	}

	err = os.WriteFile(linksPath, []byte(strings.Join(links, "\n")), 0644)
	if err != nil {
		fmt.Printf("Error while writing links: %s\n", url)
	}

	fmt.Printf("HTML content saved as %s\n", htmlPath)
	fmt.Printf("Screenshot saved as %s\n", screenshotPath)
	fmt.Printf("Links saved as %s\n", linksPath)
}
