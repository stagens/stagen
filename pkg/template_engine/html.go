package template_engine

import (
	"html/template"

	"stagen/pkg/util"
)

type HtmlTemplate struct {
	BasicTemplate

	template *template.Template
}

func NewHtmlTemplate(name string) *HtmlTemplate {
	return newHtmlTemplateFromTemplate(template.New(name))
}

func newHtmlTemplateFromTemplate(tmpl *template.Template) *HtmlTemplate {
	return &HtmlTemplate{
		BasicTemplate: tmpl,
		template:      tmpl,
	}
}

func (t *HtmlTemplate) Parse(content string) error {
	_, err := t.template.Parse(content)

	return err
}

func (t *HtmlTemplate) Templates() []BasicTemplate {
	return util.SliceOfRefsToInterfaces[template.Template, BasicTemplate](t.template.Templates())
}

func (t *HtmlTemplate) Funcs(functions template.FuncMap) {
	_ = t.template.Funcs(functions)
}

func (t *HtmlTemplate) Clone() (Template, error) {
	tmpl, err := t.template.Clone()
	if err != nil {
		return nil, err
	}

	return newHtmlTemplateFromTemplate(tmpl), nil
}
