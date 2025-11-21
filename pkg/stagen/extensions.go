package stagen

import (
	"context"
	"fmt"
	"path/filepath"
)

// nolint:unused
func (s *Impl) extensionsDir() string {
	dir := s.config.ExtensionsDir()
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

func (s *Impl) loadExtension(ctx context.Context, extensionConfig SiteExtensionConfig) error {
	extensionName := extensionConfig.Name()
	if extensionName == "" {
		return ErrNoName
	}

	s.log.GetLogger(ctx).Infof("Loading extension '%s'...", extensionName)

	// @todo
	return nil
}
