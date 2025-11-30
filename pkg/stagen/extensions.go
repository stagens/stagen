package stagen

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pixality-inc/golang-core/util"
	"gopkg.in/yaml.v3"
)

var (
	ErrExtensionConfigNotFound = errors.New("extension config not found")
	ErrExtensionAlreadyExists  = errors.New("extension already exists")
	ErrExtensionNotFound       = errors.New("extension not found")
)

// nolint:unused
func (s *Impl) extensionsDir() string {
	dir := s.config.Dirs().Extensions()
	if dir == "" {
		return filepath.Join(s.workDir(), "ext")
	}

	return dir
}

func (s *Impl) loadExtensions(ctx context.Context) error {
	s.log.GetLogger(ctx).Info("Loading extensions...")

	for _, extension := range s.siteConfig.Extensions() {
		if err := s.loadExtension(ctx, extension); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadExtension, extension.Name(), err)
		}
	}

	return nil
}

func (s *Impl) loadExtension(ctx context.Context, siteExtensionConfig SiteExtensionConfig) error {
	extensionName := siteExtensionConfig.Name()
	if extensionName == "" {
		return ErrNoName
	}

	if _, ok := s.extensions[extensionName]; ok {
		return ErrExtensionAlreadyExists
	}

	s.log.GetLogger(ctx).Infof("Loading extension '%s'...", extensionName)

	extensionDir := filepath.Join(s.extensionsDir(), extensionName)

	extensionConfig, err := s.getExtensionConfig(ctx, extensionDir)
	if err != nil {
		return fmt.Errorf("extension '%s': %w", extensionName, err)
	}

	if err = s.addExtension(extensionName, extensionDir, siteExtensionConfig, extensionConfig); err != nil {
		return fmt.Errorf("can't add extension '%s': %w", extensionName, err)
	}

	return nil
}

func (s *Impl) getExtensionConfig(ctx context.Context, extensionDir string) (ExtensionConfig, error) {
	configFiles := s.getPossibleConfigFilenames()

	for _, configFilename := range configFiles {
		configFilePath := filepath.Join(extensionDir, configFilename)

		if _, exists := util.FileExists(configFilePath); !exists {
			continue
		}

		extensionConfig, err := s.readExtensionConfig(ctx, configFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read extension config '%s': %w", configFilePath, err)
		}

		return extensionConfig, nil
	}

	return nil, ErrExtensionConfigNotFound
}

func (s *Impl) readExtensionConfig(ctx context.Context, filename string) (ExtensionConfig, error) {
	configContent, err := s.readFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir config file '%s': %w", filename, err)
	}

	var extensionConfigYaml *ExtensionConfigYaml

	if err = yaml.Unmarshal(configContent, &extensionConfigYaml); err != nil {
		return nil, fmt.Errorf("failed to parse extension config file '%s': %w", filename, err)
	}

	return extensionConfigYaml, nil
}

func (s *Impl) addExtension(
	extensionName string,
	extensionDir string,
	siteExtensionConfig SiteExtensionConfig,
	extensionConfig ExtensionConfig,
) error {
	s.extensions[extensionName] = NewExtension(
		extensionName,
		extensionDir,
		siteExtensionConfig,
		extensionConfig,
	)

	return nil
}
