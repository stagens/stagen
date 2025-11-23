package html_preprocessor

import (
	"bytes"
	"context"

	"stagen/pkg/html_tokenizer"
)

var AttributesWithoutValue = []string{"checked", "required", "crossorigin"}

type HtmlPreprocessor interface {
	Preprocess(ctx context.Context, content []byte) ([]byte, []byte, error)
	Postprocess(ctx context.Context, content []byte) ([]byte, error)
}

type Impl struct {
	htmlTokenizer          html_tokenizer.Tokenizer
	increment              int
	macroWrapper           MacroWrapper
	attributesWithoutValue []string
}

func New(
	macroWrapper MacroWrapper,
	addClosingTags []string,
	withoutClosingTags []string,
	attributesWithoutValue []string,
) *Impl {
	return &Impl{
		htmlTokenizer:          html_tokenizer.NewTokenizer(addClosingTags, withoutClosingTags),
		increment:              0,
		macroWrapper:           macroWrapper,
		attributesWithoutValue: attributesWithoutValue,
	}
}

func (p *Impl) Preprocess(ctx context.Context, content []byte) ([]byte, []byte, error) {
	tokens, err := p.htmlTokenizer.Tokenize(ctx, bytes.NewReader(content))
	if err != nil {
		return nil, nil, err
	}

	tokensRenderResult, err := p.renderTokens(tokens)
	if err != nil {
		return nil, nil, err
	}

	return tokensRenderResult.extras, tokensRenderResult.content, nil
}

func (p *Impl) Postprocess(ctx context.Context, content []byte) ([]byte, error) {
	tokens, err := p.htmlTokenizer.Tokenize(ctx, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	result, err := p.renderPostprocessTokens(tokens)
	if err != nil {
		return nil, err
	}

	return result, nil
}
