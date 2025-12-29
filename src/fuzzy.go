package main

import (
	"os"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/charmbracelet/log"
)

var index bleve.Index

func (m *model) initFuzzy() {
	if _, err := os.Stat(bleveDirName); err == nil {
		index, err = bleve.Open(bleveDirName)
	} else if os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(bleveDirName, mapping)
		if err != nil {
			log.Errorf("error creating belve: %v", err)
			panic(err)
		}
		for key, ci := range m.keyToCredInfo {
			index.Index(key, ci)
		}
	} else {
		log.Errorf("error reading file: %v", err)
		panic(err)
	}
}

func (m *model) updateFuzzy(id string, ci credInfo) {
	index.Index(id, ci)
}
func (m *model) removeFuzzy(id string) {
	index.Delete(id)
}

func queryTopIDs(query string) []string {
	lowerString := strings.ToLower(query)
	prefix := bleve.NewPrefixQuery(lowerString)

	fuzzy := bleve.NewFuzzyQuery(query)
	// bleve caps at 2, not very well documented
	fuzzy.SetFuzziness(2)

	wildcard := bleve.NewQueryStringQuery("*" + query + "*")

	bq := bleve.NewDisjunctionQuery(prefix, fuzzy, wildcard)
	searchRequest := bleve.NewSearchRequest(bq)
	searchRequest.Size = 100
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Errorf("failed to query: %v", err)
		panic(err)
	}

	orderedIDs := make([]string, 0)
	for _, result := range searchResult.Hits {
		orderedIDs = append(orderedIDs, result.ID)
	}
	return orderedIDs
}
