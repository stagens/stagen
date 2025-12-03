package stagen

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pixality-inc/golang-core/storage"
	"github.com/pixality-inc/golang-core/util"

	"github.com/stagens/stagen/pkg/filetree"
)

func (s *Impl) publicDir() string {
	return filepath.Join(s.workDir, "public")
}

// nolint:gocognit
func (s *Impl) copyPublicFiles(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Infof("Creating public dir...")

	buildPublicDir := s.buildDir()

	if err := s.storage.MkDir(ctx, buildPublicDir); err != nil {
		return fmt.Errorf("failed to create public dir: %w", err)
	}

	log.Info("Copying public files...")

	dirsToCopy := make([]string, 0)

	for _, theme := range s.themes {
		themeDir := theme.Path()
		themePublicDir := filepath.Join(themeDir, "public")

		if exists, err := s.storage.FileExists(ctx, themePublicDir); err != nil {
			return fmt.Errorf("faile to check if file %s exists: %w", themePublicDir, err)
		} else if !exists {
			continue
		}

		dirsToCopy = append(dirsToCopy, themePublicDir)
	}

	for _, extension := range s.extensions {
		extensionDir := extension.Path()
		extensionPublicDir := filepath.Join(extensionDir, "public")

		if exists, err := s.storage.FileExists(ctx, extensionPublicDir); err != nil {
			return fmt.Errorf("faile to check if file %s exists: %w", extensionPublicDir, err)
		} else if !exists {
			continue
		}

		dirsToCopy = append(dirsToCopy, extensionPublicDir)
	}

	publicDir := s.publicDir()

	if exists, err := s.storage.FileExists(ctx, publicDir); err != nil {
		return fmt.Errorf("faile to check if file %s exists: %w", publicDir, err)
	} else if exists {
		dirsToCopy = append(dirsToCopy, publicDir)
	}

	createdDirs := make(map[string]struct{})

	for _, dir := range dirsToCopy {
		log.Debugf("Copying public files from '%s'...", dir)

		tree, err := filetree.Tree(ctx, s.storage, dir, filetree.NoMaxLevel)
		if err != nil {
			return fmt.Errorf("failed to create tree for dir '%s': %w", dir, err)
		}

		err = filetree.Visit(ctx, tree, func(entry filetree.Entry) error {
			if entry.IsDir() {
				entryDir := filepath.Join(entry.Path(), entry.Name())

				entryDir, found := strings.CutPrefix(entryDir, dir+"/")
				if !found {
					entryDir, _ = strings.CutPrefix(entryDir, dir)
				}

				if entryDir != "" {
					entryPublicDir := filepath.Join(buildPublicDir, entryDir)

					if _, exists := createdDirs[entryPublicDir]; exists {
						return nil
					}

					log.Debugf("Creating public directory '%s'...", entryPublicDir)

					if err = s.storage.MkDir(ctx, entryPublicDir); err != nil {
						return fmt.Errorf("failed to create public dir '%s': %w", entryPublicDir, err)
					}

					createdDirs[entryPublicDir] = struct{}{}
				}

				return nil
			}

			entryOriginalFilename := filepath.Join(entry.Path(), entry.Name())
			entryFilename, _ := strings.CutPrefix(entryOriginalFilename, dir+"/")
			entryPublicFilename := filepath.Join(buildPublicDir, entryFilename)

			log.Debugf("Copying public file '%s' to '%s'", entryFilename, entryPublicFilename)

			// @todo!!!!
			localStorage, ok := s.storage.(storage.LocalStorage)
			if !ok {
				return ErrStorageIsNotALocalStorage
			}

			entrySourceLocalPath, err := localStorage.LocalPath(ctx, entryOriginalFilename)
			if err != nil {
				return fmt.Errorf("failed to get local file path: %w", err)
			}

			entryDestLocalPath, err := localStorage.LocalPath(ctx, entryPublicFilename)
			if err != nil {
				return fmt.Errorf("failed to get local file path: %w", err)
			}

			if err = util.CopyFile(entrySourceLocalPath, entryDestLocalPath); err != nil {
				return fmt.Errorf("failed to copy file '%s' to '%s': %w", entryOriginalFilename, entryPublicFilename, err)
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to visit dir '%s': %w", dir, err)
		}
	}

	log.Infof("Public files copied")

	return nil
}
