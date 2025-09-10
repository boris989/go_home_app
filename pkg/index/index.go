package index

import (
	"fmt"
	"homeapp/pkg/crawler"
	"sort"
	"strings"
)

type InvertedIndex map[string][]int

func BuildInvertedIndex(docs []crawler.Document) InvertedIndex {
	idx := make(InvertedIndex)
	for _, doc := range docs {
		words := strings.Fields(strings.ToLower(doc.Title))
		for _, w := range words {
			idx[w] = append(idx[w], doc.ID)
		}
	}

	for w := range idx {
		sort.Ints(idx[w])
	}

	return idx
}

func (idx InvertedIndex) Search(word string) []int {
	word = strings.ToLower(word)
	if ids, ok := idx[word]; ok {
		return ids
	}

	return []int{}
}

func FindDocument(docs []crawler.Document, id int) *crawler.Document {
	i := sort.Search(len(docs), func(i int) bool {
		return docs[i].ID >= id
	})
	if i < len(docs) && docs[i].ID == id {
		return &docs[i]
	}

	fmt.Println("NIL !!!!!")
	return nil
}
