package html_preprocessor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/stagens/stagen/pkg/html_tokenizer"
)

var ErrUnknownTokenType = errors.New("unknown token type")

func (p *Impl) renderTokens(tokens []html_tokenizer.Token) (*TokenRenderResult, error) {
	tokensRenderResult := NewTokenRenderResult()

	for _, token := range tokens {
		tokenRenderResult, err := p.renderToken(token)
		if err != nil {
			return nil, err
		}

		tokensRenderResult.Append(tokenRenderResult)
	}

	return tokensRenderResult, nil
}

func (p *Impl) renderToken(token html_tokenizer.Token) (*TokenRenderResult, error) {
	result := NewTokenRenderResult()

	switch tok := token.(type) {
	case *html_tokenizer.TextToken:
		return result.AppendContent(tok.Raw()), nil

	case *html_tokenizer.TagToken:
		originalTag := tok.Tag()
		uppercasedTag := strings.ToUpper(originalTag)

		if originalTag[0] == uppercasedTag[0] {
			childrenResults, err := p.renderTokens(tok.Children())
			if err != nil {
				return nil, err
			}

			attrs := make(map[string]any)

			for _, attr := range tok.Token().Attr {
				attrs[attr.Key] = attr.Val
			}

			macroName := originalTag

			p.increment++

			contentMacroName := "Content__Macro__" + macroName + "__" + strconv.Itoa(p.increment)

			macroWrapperResult, err := p.macroWrapper(macroName, contentMacroName, attrs)
			if err != nil {
				return nil, fmt.Errorf("macro wrapper: %w", err)
			}

			result = result.AppendExtras(macroWrapperResult.Before)

			result = result.AppendExtras(childrenResults.content)

			result = result.AppendExtras(macroWrapperResult.After)

			result = result.AppendExtras(childrenResults.extras)

			result = result.AppendContent(macroWrapperResult.Call)
		} else {
			result = result.AppendContent(tok.Raw())

			childrenResults, err := p.renderTokens(tok.Children())
			if err != nil {
				return nil, err
			}

			result = result.Append(childrenResults)

			if !tok.SelfClosing() || tok.AddClosing() {
				result = result.AppendContent([]byte("</" + tok.Tag() + ">"))
			}
		}

		return result, nil

	case *html_tokenizer.CommentToken:
		return NewTokenRenderResult(), nil

	default:
		return nil, fmt.Errorf("%w: '%s' (%T)", ErrUnknownTokenType, tok.Type(), tok)
	}
}
