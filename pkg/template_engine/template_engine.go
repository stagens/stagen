package template_engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	textTemplate "text/template"

	"github.com/pixality-inc/golang-core/json"
	"github.com/pixality-inc/golang-core/logger"
)

var (
	ErrDictRequiresEvenArgsCount = errors.New("dict requires even number of arguments")
	ErrDictKeysMustBeStrings     = errors.New("dict keys must be strings")
)

type TemplateEngine interface {
	HasBlocks(ctx context.Context, content string) (bool, error)
	Execute(ctx context.Context, layout string, content string, data any) ([]byte, error)
	Render(ctx context.Context, name string) ([]byte, error)
	RenderBlock(ctx context.Context, name string, data any) ([]byte, error)
	Import(ctx context.Context, loadType LoadType, name string) ([]byte, error)
	Include(ctx context.Context, name string, data any) ([]byte, error)
}

type Impl struct {
	log                    logger.Loggable
	name                   string
	loader                 Loader
	extraTemplateFunctions textTemplate.FuncMap
	template               Template
	context                context.Context // nolint:containedctx
	data                   any
	initialized            bool
	mutex                  sync.Mutex
}

func New(
	name string,
	format TemplateFormat,
	loader Loader,
) *Impl {
	return NewWithExtraTemplateFunctions(name, format, loader, make(textTemplate.FuncMap))
}

func NewWithExtraTemplateFunctions(
	name string,
	format TemplateFormat,
	loader Loader,
	extraTemplateFunctions textTemplate.FuncMap,
) *Impl {
	var template Template

	switch format {
	case TemplateFormatText:
		template = NewTextTemplate(name)

	case TemplateFormatHtml:
		template = NewHtmlTemplate(name)

	default:
		panic(fmt.Errorf("%w: %s", ErrUnknownTemplateFormat, format))
	}

	return &Impl{
		log: logger.NewLoggableImplWithServiceAndFields(
			"template_engine",
			logger.Fields{
				"name": name,
			},
		),
		name:                   name,
		loader:                 loader,
		extraTemplateFunctions: extraTemplateFunctions,
		template:               template,
		context:                nil,
		data:                   nil,
		initialized:            false,
		mutex:                  sync.Mutex{},
	}
}

func (e *Impl) HasBlocks(ctx context.Context, content string) (bool, error) {
	if err := e.init(ctx); err != nil {
		return false, fmt.Errorf("init: %w", err)
	}

	tmpl, err := e.template.Clone()
	if err != nil {
		return false, fmt.Errorf("clone template: %w", err)
	}

	if err = tmpl.Parse(content); err != nil {
		return false, fmt.Errorf("parse template content: %w", err)
	}

	templates := tmpl.Templates()

	noBlocks := len(templates) == 1 || (len(templates) == 1 && templates[0].Name() == e.name)

	return !noBlocks, nil
}

func (e *Impl) Execute(ctx context.Context, layout string, content string, data any) ([]byte, error) {
	if err := e.init(ctx); err != nil {
		return nil, fmt.Errorf("init: %w", err)
	}

	e.data = data

	tmpl := e.template

	writer := bytes.NewBuffer(nil)

	if layout != "" {
		importResult, err := e.Import(ctx, LoadTypeLayout, layout)
		if err != nil {
			return nil, fmt.Errorf("include layout %s: %w", layout, err)
		}

		writer.Write(importResult)
	}

	if err := tmpl.Parse(content); err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	if err := tmpl.Execute(writer, e.data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	if layout != "" {
		renderResult, err := e.Render(ctx, layout)
		if err != nil {
			return nil, fmt.Errorf("render layout %s: %w", layout, err)
		}

		writer.Write(renderResult)
	}

	return writer.Bytes(), nil
}

func (e *Impl) Render(_ context.Context, name string) ([]byte, error) {
	tmpl := e.template

	writer := bytes.NewBuffer(nil)

	if err := tmpl.ExecuteTemplate(writer, name, e.data); err != nil {
		return nil, fmt.Errorf("render template '%s': %w", name, err)
	}

	return writer.Bytes(), nil
}

func (e *Impl) RenderBlock(_ context.Context, name string, data any) ([]byte, error) {
	tmpl := e.template

	writer := bytes.NewBuffer(nil)

	if err := tmpl.ExecuteTemplate(writer, name, data); err != nil {
		return nil, fmt.Errorf("render block '%s': %w", name, err)
	}

	return writer.Bytes(), nil
}

func (e *Impl) Import(ctx context.Context, loadType LoadType, name string) ([]byte, error) {
	e.log.GetLoggerWithoutContext().Tracef("Import template type %s '%s'", loadType, name)

	content, err := e.loader.Load(ctx, loadType, name)
	if err != nil {
		return nil, fmt.Errorf("import template type %s '%s': %w", loadType, name, err)
	}

	if err = e.template.Parse(content); err != nil {
		return nil, fmt.Errorf("parse template type %s '%s': %w", loadType, name, err)
	}

	return nil, nil
}

func (e *Impl) Include(_ context.Context, name string, data any) ([]byte, error) {
	// @todo
	return []byte("[TODO:INCLUDE:" + name + "]"), nil
}

func (e *Impl) init(ctx context.Context) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.initialized {
		return nil
	}

	e.context = ctx

	e.template.Funcs(textTemplate.FuncMap{
		"dict":       e.dict,
		"json_parse": e.jsonParse,
		"import": func(path string) (string, error) {
			result, err := e.Import(e.context, LoadTypeImport, path)
			if err != nil {
				return "", err
			}

			return string(result), nil
		},
		"include": e.include,
		"render":  e.render,
		"macro":   e.macro,
	})

	e.template.Funcs(e.extraTemplateFunctions)

	e.initialized = true

	return nil
}

func (e *Impl) jsonParse(value string) (any, error) {
	var data any

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return nil, fmt.Errorf("json parse %s: %w", value, err)
	}

	return data, nil
}

func (e *Impl) dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, ErrDictRequiresEvenArgsCount
	}

	dict := make(map[string]any, len(values)/2)

	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, ErrDictKeysMustBeStrings
		}

		dict[key] = values[i+1]
	}

	return dict, nil
}

func (e *Impl) render(name string) (string, error) {
	result, err := e.Render(e.context, name)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (e *Impl) include(name string, data any) (string, error) {
	result, err := e.Include(e.context, name, data)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (e *Impl) macro(name string, uniqueName string, data map[string]any) (string, error) {
	macroContent, err := e.render(uniqueName)
	if err != nil {
		return "", fmt.Errorf("macro '%s' with unique name '%s' render: %w", name, uniqueName, err)
	}

	data["content"] = macroContent

	macroResult, err := e.RenderBlock(e.context, "macro:"+name, data)
	if err != nil {
		return "", fmt.Errorf("macro '%s' with unique name '%s' render: %w", name, uniqueName, err)
	}

	return string(macroResult), nil
}
