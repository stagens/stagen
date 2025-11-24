package stagen

import (
	"context"
	"errors"
	"fmt"
)

var ErrGeneratorAlreadyExists = errors.New("generator already exists")

func (s *Impl) loadGenerators(ctx context.Context) error {
	s.log.GetLogger(ctx).Info("Loading generators...")

	generators := make([]SiteGeneratorConfig, 0)

	generators = append(generators, s.siteConfig.Generators()...)

	for _, extension := range s.extensions {
		generators = append(generators, extension.Config().Generators()...)
	}

	for _, theme := range s.themes {
		generators = append(generators, theme.Config().Generators()...)
	}

	for _, generator := range generators {
		if err := s.loadGenerator(ctx, generator); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadGenerator, generator.Name(), err)
		}
	}

	return nil
}

func (s *Impl) loadGenerator(ctx context.Context, generatorConfig SiteGeneratorConfig) error {
	generatorName := generatorConfig.Name()
	if generatorName == "" {
		return ErrNoName
	}

	if _, ok := s.generators[generatorName]; ok {
		return fmt.Errorf("%w: %s", ErrGeneratorAlreadyExists, generatorName)
	}

	s.log.GetLogger(ctx).Infof("Loading generator '%s'...", generatorName)

	// @todo generators
	return nil
}
