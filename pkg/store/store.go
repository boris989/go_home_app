package store

import (
	"homeapp/pkg/crawler"
	"homeapp/pkg/index"
)

type Store struct {
	Docs  []crawler.Document
	Index index.InvertedIndex
}

func New() *Store {
	return &Store{}
}
