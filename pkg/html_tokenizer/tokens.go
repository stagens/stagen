package html_tokenizer

type DoctypeToken struct {
	*TokenImpl
}

func NewDoctypeToken(token *HtmlToken, position Position) Token {
	return &DoctypeToken{
		TokenImpl: NewTokenImpl(nil, TokenTypeDoctype, token, position),
	}
}

type CommentToken struct {
	*TokenImpl
}

func NewCommentToken(token *HtmlToken, position Position) Token {
	return &CommentToken{
		TokenImpl: NewTokenImpl(nil, TokenTypeComment, token, position),
	}
}

type TextToken struct {
	*TokenImpl
}

func NewTextToken(token *HtmlToken, position Position) Token {
	return &TextToken{
		TokenImpl: NewTokenImpl(nil, TokenTypeText, token, position),
	}
}

type TagToken struct {
	*TokenImpl

	tag         string
	endToken    *HtmlToken
	selfClosing bool
	children    []Token
}

func NewTagToken(
	token *HtmlToken,
	position Position,
	endToken *HtmlToken,
	selfClosing bool,
	children []Token,
) *TagToken {
	tagToken := &TagToken{
		TokenImpl:   NewTokenImpl(nil, TokenTypeTag, token, position),
		tag:         string(token.Raw[1 : len(token.Token.Data)+1]),
		endToken:    endToken,
		selfClosing: selfClosing,
		children:    nil,
	}

	for _, child := range children {
		child.SetParent(tagToken)

		tagToken.children = append(tagToken.children, child)
	}

	return tagToken
}

func (t *TagToken) Tag() string {
	return t.tag
}

func (t *TagToken) Children() []Token {
	return t.children
}

func (t *TagToken) SelfClosing() bool {
	return t.selfClosing
}

type EndTagToken struct {
	*TokenImpl

	tag string
}

func NewEndTagToken(token *HtmlToken, position Position) Token {
	return &EndTagToken{
		TokenImpl: NewTokenImpl(nil, TokenTypeEndTag, token, position),
		tag:       string(token.Raw[1 : len(token.Token.Data)+1]),
	}
}

func (t *EndTagToken) Tag() string {
	return t.tag
}
