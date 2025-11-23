package html_tokenizer

import (
	"context"
	"io"
)

var WithoutClosingTags = []string{"meta", "br", "hr", "img", "input", "link"}

type Tokenizer interface {
	Tokenize(ctx context.Context, reader io.Reader) ([]Token, error)
}

type Impl struct {
	withoutClosingTags []string
}

func NewTokenizer(withoutClosingTags []string) Tokenizer {
	return &Impl{
		withoutClosingTags: withoutClosingTags,
	}
}

func (t *Impl) Tokenize(
	ctx context.Context,
	reader io.Reader,
) ([]Token, error) {
	state := NewState(reader, t.withoutClosingTags)

	for {
		ok, err := state.Next(ctx)
		if err != nil {
			return nil, err
		}

		if !ok {
			break
		}
	}

	return state.Tokens(), nil
}
