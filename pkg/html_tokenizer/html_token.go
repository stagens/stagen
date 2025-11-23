package html_tokenizer

import "golang.org/x/net/html"

type HtmlToken struct {
	Token html.Token
	Raw   []byte
}

func NewHtmlToken(token html.Token, raw []byte) *HtmlToken {
	return &HtmlToken{
		Token: token,
		Raw:   raw,
	}
}
