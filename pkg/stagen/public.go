package stagen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pixality-inc/golang-core/util"

	"stagen/pkg/filetree"
)

func (s *Impl) publicDir() string {
	dir := s.config.PublicDir()
	if dir == "" {
		return filepath.Join(s.workDir(), "public")
	}

	return dir
}

func (s *Impl) copyPublicFiles(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Infof("Creating public dir...")

	buildPublicDir := filepath.Join(s.buildDir(), "public")

	//nolint:gosec
	if err := os.MkdirAll(buildPublicDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create public dir: %w", err)
	}

	log.Info("Copying public files...")

	dirsToCopy := make([]string, 0)

	for _, theme := range s.themes {
		themeDir := theme.Path()
		themePublicDir := filepath.Join(themeDir, "public")

		if _, exists := util.FileExists(themePublicDir); !exists {
			continue
		}

		dirsToCopy = append(dirsToCopy, themePublicDir)
	}

	for _, extension := range s.extensions {
		extensionDir := extension.Path()
		extensionPublicDir := filepath.Join(extensionDir, "public")

		if _, exists := util.FileExists(extensionPublicDir); !exists {
			continue
		}

		dirsToCopy = append(dirsToCopy, extensionPublicDir)
	}

	publicDir := s.publicDir()

	if _, exists := util.FileExists(publicDir); exists {
		dirsToCopy = append(dirsToCopy, publicDir)
	}

	createdDirs := make(map[string]struct{})

	for _, dir := range dirsToCopy {
		log.Debugf("Copying public files from '%s'...", dir)

		tree, err := filetree.Tree(ctx, dir, filetree.NoMaxLevel)
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

					//nolint:gosec
					if err = os.MkdirAll(entryPublicDir, os.ModePerm); err != nil {
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

			if err = util.CopyFile(entryOriginalFilename, entryPublicFilename); err != nil {
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
