package fuzzy

import (
	"os"
	"sort"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/charmbracelet/log"
	"github.com/dismint/dispass/internal/state"
	"github.com/dismint/dispass/internal/uconst"
)

func InitFuzzy(sm *state.Model) {
	var err error

	if _, statErr := os.Stat(uconst.BleveDirName); statErr == nil {
		sm.Index, err = bleve.Open(uconst.BleveDirName)
		if err != nil {
			log.Fatalf("error opening bleve index: %v", err)
		}
	} else if os.IsNotExist(statErr) {
		mapping := bleve.NewIndexMapping()
		sm.Index, err = bleve.New(uconst.BleveDirName, mapping)
		if err != nil {
			log.Fatalf("error creating bleve index: %v", err)
		}

		for key, ci := range sm.KeyToCredInfo {
			if err := sm.Index.Index(key, ci); err != nil {
				log.Printf("failed to index %s: %v", key, err)
			}
		}
	} else {
		log.Fatalf("failed to stat bleve dir: %v", statErr)
	}
}

func UpdateFuzzy(sm *state.Model, id string, ci state.CredInfo) {
	sm.Index.Index(id, ci)
}
func RemoveFuzzy(sm *state.Model, id string) {
	sm.Index.Delete(id)
}

func QueryTopIDs(sm *state.Model, query string) []string {
	lowerString := strings.ToLower(query)
	var searchRequest *bleve.SearchRequest
	if query != "" {
		prefix := bleve.NewPrefixQuery(lowerString)

		fuzzy := bleve.NewFuzzyQuery(query)
		// bleve caps at 2, not very well documented
		fuzzy.SetFuzziness(2)

		wildcard := bleve.NewQueryStringQuery("*" + query + "*")

		searchRequest = bleve.NewSearchRequest(
			bleve.NewDisjunctionQuery(prefix, fuzzy, wildcard),
		)
	} else {
		searchRequest = bleve.NewSearchRequest(bleve.NewMatchAllQuery())
	}
	searchRequest.Size = 10000
	searchResult, err := sm.Index.Search(searchRequest)
	if query == "" {
		sort.Slice(searchResult.Hits, func(i, j int) bool {
			secondSourceLower := strings.ToLower(
				sm.KeyToCredInfo[searchResult.Hits[j].ID].Source,
			)
			firstSourceLower := strings.ToLower(
				sm.KeyToCredInfo[searchResult.Hits[i].ID].Source,
			)
			return firstSourceLower < secondSourceLower
		})
	}
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	orderedIDs := make([]string, 0)
	for _, result := range searchResult.Hits {
		orderedIDs = append(orderedIDs, result.ID)
	}
	return orderedIDs
}
