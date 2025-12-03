package stagen

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrExtensionConfigNotFound = errors.New("extension config not found")
	ErrExtensionAlreadyExists  = errors.New("extension already exists")
)

func (s *Impl) extensionsDir() string {
	return filepath.Join(s.workDir, "ext")
}

func (s *Impl) loadExtensions(ctx context.Context) error {
	s.log.GetLogger(ctx).Info("Loading extensions...")

	for index, extension := range s.siteConfig.Extensions() {
		if err := s.loadExtension(ctx, index, extension); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadExtension, extension.Name(), err)
		}
	}

	return nil
}

func (s *Impl) loadExtension(
	ctx context.Context,
	extensionIndex int,
	siteExtensionConfig SiteExtensionConfig,
) error {
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

	if err = s.addExtension(extensionIndex, extensionName, extensionDir, siteExtensionConfig, extensionConfig); err != nil {
		return fmt.Errorf("can't add extension '%s': %w", extensionName, err)
	}

	return nil
}

func (s *Impl) getExtensionConfig(ctx context.Context, extensionDir string) (ExtensionConfig, error) {
	configFiles := s.getPossibleConfigFilenames()

	for _, configFilename := range configFiles {
		configFilePath := filepath.Join(extensionDir, configFilename)

		if exists, err := s.storage.FileExists(ctx, configFilePath); err != nil {
			return nil, fmt.Errorf("faile to check if file %s exists: %w", configFilePath, err)
		} else if !exists {
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

	var extensionConfigYaml ExtensionConfigYaml

	if err = yaml.Unmarshal(configContent, &extensionConfigYaml); err != nil {
		return nil, fmt.Errorf("failed to parse extension config file '%s': %w", filename, err)
	}

	return &extensionConfigYaml, nil
}

func (s *Impl) addExtension(
	extensionIndex int,
	extensionName string,
	extensionDir string,
	siteExtensionConfig SiteExtensionConfig,
	extensionConfig ExtensionConfig,
) error {
	s.extensions[extensionName] = NewExtension(
		extensionIndex,
		extensionName,
		extensionDir,
		siteExtensionConfig,
		extensionConfig,
	)

	return nil
}
