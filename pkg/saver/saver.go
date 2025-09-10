package saver

import (
	"encoding/json"
	"fmt"
	"homeapp/pkg/crawler"
	"homeapp/pkg/index"
	"io"
	"log"
	"os"
)

type Saver struct {
	FilePath string
}

type SavedData struct {
	Url         string              `json:"url"`
	InvertedIdx index.InvertedIndex `json:"inverted_idx"`
	Docs        []crawler.Document  `json:"docs"`
}

func (s *Saver) Read() map[string]SavedData {
	file, err := os.OpenFile(s.FilePath, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var savedDataMap = make(map[string]SavedData)
	if err = json.NewDecoder(file).Decode(&savedDataMap); err != nil {
		if err == io.EOF {
			fmt.Println("Файл пустой — создаём пустую карту")
		} else {
			log.Fatal(err)
		}
	}

	return savedDataMap
}

func (s *Saver) Save(data map[string]SavedData) {
	file, err := os.OpenFile(s.FilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatal("1", err)
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(&data); err != nil {
		log.Fatal(err)
	}
}
