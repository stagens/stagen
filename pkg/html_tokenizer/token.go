package html_tokenizer

import "golang.org/x/net/html"

type TokenType string

const (
	TokenTypeDoctype TokenType = "doctype"
	TokenTypeComment TokenType = "comment"
	TokenTypeText    TokenType = "text"
	TokenTypeTag     TokenType = "tag"
	TokenTypeEndTag  TokenType = "end_tag"
)

type Token interface {
	Type() TokenType
	Parent() Token
	Depth() int
	Token() html.Token
	Raw() []byte
	isToken() bool
	SetParent(parent Token)
}
