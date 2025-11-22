package template_engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	textTemplate "text/template"

	"github.com/pixality-inc/golang-core/logger"
)

var (
	ErrDictRequiresEvenArgsCount = errors.New("dict requires even number of arguments")
	ErrDictKeysMustBeStrings     = errors.New("dict keys must be strings")
)

type TemplateEngine interface {
	HasBlocks(ctx context.Context, content string) (bool, error)
	Execute(ctx context.Context, layout string, content string, data any) ([]byte, error)
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
		_, err := e.includeTemplate(LoadTypeLayout, layout)
		if err != nil {
			return nil, fmt.Errorf("include layout %s: %w", layout, err)
		}
	}

	if err := tmpl.Parse(content); err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	if err := tmpl.Execute(writer, e.data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	if layout != "" {
		renderContents, err := e.render(layout)
		if err != nil {
			return nil, fmt.Errorf("render layout %s: %w", layout, err)
		}

		writer.WriteString(renderContents)
	}

	return writer.Bytes(), nil
}

func (e *Impl) init(ctx context.Context) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.initialized {
		return nil
	}

	e.context = ctx

	e.template.Funcs(textTemplate.FuncMap{
		"dict": e.dict,
		"include": func(path string) (string, error) {
			return e.includeTemplate(LoadTypeInclude, path)
		},
		"render": e.render,
		"macro":  e.macro,
	})

	e.template.Funcs(e.extraTemplateFunctions)

	e.initialized = true

	return nil
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

func (e *Impl) includeTemplate(loadType LoadType, path string) (string, error) {
	e.log.GetLoggerWithoutContext().Tracef("Include template type %s '%s'", loadType, path)

	content, err := e.loader.Load(e.context, loadType, path)
	if err != nil {
		return "", fmt.Errorf("include template type %s '%s': %w", loadType, path, err)
	}

	if err = e.template.Parse(content); err != nil {
		return "", fmt.Errorf("parse template type %s '%s': %w", loadType, path, err)
	}

	return "", nil
}

func (e *Impl) render(name string) (string, error) {
	tmpl := e.template

	writer := bytes.NewBuffer(nil)

	if err := tmpl.ExecuteTemplate(writer, name, e.data); err != nil {
		return "", fmt.Errorf("execute template '%s': %w", name, err)
	}

	return writer.String(), nil
}

func (e *Impl) macro(name string, data any) (string, error) {
	e.log.GetLoggerWithoutContext().Infof("Macro '%s'", name)

	return "<MACRO::" + name + ">", nil
}
