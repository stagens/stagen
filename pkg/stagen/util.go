package stagen

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

func (s *Impl) readFile(ctx context.Context, filename string) ([]byte, error) {
	storageFile, err := s.storage.ReadFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}

	defer func() {
		if fErr := storageFile.Close(); fErr != nil {
			s.log.GetLogger(ctx).WithError(fErr).Errorf("failed to close storage file: %s", filename)
		}
	}()

	storageFileContent, err := io.ReadAll(storageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return storageFileContent, nil
}

func (s *Impl) getPossibleConfigFilenames() []string {
	configFiles := make([]string, 0, len(configFilenames)+2)

	configFiles = append(configFiles, configFilenames...)

	env := s.config.Env()

	configFiles = append(configFiles, "config."+env+".yaml", "config."+env+".yml")

	return configFiles
}

func removeFileExtension(filename string) (string, string) {
	resultExtensions := make([]string, 0)
	resultFilename := filename

	for {
		ext := filepath.Ext(resultFilename)
		if ext == "" {
			break
		}

		resultFilename = strings.TrimSuffix(resultFilename, ext)

		resultExtensions = append(resultExtensions, ext)
	}

	slices.Reverse(resultExtensions)

	return resultFilename, strings.Join(resultExtensions, "")
}
