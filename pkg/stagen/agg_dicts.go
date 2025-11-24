package stagen

import (
	"context"
	"errors"
	"fmt"
)

var ErrAggDictAlreadyExists = errors.New("agg dict already exists")

func (s *Impl) loadAggDicts(ctx context.Context) error {
	s.log.GetLogger(ctx).Info("Loading agg dicts...")

	aggDicts := make([]SiteAggDictConfig, 0)

	aggDicts = append(aggDicts, s.siteConfig.AggDicts()...)

	for _, extension := range s.extensions {
		aggDicts = append(aggDicts, extension.Config().AggDicts()...)
	}

	for _, theme := range s.themes {
		aggDicts = append(aggDicts, theme.Config().AggDicts()...)
	}

	for _, aggDict := range aggDicts {
		if err := s.loadAggDict(ctx, aggDict); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadAggDict, aggDict.Name(), err)
		}
	}

	return nil
}

func (s *Impl) loadAggDict(ctx context.Context, aggDictConfig SiteAggDictConfig) error {
	log := s.log.GetLogger(ctx)

	aggDictName := aggDictConfig.Name()
	if aggDictName == "" {
		return ErrNoName
	}

	if _, ok := s.aggDictsData[aggDictName]; ok {
		return fmt.Errorf("%w: %s", ErrAggDictAlreadyExists, aggDictName)
	}

	log.Infof("Loading agg dict '%s'...", aggDictName)

	s.aggDictsData[aggDictName] = make(map[string]map[string][]Page)

	aggDictKeys := aggDictConfig.Keys()

	for _, aggDictKey := range aggDictKeys {
		log.Infof("Collecting agg dict key '%s'...", aggDictKey)

		aggDictKeyData := make(map[string][]Page)

		for _, page := range s.pages {
			pageConfig := page.Config()
			pageVariables := pageConfig.Variables()

			pageVariable, ok := pageVariables[aggDictKey]
			if !ok {
				continue
			}

			switch variable := pageVariable.(type) {
			case string:
				aggDictKeyData[variable] = append(aggDictKeyData[variable], page)

			case []any:
				for _, arrayValueAny := range variable {
					switch arrayValue := arrayValueAny.(type) {
					case string:
						aggDictKeyData[arrayValue] = append(aggDictKeyData[arrayValue], page)

					default:
						// skip
					}
				}

			default:
				// skip
			}
		}

		s.aggDictsData[aggDictName][aggDictKey] = aggDictKeyData
	}

	return nil
}
