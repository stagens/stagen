package stagen

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pixality-inc/golang-core/storage"
)

func (s *Impl) Watch(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Infof("Starting fs notify...")

	watcher, err := s.initWatcher(ctx)
	if err != nil {
		return fmt.Errorf("failed to create fs watcher: %w", err)
	}

	defer func() {
		if wErr := watcher.Close(); wErr != nil {
			log.WithError(wErr).Errorf("failed to close watcher")
		}
	}()

	s.watcherWatch(ctx, watcher)

	return nil
}

func (s *Impl) initWatcher(ctx context.Context) (*fsnotify.Watcher, error) {
	log := s.log.GetLogger(ctx)

	sourceDir := s.workDir
	buildDir := s.buildDir()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create new watcher: %w", err)
	}

	var addDir func(path string)

	addDir = func(dir string) {
		dirClean := path.Clean(dir)

		if dirClean == buildDir {
			return
		}

		// @todo remove prefix of work dir, it's not working

		if (strings.HasPrefix(dirClean, ".") && dirClean != ".") || strings.HasPrefix(dirClean, "~") {
			return
		}

		log.Infof("Adding watching directory: %s", dir)

		// @todo!!!!
		localStorage, ok := s.storage.(storage.LocalStorage)
		if !ok {
			log.WithError(ErrStorageIsNotALocalStorage).Errorf("Storage is not a local storage")

			return
		}

		localDir, err := localStorage.LocalPath(ctx, dir)
		if err != nil {
			log.WithError(err).Errorf("Failed to get local directory: %s", dir)

			return
		}

		if err = watcher.Add(localDir); err != nil {
			log.WithError(err).Errorf("failed to watch directory %s", dir)
		}

		entries, err := s.storage.ReadDir(ctx, dir)
		if err != nil {
			log.WithError(err).Errorf("failed to read directory %s", dir)

			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				addDir(path.Join(dir, entry.Name()))
			}
		}
	}

	addDir(sourceDir)

	return watcher, nil
}

func (s *Impl) watcherWatch(ctx context.Context, watcher *fsnotify.Watcher) {
	log := s.log.GetLogger(ctx)

	buildDir := s.buildDir()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if strings.HasSuffix(event.Name, "~") {
				continue
			}

			doIt := event.Has(fsnotify.Write) ||
				event.Has(fsnotify.Rename) ||
				event.Has(fsnotify.Remove) ||
				event.Has(fsnotify.Create)

			if doIt {
				if path.Clean(event.Name) == buildDir || strings.HasSuffix(path.Clean(event.Name), buildDir) { // @todo!!!!
					continue
				}

				log.Infof("Rebuild becase of %s of %s", event.Op, event.Name)

				if err := s.rebuild(ctx); err != nil {
					log.WithError(err).Errorf("failed to build after fs notify")
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.WithError(err).Errorf("Error")
		}
	}
}

func (s *Impl) rebuild(ctx context.Context) error {
	s.rebuildMutex.Lock()
	defer s.rebuildMutex.Unlock()

	buildDir := s.buildDir()

	if err := s.storage.DeleteDir(ctx, buildDir); err != nil {
		return fmt.Errorf("failed to remove build directory: %w", err)
	}

	s.initialized = false
	s.extensions = make(map[string]Extension)
	s.databases = make(map[string]Database)
	s.aggDicts = make(map[string]SiteAggDictConfig)
	s.aggDictsData = make(map[string]map[string]map[string][]Page)
	s.generators = make(map[string]Generator)
	s.pages = make(map[string]Page)
	s.themes = make(map[string]Theme)
	s.createdDirs = make(map[string]struct{})

	return s.Build(ctx)
}
