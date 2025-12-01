package stagen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"stagen/pkg/filetree"
)

func (s *Impl) databasesDir() string {
	dir := s.config.Dirs().Databases()
	if dir == "" {
		return filepath.Join(s.workDir(), "databases")
	}

	return dir
}

func (s *Impl) loadDatabases(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Info("Loading databases...")

	tree, err := filetree.Tree(ctx, s.databasesDir(), 1)
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

	databaseContent, err := os.ReadFile(databaseFilename)
	if err != nil {
		return fmt.Errorf("%w: failed to read database file '%s': %w", ErrLoadDatabase, databaseFilename, err)
	}

	var databaseYaml *DatabaseConfigYaml

	if err = yaml.Unmarshal(databaseContent, &databaseYaml); err != nil {
		return fmt.Errorf("%w: failed to parse database file '%s': %w", ErrLoadDatabase, databaseFilename, err)
	}

	database := NewDatabase(
		databaseYaml.Name(),
		databaseYaml.Data(),
		databaseYaml,
	)

	s.databases[databaseYaml.Name()] = database

	return nil
}
