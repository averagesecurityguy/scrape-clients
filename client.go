package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	url        = "https://scrape.asg.cx/api/scrapes"
	authHeader = "Bearer scrape-YWb36X7FvLWVK93crQec9neaJji3M3Ud+50JAtdTNpo"
)

type ScrapeFile struct {
	Key      string   `json:"key"`
	Location string   `json:"location"`
	Size     int      `json:"size"`
	User     string   `json:"user"`
	Domain   string   `json:"domain"`
	Sha256   string   `json:"sha256"`
	Tags     []string `json:"tags"`
}

func (sf *ScrapeFile) HasTag(tag string) bool {
	for _, t := range sf.Tags {
		return tag == t
	}

	return false
}

// Download Scrape files
func getScrapeFiles() ([]ScrapeFile, error) {
	var files []ScrapeFile

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return files, err
	}

	// Each request to the API must include an authorization header with your
	// bearer token.
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return files, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return files, err
	}

	defer resp.Body.Close()

	err = json.Unmarshal(body, &files)
	if err != nil {
		return files, err
	}

	return files, nil
}

func main() {
	// Each time you call the API you will likely get scrape files you have
	// already processed so you need to keep track of the files you have seen
	// so you don't process them a second time. Each scrape file has a unique
	// key so you can keep a list of keys you have already seen and ignore any
	// files whose key is in the list.
	seen := make(map[string]struct{})

	for {
		files, err := getScrapeFiles()
		if err != nil {
			fmt.Printf("Could not get scrape files: %s\n", err)
			os.Exit(1)
		}

		// The response will be a JSON document that contains a list of
		// scrape files, each of which look like the following:
		//
		// {
		//    'key': 'Z9Ji/f/2Uy/mT38KOHX/mRp5W+9Nfl+4rSks8uqeCQE',
		//    'location': '',
		//    'size': 131,
		//    'user': 'yulei',
		//    'domain': 'github.com',
		//    'sha256': '34dcfdb4d756318cb04ba7ca9b6cb4872ae7a447a403df9c56eb00e93576be94',
		//    'tags': ['amazon-aws-secret', 'match-regex']
		// }
		//
		// You can use the metadata from the scrape file to determine which
		// interesting files you want to download.
		for _, file := range files {
			// Skip scrape files we've already seen
			if _, ok := seen[file.Key]; ok {
				continue
			}

			// Print the location of any scrape file that is tagged 'emails'
			if file.HasTag("emails") {
				fmt.Println(file.Location)
			}

			// Append the key to our seen list so we don't process this file
			// again.
			seen[file.Key] = struct{}{}
		}

		// There is no need to query the API more often than every 30 seconds.
		fmt.Println("Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)
	}
}
