package stagen

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (s *Impl) workDir() string {
	return s.config.WorkDir()
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

func (s *Impl) readFile(_ context.Context, filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (s *Impl) getPossibleConfigFilenames() []string {
	configFiles := make([]string, 0, len(configFilenames)+2)

	configFiles = append(configFiles, configFilenames...)

	env := s.config.Env()

	configFiles = append(configFiles, "config."+env+".yaml", "config."+env+".yml")

	return configFiles
}
