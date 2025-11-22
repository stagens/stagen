package template_engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pixality-inc/golang-core/util"
)

type FsLoader struct {
	includePaths map[LoadType][]string
	extensions   []string
}

func NewFsLoader(includePaths map[LoadType][]string, extensions []string) *FsLoader {
	return &FsLoader{
		includePaths: includePaths,
		extensions:   extensions,
	}
}

func (t *FsLoader) Load(_ context.Context, loadType LoadType, path string) (string, error) {
	includePaths, ok := t.includePaths[loadType]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrLoadTypeNotFound, loadType)
	}

	for _, includePath := range includePaths {
		for _, extension := range t.extensions {
			filename := filepath.Join(includePath, path+extension)

			if _, exists := util.FileExists(filename); !exists {
				continue
			}

			content, err := os.ReadFile(filename)
			if err != nil {
				return "", fmt.Errorf("failed to read file %s: %w", filename, err)
			}

			return string(content), nil
		}
	}

	return "", fmt.Errorf("%w: %s", ErrTemplateNotFound, path)
}
