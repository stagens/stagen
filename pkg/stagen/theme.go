package stagen

import (
	"context"
	"fmt"

	"stagen/pkg/template_engine"
)

type Theme interface {
	Name() string

	Config() ThemeConfig

	Render(
		ctx context.Context,
		layout string,
		content []byte,
		data any,
	) ([]byte, error)
}

type ThemeImpl struct {
	name   string
	config ThemeConfig
	loader template_engine.Loader
}

func NewTheme(
	name string,
	config ThemeConfig,
	layoutsIncludePaths []string,
	includePaths []string,
) *ThemeImpl {
	templateLoader := template_engine.NewFsLoader(
		map[template_engine.LoadType][]string{
			template_engine.LoadTypeLayout:  layoutsIncludePaths,
			template_engine.LoadTypeInclude: includePaths,
		},
		[]string{
			".html.tmpl",
		},
	)

	return &ThemeImpl{
		name:   name,
		config: config,
		loader: templateLoader,
	}
}

func (t *ThemeImpl) Name() string {
	return t.name
}

func (t *ThemeImpl) Config() ThemeConfig {
	return t.config
}

func (t *ThemeImpl) Render(
	ctx context.Context,
	layout string,
	content []byte,
	data any,
) ([]byte, error) {
	templateEngine := template_engine.New(
		t.name,
		template_engine.TemplateFormatText,
		t.loader,
	)

	contentStr := string(content)

	hasBlocks, err := templateEngine.HasBlocks(ctx, contentStr)
	if err != nil {
		return nil, fmt.Errorf("failed to check block: %w", err)
	}

	if !hasBlocks {
		contentStr = `{{- define "page_content" -}}` + contentStr + `{{- end -}}`
	}

	return templateEngine.Execute(ctx, layout, contentStr, data)
}
