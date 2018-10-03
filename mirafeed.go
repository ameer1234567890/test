package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	origFeed  = "MIRAFeed.rss"
	fixedFeed = "FixedFeed.rss"
	newData   string
	newLine   string
)

func DoReplacements() {
	fmt.Println("Performing replacements...")

	// For replacing pubDate
	patternPubDate := "\\<pubDate>(\\d{2} \\w{3} \\d{4})\\</pubDate>"
	rePubDate, err := regexp.Compile(patternPubDate)
	if err != nil {
		panic(err)
	}

	// For replacing description
	patternDescription := "<description>Lates news &amp; updates of MIRA</description>"
	reDescription, err := regexp.Compile(patternDescription)
	if err != nil {
		panic(err)
	}

	// For replacing rss tag
	patternRssTag := "<rss version=\"2.0\">"
	reRssTag, err := regexp.Compile(patternRssTag)
	if err != nil {
		panic(err)
	}

	// For replacing space in URL
	patternSpaceInUrl := "<link>.* .*</link>"
	reSpaceInUrl, err := regexp.Compile(patternSpaceInUrl)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(fixedFeed)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newLine = scanner.Text()

		// Replace pubDate
		dateString := strings.TrimRight(strings.TrimLeft(rePubDate.FindString(scanner.Text()), "<pubDate>"), "</pubDate>")
		if dateString != "" {
			layout := "02 Jan 2006"
			t, err := time.Parse(layout, dateString)
			if err != nil {
				panic(err)
			}
			rssDate := strings.TrimRight(t.Format("Mon, 02 Jan 2006 15:04:05 -0700"), "+0000") + "+0500"
			newLine = "<pubDate>" + rssDate + "</pubDate>"
		}

		// Replace Descrription
		descriptionString := strings.TrimRight(strings.TrimLeft(reDescription.FindString(scanner.Text()), "<description>"), "</description>")
		if descriptionString != "" {
			newLine = "<description>Latest news &amp; updates of MIRA</description>"
		}

		// Replace RSS Tag
		rssTagString := reRssTag.FindString(scanner.Text())
		if rssTagString != "" {
			newLine = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"
			newLine += "<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n"
			newLine += "<channel>\n"
			newLine += "<atom:link href=\"https://ameer.io/test/FixedFeed.rss\" rel=\"self\" type=\"application/rss+xml\" />\n"
			newLine += "<title>MIRA RSS</title>"
		}

		// Replace space in URL
		spaceInUrlString := strings.TrimRight(strings.TrimLeft(reSpaceInUrl.FindString(scanner.Text()), "<link>"), "</link>")
		if spaceInUrlString != "" {
			newLine = "<link>" + strings.Replace(spaceInUrlString, " ", "%20", -1) + "</link>"
		}

		newData += newLine + "\n"
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	writer, err := os.Create(fixedFeed)
	if err != nil {
		panic(err)
	}
	defer writer.Close()
	writer.Write([]byte(newData))
}

func addLineBreaks() {
	fmt.Println("Adding line breaks to file...")
	pattern := "(</\\w*>)"
	re, err := regexp.Compile(pattern)
	repl := "${1}\n"
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile(origFeed)
	if err != nil {
		panic(err)
	}
	err = os.Remove(fixedFeed)
	if err != nil {
		panic(err)
	}
	writer, err := os.Create(fixedFeed)
	if err != nil {
		panic(err)
	}
	writer.Write(re.ReplaceAll(b, []byte(repl)))
}

func SaveFeedToFile() {
	fmt.Println("Saving feed to file...")
	url := "https://www.mira.gov.mv/MIRAFeed.aspx"
	response, e := http.Get(url)
	if e != nil {
		panic(e)
	}
	defer response.Body.Close()
	file, err := os.Create("MIRAFeed.rss")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func main() {
	SaveFeedToFile()
	addLineBreaks()
	DoReplacements()
}
