package html_preprocessor

type TokenRenderResult struct {
	extras  []byte
	content []byte
}

func NewTokenRenderResult() *TokenRenderResult {
	return &TokenRenderResult{
		extras:  make([]byte, 0),
		content: make([]byte, 0),
	}
}

func (tr *TokenRenderResult) AppendExtras(buf []byte) *TokenRenderResult {
	tr.extras = append(tr.extras, buf...)

	return tr
}

func (tr *TokenRenderResult) AppendContent(buf []byte) *TokenRenderResult {
	tr.content = append(tr.content, buf...)

	return tr
}

func (tr *TokenRenderResult) Append(res *TokenRenderResult) *TokenRenderResult {
	tr.AppendExtras(res.extras)
	tr.AppendContent(res.content)

	return tr
}
