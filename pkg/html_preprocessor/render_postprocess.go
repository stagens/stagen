package html_preprocessor

import (
	"fmt"
	"strconv"
	"strings"

	"stagen/pkg/html_tokenizer"
)

func (p *Impl) renderPostprocessTokens(tokens []html_tokenizer.Token) ([]byte, error) {
	result := make([]byte, 0)

	for _, token := range tokens {
		tokenResult, err := p.renderPostprocessToken(token)
		if err != nil {
			return nil, err
		}

		result = append(result, tokenResult...)
	}

	return result, nil
}

func (p *Impl) renderPostprocessToken(token html_tokenizer.Token) ([]byte, error) {
	switch tok := token.(type) {
	case *html_tokenizer.DoctypeToken:
		return tok.Raw(), nil
	case *html_tokenizer.CommentToken:
		return nil, nil
	case *html_tokenizer.TextToken:
		return tok.Raw(), nil
	case *html_tokenizer.TagToken:
		if tok.Tag() == "blockquote" {
			var err error

			tok, err = p.postprocessBlockquote(tok)
			if err != nil {
				return nil, err
			}
		}

		result := make([]byte, 0)

		attributes := make([]string, 0)

		for _, attr := range tok.Token().Attr {
			attribute := attr.Key + "=" + strconv.Quote(attr.Val)

			attributes = append(attributes, attribute)
		}

		result = append(result, '<')
		result = append(result, []byte(tok.Tag())...)

		if len(attributes) > 0 {
			result = append(result, ' ')
			result = append(result, []byte(strings.Join(attributes, " "))...)
		}

		if tok.SelfClosing() {
			result = append(result, '/', '>')
		} else {
			result = append(result, '>')
		}

		if !tok.SelfClosing() {
			childrenResult, err := p.renderPostprocessTokens(tok.Children())
			if err != nil {
				return nil, err
			}

			result = append(result, childrenResult...)

			result = append(result, '<', '/')
			result = append(result, []byte(tok.Tag())...)
			result = append(result, '>')
		}

		return result, nil

	case *html_tokenizer.EndTagToken:
		result := make([]byte, 0)

		result = append(result, '<', '/')
		result = append(result, []byte(tok.Tag())...)
		result = append(result, '>')

		return result, nil

	default:
		return nil, fmt.Errorf("%w: '%s' (%T)", ErrUnknownTokenType, tok.Type(), tok)
	}
}

func (p *Impl) postprocessBlockquote(token *html_tokenizer.TagToken) (*html_tokenizer.TagToken, error) {
	var details *html_tokenizer.TagToken

	others := make([]html_tokenizer.Token, 0)

	for _, child := range token.Children() {
		if childTag, ok := child.(*html_tokenizer.TagToken); ok && childTag.Tag() == "details" {
			details = childTag
		} else {
			others = append(others, child)
		}
	}

	newChildren := make([]html_tokenizer.Token, 0)

	if details != nil {
		newDetailsChildren := details.Children()

		newDetailsChildren = append(newDetailsChildren, others...)

		newDetails := html_tokenizer.NewTagToken(
			html_tokenizer.NewHtmlToken(details.Token(), details.Raw()),
			details.Position(),
			nil,
			details.SelfClosing(),
			newDetailsChildren,
		)

		newChildren = append(newChildren, newDetails)
	} else {
		newChildren = append(newChildren, others...)
	}

	newToken := html_tokenizer.NewTagToken(
		html_tokenizer.NewHtmlToken(token.Token(), token.Raw()),
		token.Position(),
		nil,
		token.SelfClosing(),
		newChildren,
	)

	return newToken, nil
}
