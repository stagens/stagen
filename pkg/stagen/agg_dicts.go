package stagen

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrAggDictAlreadyExists = errors.New("agg dict already exists")
	ErrLoadAggDictData      = errors.New("load agg dict data")
)

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

	if _, ok := s.aggDicts[aggDictName]; ok {
		return fmt.Errorf("%w: %s", ErrAggDictAlreadyExists, aggDictName)
	}

	log.Infof("Loading agg dict '%s'...", aggDictName)

	s.aggDicts[aggDictName] = aggDictConfig

	return nil
}

func (s *Impl) loadAggDictsData(ctx context.Context) error {
	s.log.GetLogger(ctx).Info("Loading agg dicts data...")

	for _, aggDictConfig := range s.aggDicts {
		if err := s.loadAggDictData(ctx, aggDictConfig); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadAggDictData, aggDictConfig.Name(), err)
		}
	}

	return nil
}

func (s *Impl) loadAggDictData(ctx context.Context, aggDictConfig SiteAggDictConfig) error {
	log := s.log.GetLogger(ctx)

	aggDictName := aggDictConfig.Name()

	log.Infof("Loading agg dict data '%s'...", aggDictName)

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
				log.Trace("Added key", aggDictKey, "value", variable, "for page", page.Id())
				aggDictKeyData[variable] = append(aggDictKeyData[variable], page)

			case bool:
				boolStr := strconv.FormatBool(variable)
				log.Trace("Added key", aggDictKey, "value", boolStr, "for page", page.Id())
				aggDictKeyData[boolStr] = append(aggDictKeyData[boolStr], page)

			case []any:
				for _, arrayValueAny := range variable {
					switch arrayValue := arrayValueAny.(type) {
					case string:
						log.Trace("Added key", aggDictKey, "value", arrayValue, "for page", page.Id())
						aggDictKeyData[arrayValue] = append(aggDictKeyData[arrayValue], page)

					default:
						log.Warnf("Unsupported variable array value type: %T (%#v)", arrayValue, arrayValue)
						// skip
					}
				}

			default:
				log.Warnf("Unsupported variable type: %T (%#v)", variable, variable)
				// skip
			}
		}

		s.aggDictsData[aggDictName][aggDictKey] = aggDictKeyData
	}

	return nil
}
