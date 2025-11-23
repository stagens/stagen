package html_preprocessor

import (
	"bytes"
	"context"

	"stagen/pkg/html_tokenizer"
)

type HtmlPreprocessor interface {
	Preprocess(ctx context.Context, content []byte) ([]byte, []byte, error)
	Postprocess(ctx context.Context, content []byte) ([]byte, error)
}

type Impl struct {
	htmlTokenizer html_tokenizer.Tokenizer
	increment     int
	macroWrapper  MacroWrapper
}

func New(macroWrapper MacroWrapper) *Impl {
	return &Impl{
		htmlTokenizer: html_tokenizer.NewTokenizer(),
		increment:     0,
		macroWrapper:  macroWrapper,
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
