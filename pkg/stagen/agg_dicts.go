package stagen

import (
	"context"
	"fmt"
)

func (s *Impl) loadAggDicts(ctx context.Context) error {
	s.log.GetLogger(ctx).Info("Loading agg dicts...")

	for _, aggDict := range s.siteConfig.AggDicts() {
		if err := s.loadAggDict(ctx, aggDict); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadAggDict, aggDict.Name(), err)
		}
	}

	return nil
}

func (s *Impl) loadAggDict(ctx context.Context, aggDictConfig SiteAggDictConfig) error {
	aggDictName := aggDictConfig.Name()
	if aggDictName == "" {
		return ErrNoName
	}

	s.log.GetLogger(ctx).Infof("Loading agg dict '%s'...", aggDictName)

	// @todo
	return nil
}
