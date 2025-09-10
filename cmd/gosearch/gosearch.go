package main

import (
	"flag"
	"fmt"
	"homeapp/pkg/crawler"
	"homeapp/pkg/crawler/spider"
	"homeapp/pkg/index"
	"homeapp/pkg/saver"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

var urls = [2]string{
	"https://go.dev",
	"https://golang.org",
}

func main() {
	go run()
	http.ListenAndServe(":8080", nil)
}

func run() {
	sFlag := flag.String("s", "", "строковое значение для обрабработки")

	flag.Parse()

	fmt.Println(*sFlag)

	if *sFlag != "" {

		saverService := saver.Saver{
			FilePath: fmt.Sprintf("%s.json", *sFlag),
		}
		cache := saverService.Read()
		lsFlag := strings.ToLower(*sFlag)

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
					return
				}

				iIndex = index.BuildInvertedIndex(allDocs)
				shouldSave = true
			}

			docIds = iIndex[lsFlag]
			docs := make([]crawler.Document, 0)
			for _, id := range docIds {
				doc := index.FindDocument(allDocs, id)

				if doc.ID != 0 {
					docs = append(docs, *doc)
				}
			}

			for _, doc := range docs {
				fmt.Println(doc.URL, doc.Title)
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
	} else {
		log.Fatal("Ошибка: нужно указать значение флаг -s")
	}

	http.ListenAndServe(":8080", nil)
}
