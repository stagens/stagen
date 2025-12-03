package stagen

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
)

var (
	ErrGeneratorAlreadyExists  = errors.New("generator already exists")
	ErrGeneratorSourceNotFound = errors.New("generator source not found")
	ErrGeneratorUnknownSource  = errors.New("unknown generator source")
	ErrGeneratorBuildFailed    = errors.New("generator build failed")
)

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

	generatorSourceConfig := generatorConfig.Source()

	generatorSourceType := generatorSourceConfig.Type()
	generatorSourceName := generatorSourceConfig.Name()

	var generatorSource GeneratorSource

	switch generatorSourceType {
	case GeneratorSourceTypeAggDict:
		aggDictConfig, ok := s.aggDicts[generatorSourceName]
		if !ok {
			return fmt.Errorf("%w: generator '%s'", ErrGeneratorSourceNotFound, generatorName)
		}

		aggDictData, ok := s.aggDictsData[generatorSourceName]
		if !ok {
			return fmt.Errorf("%w: generator data '%s'", ErrGeneratorSourceNotFound, generatorName)
		}

		generatorSource = NewAggDictGeneratorSource(
			aggDictConfig,
			aggDictData,
		)

	case GeneratorSourceTypeDatabase:
		database, ok := s.databases[generatorSourceName]
		if !ok {
			return fmt.Errorf("%w: database '%s'", ErrGeneratorSourceNotFound, generatorName)
		}

		generatorSource = NewDatabaseGeneratorSource(
			database,
		)

	case GeneratorSourceTypeData:
		generatorSource = NewDataGeneratorSource(
			generatorConfig.Data(),
		)

	default:
		return fmt.Errorf("%w: %s", ErrGeneratorUnknownSource, generatorSourceType)
	}

	templateDirs := make([]string, 0)

	templateDirs = append(templateDirs, filepath.Join(s.templatesDir(), "templates"))

	for _, theme := range s.themes {
		templateDirs = append(templateDirs, filepath.Join(theme.Path(), "templates"))
	}

	for _, extension := range s.extensions {
		templateDirs = append(templateDirs, filepath.Join(extension.Path(), "templates"))
	}

	s.generators[generatorName] = NewGenerator(
		generatorConfig,
		generatorSource,
		s.clock,
		s.storage,
		templateDirs,
		s.createPage,
	)

	return nil
}

func (s *Impl) buildGenerators(ctx context.Context) error {
	s.log.GetLogger(ctx).Info("Building generators...")

	for _, generator := range s.generators {
		if err := s.buildGenerator(ctx, generator); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrGeneratorBuildFailed, generator.Config().Name(), err)
		}
	}

	return nil
}

func (s *Impl) buildGenerator(ctx context.Context, generator Generator) error {
	generatorName := generator.Config().Name()

	s.log.GetLogger(ctx).Infof("Building generator '%s'...", generatorName)

	pages, err := generator.Generate(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s: %w", ErrGeneratorBuildFailed, generatorName, err)
	}

	for _, page := range pages {
		s.pages[page.Id()] = page
	}

	return nil
}
