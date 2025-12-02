package template_engine

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/pixality-inc/golang-core/logger"
	"github.com/pixality-inc/golang-core/storage"
)

type FsLoader struct {
	log          logger.Loggable
	storage      storage.Storage
	includePaths map[LoadType][]string
	extensions   []string
}

func NewFsLoader(
	storage storage.Storage,
	includePaths map[LoadType][]string,
	extensions []string,
) *FsLoader {
	return &FsLoader{
		log:          logger.NewLoggableImplWithService("template_fs_loader"),
		storage:      storage,
		includePaths: includePaths,
		extensions:   extensions,
	}
}

func (t *FsLoader) Load(ctx context.Context, loadType LoadType, path string) (string, error) {
	includePaths, ok := t.includePaths[loadType]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrLoadTypeNotFound, loadType)
	}

	loadFile := func(filename string) (string, error) {
		file, err := t.storage.ReadFile(ctx, filename)
		if err != nil {
			return "", fmt.Errorf("faile to open file %s: %w", filename, err)
		}

		defer func() {
			if fErr := file.Close(); fErr != nil {
				t.log.GetLogger(ctx).WithError(fErr).Errorf("faile to close file %s: %w", filename, err)
			}
		}()

		content, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", filename, err)
		}

		return string(content), nil
	}

	for _, includePath := range includePaths {
		for _, extension := range t.extensions {
			filename := filepath.Join(includePath, path+extension)

			if exists, err := t.storage.FileExists(ctx, filename); err != nil {
				return "", fmt.Errorf("faile to check if file %s exists: %w", filename, err)
			} else if !exists {
				continue
			}

			return loadFile(filename)
		}
	}

	return "", fmt.Errorf("%w: %s", ErrTemplateNotFound, path)
}
