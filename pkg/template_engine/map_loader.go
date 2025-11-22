package template_engine

import (
	"context"
	"fmt"
)

type MapLoader struct {
	templates map[LoadType]map[string]string
}

func NewMapLoader(templates map[LoadType]map[string]string) *MapLoader {
	return &MapLoader{
		templates: templates,
	}
}

func (t *MapLoader) Load(_ context.Context, loadType LoadType, path string) (string, error) {
	templates, ok := t.templates[loadType]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrLoadTypeNotFound, loadType)
	}

	if content, ok := templates[path]; ok {
		return content, nil
	}

	return "", fmt.Errorf("%w: %s", ErrTemplateNotFound, path)
}
