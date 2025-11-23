package main

import (
	"context"
	"os"

	"stagen/pkg/template_engine"
)

var (
	layoutTemplate = `
{{- define "layout" -}}
<html>
	<head>
		<meta charset="utf-8" />
	</head>
	<body>
		<header>
		{{- block "header" . -}}
		DEFAULT HEADER
		{{- end -}}
		</header>
		<main>
		{{- block "main" . -}}
			DEFAULT BODY
			<content>
				{{- block "content" . -}}
				{{- block "page_content" . -}}
					DEFAULT PAGE CONTENT
				{{- end -}}
				{{- end -}}
			</content>
		{{- end -}}
		</main>
		<footer>
		{{- block "footer" . -}}
		DEFAULT FOOTER
		{{- end -}}
		</footer>
	</body>
</html>
{{- end -}}
`

	anotherBodyTemplate = `
{{- define "main" }}
	ANOTHER BODY
	{{ render "content" }}
{{ end -}}`

	testTemplate = `
{{- define "test" }}
	This is a test: {{.A}}
{{ end -}}
`

	childTemplate = `
{{ define "content" }}
	<h1>Hello, {{ .Name }}!</h1>
	<p>Inner content</p>
	{{$Name := "test"}}
	{{- include $Name -}}
	{{- render $Name -}}
{{- end -}}

{{- define "main" -}}
	This is a child body
	{{ render "content" }}
{{- end -}}

{{- define "footer" -}}
	FOOTER FROM CHILD
{{- end -}}
`
)

func main() {
	ctx := context.Background()

	myTemplateLoader := template_engine.NewMapLoader(map[template_engine.LoadType]map[string]string{
		template_engine.LoadTypeLayout: {
			"layout": layoutTemplate,
		},
		template_engine.LoadTypeImport: {
			"test":         testTemplate,
			"another_body": anotherBodyTemplate,
		},
	})

	templateEngine := template_engine.New("default", template_engine.TemplateFormatText, myTemplateLoader)

	templateContent := "My Template"

	if false {
		templateContent = childTemplate
	}

	hasBlocks, err := templateEngine.HasBlocks(ctx, templateContent)
	if err != nil {
		panic(err)
	}

	if !hasBlocks {
		templateContent = `{{- define "page_content" -}}` + templateContent + `{{- end -}}`
	}

	result, err := templateEngine.Execute(ctx, "layout", templateContent, nil)
	if err != nil {
		panic(err)
	}

	_, _ = os.Stdout.Write(result)
}
