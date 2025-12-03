---
color: white
colors: [green, blue]
---

Colors: {{ .Databases.colors.Data|len }}
{{- range .Databases.colors.Data }}
- Color: {{ .id }}
{{ end -}}
