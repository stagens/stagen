package markdown

import (
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// MarkdownSpoiler

type MarkdownSpoiler struct{}

func NewMarkdownSpoiler() *MarkdownSpoiler {
	return &MarkdownSpoiler{}
}

func (p *MarkdownSpoiler) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewMarkdownSpoilerParser(), 500),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewMarkdownSpoilerHTMLRenderer(), 500),
	))
}

// Spoiler

type Spoiler struct {
	gast.BaseInline
}

func NewSpoiler() *Spoiler {
	return &Spoiler{}
}

// Dump implements Node.Dump.
func (n *Spoiler) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

var KindSpoiler = gast.NewNodeKind("Spoiler")

// Kind implements Node.Kind.
func (n *Spoiler) Kind() gast.NodeKind {
	return KindSpoiler
}

// Delimiter Processor

type spoilerDelimiterProcessor struct{}

func (p *spoilerDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '|'
}

func (p *spoilerDelimiterProcessor) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *spoilerDelimiterProcessor) OnMatch(consumes int) gast.Node {
	return NewSpoiler()
}

var defaultSpoilerDelimiterProcessor = &spoilerDelimiterProcessor{}

// Parser

type MarkdownSpoilerParser struct {
	parser.InlineParser
}

func NewMarkdownSpoilerParser() *MarkdownSpoilerParser {
	return &MarkdownSpoilerParser{}
}

func (s *MarkdownSpoilerParser) Trigger() []byte {
	return []byte{'|'}
}

func (s *MarkdownSpoilerParser) Parse(_ gast.Node, block text.Reader, parserContext parser.Context) gast.Node {
	before := block.PrecendingCharacter()

	line, segment := block.PeekLine()

	node := parser.ScanDelimiter(line, before, 1, defaultSpoilerDelimiterProcessor)
	if node == nil || node.OriginalLength > 2 || before == '|' {
		return nil
	}

	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)

	block.Advance(node.OriginalLength)

	parserContext.PushDelimiter(node)

	return node
}

func (s *MarkdownSpoilerParser) CloseBlock(_ gast.Node, _ parser.Context) {
	// nothing to do
}

// HTML Renderer

type MarkdownSpoilerHTMLRenderer struct{}

func NewMarkdownSpoilerHTMLRenderer() *MarkdownSpoilerHTMLRenderer {
	return &MarkdownSpoilerHTMLRenderer{}
}

func (r *MarkdownSpoilerHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindSpoiler, r.renderSpoiler)
}

var SpoilerAttributeFilter = html.GlobalAttributeFilter

func (r *MarkdownSpoilerHTMLRenderer) renderSpoiler(
	writer util.BufWriter,
	source []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		if node.Attributes() != nil {
			_, _ = writer.WriteString(`<span class="spoiler"`) //nolint:errcheck

			html.RenderAttributes(writer, node, SpoilerAttributeFilter)

			_ = writer.WriteByte('>') //nolint:errcheck
		} else {
			_, _ = writer.WriteString(`<span class="spoiler">`) //nolint:errcheck
		}

		_, _ = writer.WriteString(`<span class="spoiler-content">`) //nolint:errcheck
	} else {
		_, _ = writer.WriteString("</span>") //nolint:errcheck
		_, _ = writer.WriteString("</span>") //nolint:errcheck
	}

	return gast.WalkContinue, nil
}
