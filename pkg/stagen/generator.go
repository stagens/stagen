package stagen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/pixality-inc/golang-core/clock"
	"github.com/pixality-inc/golang-core/storage"
)

var (
	ErrNoSource                  = errors.New("no source")
	ErrFailedToGenEntries        = errors.New("failed to generate entries from source")
	ErrGeneratorTemplateNotFound = errors.New("generator template not found")
)

type CreatePageFunction = func(
	ctx context.Context,
	pageFileInfo *PageFileInfo,
	content []byte,
	extraVariables map[string]any,
	dirConfigs []PageConfig,
) (Page, error)

type GeneratorSourceEntry interface {
	Filename() string
	Variables() map[string]any
}

type GeneratorSourceEntryImpl struct {
	filename  string
	variables map[string]any
}

func NewGeneratorSourceEntry(
	filename string,
	variables map[string]any,
) *GeneratorSourceEntryImpl {
	return &GeneratorSourceEntryImpl{
		filename:  filename,
		variables: variables,
	}
}

func (e *GeneratorSourceEntryImpl) Filename() string {
	return e.filename
}

func (e *GeneratorSourceEntryImpl) Variables() map[string]any {
	return e.variables
}

type GeneratorSource interface {
	Entries(ctx context.Context) ([]GeneratorSourceEntry, error)
	Variables() map[string]any
}

type Generator interface {
	Config() SiteGeneratorConfig
	Generate(ctx context.Context) ([]Page, error)
}

type GeneratorImpl struct {
	config             SiteGeneratorConfig
	source             GeneratorSource
	clock              clock.Clock
	storage            storage.Storage
	templateDirs       []string
	createPageFunction CreatePageFunction
}

func NewGenerator(
	config SiteGeneratorConfig,
	source GeneratorSource,
	clocks clock.Clock,
	storage storage.Storage,
	templateDirs []string,
	createPageFunction CreatePageFunction,
) *GeneratorImpl {
	return &GeneratorImpl{
		config:             config,
		source:             source,
		clock:              clocks,
		storage:            storage,
		templateDirs:       templateDirs,
		createPageFunction: createPageFunction,
	}
}

func (g *GeneratorImpl) Config() SiteGeneratorConfig {
	return g.config
}

func (g *GeneratorImpl) Generate(ctx context.Context) ([]Page, error) {
	if g.source == nil {
		return nil, ErrNoSource
	}

	var err error

	templateExtensions := []string{
		".tmpl",
		".html.tmpl",
		".md.tmpl",
		".html",
		".md",
		"",
	}
	templateFilename := g.config.Template().Name()
	templatePath := ""
	found := false
	pageContent := make([]byte, 0)

	for _, templateDir := range g.templateDirs {
		if found {
			break
		}

		for _, templateExtension := range templateExtensions {
			templatePath = filepath.Join(templateDir, templateFilename+templateExtension)

			if exists, err := g.storage.FileExists(ctx, templatePath); err != nil {
				return nil, fmt.Errorf("faile to check if file %s exists: %w", templatePath, err)
			} else if !exists {
				continue
			}

			found = true

			pageFile, err := g.storage.ReadFile(ctx, templatePath)
			if err != nil {
				return nil, fmt.Errorf("%w: failed to read template file '%s': %w", ErrGeneratorTemplateNotFound, templatePath, err)
			}

			pageContent, err = io.ReadAll(pageFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read template file '%s': %w", templatePath, err)
			}

			break
		}
	}

	if !found {
		return nil, ErrGeneratorTemplateNotFound
	}

	timeSpec := NewFakeTimeSpec(g.clock.Now())

	templatePageFileInfo := NewPageFileInfo(
		templatePath,
		"",
		"",
		timeSpec,
	)

	entries, err := g.source.Entries(ctx)
	if err != nil {
		return nil, errors.Join(ErrFailedToGenEntries, err)
	}

	generatorOutput := g.config.Output()

	outputDir := generatorOutput.Dir()
	outputFilenameTemplate := generatorOutput.FilenameTemplate()

	dirConfig := NewPageConfig(
		"generator:source",
		"",
		"",
		"",
		false,
		false,
		false,
		g.source.Variables(),
		nil,
		nil,
		nil,
	)

	pages := make([]Page, 0, len(entries))

	for _, entry := range entries {
		entryFilename := entry.Filename() + templatePageFileInfo.FileExtension
		entryVariables := entry.Variables()

		if outputFilenameTemplate != "" {
			tmpl := template.New(entryFilename)
			if _, err = tmpl.Parse(outputFilenameTemplate); err != nil {
				return nil, fmt.Errorf("failed to parse template filename '%s': %w", outputFilenameTemplate, err)
			}

			writer := bytes.NewBuffer(nil)

			if err = tmpl.Execute(writer, entryVariables); err != nil {
				return nil, fmt.Errorf("failed to execute template filename '%s': %w", outputFilenameTemplate, err)
			}

			entryFilename = writer.String()
		}

		pageFileInfo := NewPageFileInfo(
			filepath.Join(outputDir, entryFilename),
			"",
			"",
			timeSpec,
		)
		pageFileInfo.IsTemplate = templatePageFileInfo.IsTemplate
		pageFileInfo.IsMarkdown = templatePageFileInfo.IsMarkdown
		pageFileInfo.IsHtml = templatePageFileInfo.IsHtml

		page, err := g.createPageFunction(
			ctx,
			pageFileInfo,
			pageContent,
			entryVariables,
			[]PageConfig{
				dirConfig,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create page for entry '%s': %w", entry.Filename(), err)
		}

		pages = append(pages, page)
	}

	return pages, nil
}
