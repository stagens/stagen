package html_tokenizer

import (
	"context"
	"io"
)

type Tokenizer interface {
	Tokenize(ctx context.Context, reader io.Reader) ([]Token, error)
}

type Impl struct{}

func NewTokenizer() Tokenizer {
	return &Impl{}
}

func (t *Impl) Tokenize(ctx context.Context, reader io.Reader) ([]Token, error) {
	state := NewState(reader)

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
