package template_engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
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
	Import(ctx context.Context, loadType LoadType, name string, withCache bool) ([]byte, error)
	Include(ctx context.Context, name string, data any) ([]byte, error)
}

type Impl struct {
	log                    logger.Loggable
	name                   string
	loader                 Loader
	format                 TemplateFormat
	template               Template
	extraTemplateFunctions textTemplate.FuncMap
	context                context.Context // nolint:containedctx
	data                   any
	imported               map[string]struct{}
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
	tmpl := newTemplate(format, name)

	impl := &Impl{
		log: logger.NewLoggableImplWithServiceAndFields(
			"template_engine",
			logger.Fields{
				"name": name,
			},
		),
		name:                   name,
		loader:                 loader,
		format:                 format,
		template:               tmpl,
		extraTemplateFunctions: extraTemplateFunctions,
		context:                nil,
		data:                   nil,
		imported:               make(map[string]struct{}),
		mutex:                  sync.Mutex{},
	}

	impl.addFuncs(tmpl)

	return impl
}

func (e *Impl) HasBlocks(ctx context.Context, content string) (bool, error) {
	tmpl := newTemplate(e.format, e.name)
	e.addFuncs(tmpl)

	if err := tmpl.Parse(content); err != nil {
		return false, fmt.Errorf("parse template content: %w", err)
	}

	templates := tmpl.Templates()

	noBlocks := len(templates) == 1 || (len(templates) == 1 && templates[0].Name() == e.name)

	return !noBlocks, nil
}

func (e *Impl) Execute(ctx context.Context, layout string, content string, data any) ([]byte, error) {
	e.context = ctx
	e.addFuncs(e.template)

	e.data = data

	tmpl := e.template

	writer := bytes.NewBuffer(nil)

	if layout != "" {
		importResult, err := e.Import(ctx, LoadTypeLayout, layout, true)
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

func (e *Impl) Import(ctx context.Context, loadType LoadType, name string, withCache bool) ([]byte, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	uniqueName := string(loadType) + "::" + name

	if withCache {
		if _, ok := e.imported[uniqueName]; ok {
			return nil, nil
		}
	}

	e.log.GetLogger(ctx).Tracef("Import template type %s '%s'", loadType, name)

	content, err := e.loader.Load(ctx, loadType, name)
	if err != nil {
		return nil, fmt.Errorf("import template type %s '%s': %w", loadType, name, err)
	}

	if err = e.template.Parse(content); err != nil {
		return nil, fmt.Errorf("parse template type %s '%s': %w", loadType, name, err)
	}

	if withCache {
		e.imported[uniqueName] = struct{}{}
	}

	return nil, nil
}

func (e *Impl) Include(ctx context.Context, name string, data any) ([]byte, error) {
	importResult, err := e.Import(ctx, LoadTypeInclude, name, false)
	if err != nil {
		return nil, fmt.Errorf("include '%s': %w", name, err)
	}

	renderResult, err := e.RenderBlock(ctx, name, data)
	if err != nil {
		return nil, fmt.Errorf("render block '%s': %w", name, err)
	}

	result := make([]byte, len(importResult)+len(renderResult))
	copy(result, importResult)
	copy(result, renderResult)

	return result, nil
}

func (e *Impl) addFuncs(tmpl Template) {
	tmpl.Funcs(textTemplate.FuncMap{
		"dict":       e.dict,
		"json_parse": e.jsonParse,
		"has_prefix": strings.HasPrefix,
		"has_suffix": strings.HasSuffix,
		"extends": func(name string) (string, error) {
			result, err := e.Import(e.context, LoadTypeLayout, name, true)
			if err != nil {
				return "", err
			}

			return string(result), nil
		},
		"import": func(name string) (string, error) {
			result, err := e.Import(e.context, LoadTypeImport, name, true)
			if err != nil {
				return "", err
			}

			return string(result), nil
		},
		"include": e.include,
		"render":  e.render,
		"macro":   e.macro,
	})

	tmpl.Funcs(e.extraTemplateFunctions)
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

func newTemplate(format TemplateFormat, name string) Template {
	var template Template

	switch format {
	case TemplateFormatText:
		template = NewTextTemplate(name)

	case TemplateFormatHtml:
		template = NewHtmlTemplate(name)

	default:
		panic(fmt.Errorf("%w: %s", ErrUnknownTemplateFormat, format))
	}

	return template
}
