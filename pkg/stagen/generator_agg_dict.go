package stagen

import (
	"context"
)

type AggDictGeneratorSource struct {
	aggDict     SiteAggDictConfig
	aggDictData map[string]map[string][]Page
}

func NewAggDictGeneratorSource(
	aggDict SiteAggDictConfig,
	aggDictData map[string]map[string][]Page,
) *AggDictGeneratorSource {
	return &AggDictGeneratorSource{
		aggDict:     aggDict,
		aggDictData: aggDictData,
	}
}

func (s *AggDictGeneratorSource) Entries(_ context.Context) ([]GeneratorSourceEntry, error) {
	entries := make([]GeneratorSourceEntry, 0)

	for aggDictKey, aggDictValues := range s.aggDictData {
		for aggDictValue, pages := range aggDictValues {
			pagesIds := make([]string, len(pages))
			for index, page := range pages {
				pagesIds[index] = page.Id()
			}

			entry := NewGeneratorSourceEntry(
				aggDictKey+"_"+aggDictValue,
				map[string]any{
					"AggDictKey":      aggDictKey,
					"AggDictValue":    aggDictValue,
					"AggDictPagesIds": pagesIds,
				},
			)

			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func (s *AggDictGeneratorSource) Variables() map[string]any {
	return map[string]any{
		"AggDict": map[string]any{
			"Name": s.aggDict.Name(),
			"Keys": s.aggDict.Keys(),
			"Data": s.aggDictData,
		},
	}
}
