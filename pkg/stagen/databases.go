package stagen

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/stagens/stagen/pkg/filetree"
)

func (s *Impl) databasesDir() string {
	return filepath.Join(s.workDir, "databases")
}

func (s *Impl) loadDatabases(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Info("Loading databases...")

	if exists, err := s.storage.FileExists(ctx, s.databasesDir()); err != nil {
		return fmt.Errorf("failed to check if databases exists: %w", err)
	} else if !exists {
		return nil
	}

	tree, err := filetree.Tree(ctx, s.storage, s.databasesDir(), 1)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	for _, dirEntry := range tree.Children() {
		databaseFilename := filepath.Join(dirEntry.Path(), dirEntry.Name())

		if !databaseFilenameRegexp.MatchString(dirEntry.Name()) {
			log.Debugf("Skipping database file %s...", databaseFilename)

			continue
		}

		if err = s.loadDatabase(ctx, databaseFilename); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadDatabase, databaseFilename, err)
		}
	}

	return nil
}

func (s *Impl) loadDatabase(ctx context.Context, databaseFilename string) error {
	if databaseFilename == "" {
		return ErrNoName
	}

	log := s.log.GetLogger(ctx)

	log.Infof("Loading database %s...", databaseFilename)

	databaseFile, err := s.storage.ReadFile(ctx, databaseFilename)
	if err != nil {
		return fmt.Errorf("%w: failed to read database file '%s': %w", ErrLoadDatabase, databaseFilename, err)
	}

	defer func() {
		if fErr := databaseFile.Close(); fErr != nil {
			log.WithError(fErr).Errorf("Failed to close database file: %s", databaseFilename)
		}
	}()

	databaseContent, err := io.ReadAll(databaseFile)
	if err != nil {
		return fmt.Errorf("%w: failed to read database file '%s': %w", ErrLoadDatabase, databaseFilename, err)
	}

	var databaseYaml DatabaseConfigYaml

	if err = yaml.Unmarshal(databaseContent, &databaseYaml); err != nil {
		return fmt.Errorf("%w: failed to parse database file '%s': %w", ErrLoadDatabase, databaseFilename, err)
	}

	database := NewDatabase(
		databaseYaml.Name(),
		databaseYaml.Data(),
		&databaseYaml,
	)

	s.databases[databaseYaml.Name()] = database

	return nil
}
