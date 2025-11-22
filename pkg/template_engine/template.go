package template_engine

import (
	"errors"
	"html/template"
	"io"
)

var ErrUnknownTemplateFormat = errors.New("unknown template format")

type TemplateFormat string

const (
	TemplateFormatText TemplateFormat = "text"
	TemplateFormatHtml TemplateFormat = "html"
)

type BasicTemplate interface {
	Name() string
	Execute(writer io.Writer, data any) error
	ExecuteTemplate(writer io.Writer, name string, data any) error
}

type Template interface {
	BasicTemplate

	Parse(content string) error
	Templates() []BasicTemplate
	Funcs(functions template.FuncMap)
	Clone() (Template, error)
}
