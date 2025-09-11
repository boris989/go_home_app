package main

import (
	"fmt"
	"homeapp/pkg/crawler"
	"homeapp/pkg/crawler/spider"
	"homeapp/pkg/index"
	"homeapp/pkg/netsrv"
	"homeapp/pkg/saver"
	"log"
	"strings"
)

var urls = [2]string{
	"https://go.dev",
	"https://golang.org",
}

func main() {
	netsrv.StartServer("8080", search)
}

func search(word string) []string {
	saverService := saver.Saver{
		FilePath: fmt.Sprintf("%s.json", word),
	}
	cache := saverService.Read()
	lsWord := strings.ToLower(word)

	result := []string{}

	var iIndex index.InvertedIndex
	var docIds []int
	var allDocs []crawler.Document
	var shouldSave bool

	for _, url := range urls {
		savedData, ok := cache[url]
		if ok {
			iIndex = savedData.InvertedIdx
			allDocs = savedData.Docs
		} else {
			var err error
			spiderService := spider.New()
			allDocs, err = spiderService.Scan(url, 2)

			if err != nil {
				log.Fatal(err)
				return result
			}

			iIndex = index.BuildInvertedIndex(allDocs)
			shouldSave = true
		}

		docIds = iIndex[lsWord]
		docs := make([]crawler.Document, 0)
		for _, id := range docIds {
			doc := index.FindDocument(allDocs, id)

			if doc.ID != 0 {
				docs = append(docs, *doc)
			}
		}

		for _, doc := range docs {
			fmt.Println(doc.URL, doc.Title)

			result = append(result, fmt.Sprintf("%s %s", doc.URL, doc.Title))
		}

		if shouldSave {
			cache[url] = saver.SavedData{
				Url:         url,
				InvertedIdx: iIndex,
				Docs:        docs,
			}
		}
	}
	if shouldSave {
		saverService.Save(cache)
	}

	return result
}
