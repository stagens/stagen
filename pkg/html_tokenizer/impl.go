package html_tokenizer

import "golang.org/x/net/html"

type TokenImpl struct {
	parent    Token
	tokenType TokenType
	token     *HtmlToken
	position  Position
}

func NewTokenImpl(
	parent Token,
	tokenType TokenType,
	token *HtmlToken,
	position Position,
) *TokenImpl {
	return &TokenImpl{
		parent:    parent,
		tokenType: tokenType,
		token:     token,
		position:  position,
	}
}

func (t *TokenImpl) Type() TokenType {
	return t.tokenType
}

func (t *TokenImpl) Depth() int {
	return t.position.depth
}

func (t *TokenImpl) Parent() Token {
	return t.parent
}

func (t *TokenImpl) Token() html.Token {
	return t.token.Token
}

func (t *TokenImpl) Raw() []byte {
	return t.token.Raw
}

func (t *TokenImpl) Position() Position {
	return t.position
}

func (t *TokenImpl) SetParent(parent Token) {
	t.parent = parent
}

func (t *TokenImpl) isToken() bool {
	return true
}
