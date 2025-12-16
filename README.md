# Web Scraper

Hey there, this is a simple but handy web scraping tool I built with Go. Basically, it visits the URLs you give it, takes a full-page screenshot, downloads the source HTML, and extracts all the links it finds into the `output` folder.

It runs a real browser instance (Chrome) in the background using `chromedp`, so it handles Javascript-loaded content pretty well too.

## How to use it?

You'll need Go and Chrome (or a Chromium-based browser) installed on your machine. Once you have those, just hop into the terminal and run it like this:

**For a single site:**
```bash
go run main.go -url "https://example.com"
```

**If you have a list of sites:**
Create a file like `sites.txt` and dump your links in there one by one (there's an example file in the repo you can check out), then run:
```bash
go run main.go -file sites.txt
```

When it runs, it'll create a folder for each site inside `output` and save everything there.

That's pretty much it. Cheers!
