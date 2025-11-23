package template_engine

import (
	"text/template"
	"text/template/parse"

	"stagen/pkg/util"
)

type TextTemplate struct {
	BasicTemplate

	template *template.Template
}

func NewTextTemplate(name string) *TextTemplate {
	return newTextTemplateFromTemplate(template.New(name))
}

func newTextTemplateFromTemplate(tmpl *template.Template) *TextTemplate {
	return &TextTemplate{
		BasicTemplate: tmpl,
		template:      tmpl,
	}
}

func (t *TextTemplate) Parse(content string) error {
	_, err := t.template.Parse(content)

	return err
}

func (t *TextTemplate) Templates() []BasicTemplate {
	return util.SliceOfRefsToInterfaces[template.Template, BasicTemplate](t.template.Templates())
}

func (t *TextTemplate) Funcs(functions template.FuncMap) {
	_ = t.template.Funcs(functions)
}

func (t *TextTemplate) ParseTree() *parse.Tree {
	return t.template.Tree
}
