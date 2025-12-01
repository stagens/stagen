package markdown

import (
	"bytes"

	callout "github.com/VojtaStruhar/goldmark-obsidian-callout"
	enclave "github.com/quailyquaily/goldmark-enclave"
	"github.com/quailyquaily/goldmark-enclave/core"
	enclaveMark "github.com/quailyquaily/goldmark-enclave/mark"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type Markdown interface {
	Render(content []byte) ([]byte, error)
}

type Impl struct {
	markdown goldmark.Markdown
}

func New() *Impl {
	return &Impl{
		markdown: goldmark.New(
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithExtensions(
				extension.NewTable(extension.WithTableHTMLOptions()),
				NewMarkdownSpoiler(),
				extension.Strikethrough,
				extension.Linkify,
				extension.DefinitionList,
				extension.TaskList,
				callout.ObsidianCallout,
				extension.NewTypographer(
					extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
						extension.LeftDoubleQuote:  []byte(`«`),
						extension.RightDoubleQuote: []byte(`»`),
					}),
				),
				enclave.New(&core.Config{
					DefaultImageAltPrefix: "",
					IframeDisabled:        false,
					VideoDisabled:         true,
					TwitterDisabled:       true,
					TradingViewDisabled:   true,
					DifyWidgetDisabled:    true,
					QuailWidgetDisabled:   true,
				}),
				enclaveMark.New(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
				html.WithUnsafe(),
			),
		),
	}
}

func (m *Impl) Render(content []byte) ([]byte, error) {
	writer := bytes.NewBuffer(nil)

	doc := m.markdown.Parser().Parse(text.NewReader(content))

	// tree, err := toc.Inspect(doc, content)
	// if err != nil {
	// 	return nil, err
	// }

	if err := m.markdown.Renderer().Render(writer, content, doc); err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
