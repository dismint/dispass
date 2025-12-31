package fuzzy

import (
	"os"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/charmbracelet/log"
	"github.com/dismint/dispass/internal/state"
	"github.com/dismint/dispass/internal/uconst"
)

func InitFuzzy(sm *state.Model) {
	if _, err := os.Stat(uconst.BleveDirName); err == nil {
		sm.Index, err = bleve.Open(uconst.BleveDirName)
	} else if os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		sm.Index, err = bleve.New(uconst.BleveDirName, mapping)
		if err != nil {
			log.Fatalf("error creating belve: %v", err)
		}
		for key, ci := range sm.KeyToCredInfo {
			sm.Index.Index(key, ci)
		}
	} else {
		log.Fatalf("error reading file: %v", err)
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
	prefix := bleve.NewPrefixQuery(lowerString)

	fuzzy := bleve.NewFuzzyQuery(query)
	// bleve caps at 2, not very well documented
	fuzzy.SetFuzziness(2)

	wildcard := bleve.NewQueryStringQuery("*" + query + "*")

	bq := bleve.NewDisjunctionQuery(prefix, fuzzy, wildcard)
	searchRequest := bleve.NewSearchRequest(bq)
	searchRequest.Size = 100
	searchResult, err := sm.Index.Search(searchRequest)
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	orderedIDs := make([]string, 0)
	for _, result := range searchResult.Hits {
		orderedIDs = append(orderedIDs, result.ID)
	}
	return orderedIDs
}
