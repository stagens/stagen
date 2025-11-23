package html_tokenizer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"

	"golang.org/x/net/html"
)

var (
	ErrUnknownTokenType    = errors.New("unknown token type")
	ErrUnexpectedTokenType = errors.New("unexpected token type")
)

type Position struct {
	depth int
}

type State struct {
	withoutClosingTags []string
	tokens             []Token
	reader             io.Reader
	htmlTokenizer      *html.Tokenizer
	position           Position
	nextTokenBuf       *HtmlToken
}

func NewState(reader io.Reader, withoutClosingTags []string) *State {
	return &State{
		withoutClosingTags: withoutClosingTags,
		tokens:             []Token{},
		reader:             reader,
		htmlTokenizer:      html.NewTokenizer(reader),
		position:           Position{depth: 0},
	}
}

func (s *State) Next(ctx context.Context) (bool, error) {
	token, err := s.nextToken(ctx)
	if err != nil {
		return false, err
	}

	if token == nil {
		return false, nil
	}

	parsedToken, err := s.parseToken(ctx, token)
	if err != nil {
		return false, fmt.Errorf("failed to parse token %s: %w", token.Token.Type, err)
	}

	s.tokens = append(s.tokens, parsedToken)

	return true, nil
}

func (s *State) Tokens() []Token {
	return s.tokens
}

func (s *State) nextToken(_ context.Context) (*HtmlToken, error) {
	if s.nextTokenBuf != nil {
		s.nextTokenBuf = nil

		return s.nextTokenBuf, nil
	}

	tt := s.htmlTokenizer.Next()

	if tt == html.ErrorToken {
		err := s.htmlTokenizer.Err()
		switch {
		case errors.Is(err, io.EOF):
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return nil, nil
		}
	}

	raw := s.htmlTokenizer.Raw()
	rawCopy := make([]byte, len(raw))
	copy(rawCopy, raw)

	token := s.htmlTokenizer.Token()

	return NewHtmlToken(token, rawCopy), nil
}

func (s *State) parseToken(ctx context.Context, token *HtmlToken) (Token, error) {
	switch token.Token.Type {
	case html.DoctypeToken:
		return s.parseDoctypeToken(ctx, token)

	case html.CommentToken:
		return s.parseCommentToken(ctx, token)

	case html.TextToken:
		return s.parseTextToken(ctx, token)

	case html.StartTagToken:
		return s.parseTagToken(ctx, token)

	case html.SelfClosingTagToken:
		return s.parseSelfClosingTagToken(ctx, token)

	case html.EndTagToken:
		return s.parseEndTagToken(ctx, token)

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedTokenType, token.Token.Type)
	}
}

func (s *State) parseDoctypeToken(_ context.Context, token *HtmlToken) (Token, error) {
	result := NewDoctypeToken(token, s.position)

	return result, nil
}

func (s *State) parseCommentToken(_ context.Context, token *HtmlToken) (Token, error) {
	result := NewCommentToken(token, s.position)

	return result, nil
}

func (s *State) parseTextToken(_ context.Context, token *HtmlToken) (Token, error) {
	result := NewTextToken(token, s.position)

	return result, nil
}

func (s *State) parseTagToken(ctx context.Context, token *HtmlToken) (Token, error) {
	position := s.position

	s.position.depth++

	defer func() {
		s.position.depth--
	}()

	var endToken *HtmlToken

	children := make([]Token, 0)

	selfClosing := false

	for {
		//nolint:staticcheck // @todo
		if slices.Contains(s.withoutClosingTags, token.Token.Data) {
			selfClosing = true

			break
		}

		nextToken, err := s.nextToken(ctx)
		if err != nil {
			return nil, err
		}

		if nextToken == nil {
			break
		}

		if nextToken.Token.Type == html.EndTagToken {
			endToken = nextToken

			break
		} else {
			child, err := s.parseToken(ctx, nextToken)
			if err != nil {
				return nil, fmt.Errorf("failed to parse child token %s: %w", nextToken.Token.Type, err)
			}

			children = append(children, child)
		}
	}

	if endToken != nil {
		if token.Token.Data != endToken.Token.Data {
			s.nextTokenBuf = endToken
			endToken = nil
		}

		selfClosing = false
	}

	result := NewTagToken(token, position, endToken, selfClosing, children)

	return result, nil
}

func (s *State) parseSelfClosingTagToken(_ context.Context, token *HtmlToken) (Token, error) {
	result := NewTagToken(token, s.position, nil, true, nil)

	return result, nil
}

func (s *State) parseEndTagToken(_ context.Context, token *HtmlToken) (Token, error) {
	result := NewEndTagToken(token, s.position)

	return result, nil
}
