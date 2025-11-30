package stagen

import (
	"context"
	"fmt"
	"strconv"
	"text/template"

	"github.com/pixality-inc/golang-core/json"

	"stagen/pkg/html_preprocessor"
	"stagen/pkg/html_tokenizer"
	"stagen/pkg/markdown"
	"stagen/pkg/template_engine"
)

type Theme interface {
	Name() string

	Path() string

	Config() ThemeConfig

	Render(
		ctx context.Context,
		imports map[string][]SiteConfigTemplateImport,
		layout string,
		content []byte,
		isMarkdown bool,
		data map[string]any,
	) ([]byte, error)
}

type ThemeImpl struct {
	name             string
	path             string
	config           ThemeConfig
	loader           template_engine.Loader
	markdown         markdown.Markdown
	htmlPreprocessor html_preprocessor.HtmlPreprocessor
}

func NewTheme(
	name string,
	path string,
	config ThemeConfig,
	layoutsIncludePaths []string,
	importPaths []string,
	includePaths []string,
) *ThemeImpl {
	templateLoader := template_engine.NewFsLoader(
		map[template_engine.LoadType][]string{
			template_engine.LoadTypeLayout:  layoutsIncludePaths,
			template_engine.LoadTypeImport:  importPaths,
			template_engine.LoadTypeInclude: includePaths,
		},
		[]string{
			".html.tmpl",
		},
	)

	macroWrapper := func(
		macroName string,
		uniqueName string,
		attributes map[string]any,
	) (*html_preprocessor.MacroWrapperResult, error) {
		jsonAttributes, err := json.Marshal(attributes)
		if err != nil {
			return nil, err
		}

		wrapperResult := &html_preprocessor.MacroWrapperResult{
			Before: fmt.Appendf(nil, `{{- define %s }}`, strconv.Quote(uniqueName)),
			After:  []byte(`{{ end -}}`),
			Call: fmt.Appendf(
				nil,
				`{{ macro_render %s %s (%s|json_parse) }}`,
				strconv.Quote(macroName),
				strconv.Quote(uniqueName),
				strconv.Quote(string(jsonAttributes)),
			),
		}

		return wrapperResult, nil
	}

	addClosingTags := []string{"no"}

	withoutClosingTags := make([]string, 0, len(html_tokenizer.WithoutClosingTags)+len(addClosingTags))
	copy(withoutClosingTags, html_tokenizer.WithoutClosingTags)
	withoutClosingTags = append(withoutClosingTags, addClosingTags...)

	return &ThemeImpl{
		name:     name,
		path:     path,
		config:   config,
		loader:   templateLoader,
		markdown: markdown.New(),
		htmlPreprocessor: html_preprocessor.New(
			macroWrapper,
			addClosingTags,
			withoutClosingTags,
			html_preprocessor.AttributesWithoutValue,
		),
	}
}

func (t *ThemeImpl) Name() string {
	return t.name
}

func (t *ThemeImpl) Path() string {
	return t.path
}

func (t *ThemeImpl) Config() ThemeConfig {
	return t.config
}

func (t *ThemeImpl) Render(
	ctx context.Context,
	imports map[string][]SiteConfigTemplateImport,
	layout string,
	content []byte,
	isMarkdown bool,
	data map[string]any,
) ([]byte, error) {
	var templateEngine template_engine.TemplateEngine

	templateEngine = template_engine.NewWithExtraTemplateFunctions(
		t.name,
		template_engine.TemplateFormatText,
		t.loader,
		template.FuncMap{
			"page_content": func() (string, error) {
				return t.renderPageContent(ctx, templateEngine, isMarkdown)
			},
			"markdown": func(text string) (string, error) {
				return t.renderMarkdown(ctx, text)
			},
			"includes": func(includes []SiteConfigTemplateInclude) (string, error) {
				return t.includes(ctx, templateEngine, data, includes)
			},
		},
	)

	importsValues, ok := imports["imports"]
	if !ok {
		importsValues = nil
	}

	for _, importValue := range importsValues {
		if _, err := templateEngine.Import(ctx, template_engine.LoadTypeImport, importValue.Name(), true); err != nil {
			return nil, fmt.Errorf("import '%s': %w", importValue.Name(), err)
		}
	}

	var extras []byte

	extras, content, err := t.htmlPreprocessor.Preprocess(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess content: %w", err)
	}

	contentStr := string(content)

	hasBlocks, err := templateEngine.HasBlocks(ctx, contentStr)
	if err != nil {
		return nil, fmt.Errorf("failed to check block: %w", err)
	}

	if !hasBlocks {
		contentStr = `{{- define "page_content" -}}` + contentStr + `{{- end -}}`
	}

	contentStr = string(extras) + contentStr

	templateResult, err := templateEngine.Execute(ctx, layout, contentStr, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render layout: %w", err)
	}

	templateResult, err = t.htmlPreprocessor.Postprocess(ctx, templateResult)
	if err != nil {
		return nil, fmt.Errorf("failed to postprocess layout: %w", err)
	}

	return templateResult, err
}

func (t *ThemeImpl) renderPageContent(
	ctx context.Context,
	templateEngine template_engine.TemplateEngine,
	isMarkdown bool,
) (string, error) {
	renderResult, err := templateEngine.Render(ctx, "page_content")
	if err != nil {
		return "", err
	}

	if isMarkdown {
		markdownResult, err := t.markdown.Render(renderResult)
		if err != nil {
			return "", fmt.Errorf("failed to render markdown: %w", err)
		}

		return string(markdownResult), nil
	}

	return string(renderResult), err
}

func (t *ThemeImpl) renderMarkdown(
	_ context.Context,
	text string,
) (string, error) {
	markdownResult, err := t.markdown.Render([]byte(text))
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return string(markdownResult), nil
}

func (t *ThemeImpl) includes(
	ctx context.Context,
	templateEngine template_engine.TemplateEngine,
	data map[string]any,
	includes []SiteConfigTemplateInclude,
) (string, error) {
	var results []byte

	for _, includeValue := range includes {
		includeResult, err := templateEngine.Include(ctx, includeValue.Name(), data)
		if err != nil {
			return "", fmt.Errorf("failed to include '%s': %w", includeValue.Name(), err)
		}

		results = append(results, includeResult...)
	}

	return string(results), nil
}
