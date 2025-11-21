package filetree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const NoMaxLevel = -1

func Tree(ctx context.Context, dir string, maxLevels int) (Entry, error) {
	entries, err := getDirEntries(ctx, dir, 0, maxLevels)
	if err != nil {
		return nil, err
	}

	basename := filepath.Base(dir)
	fullName := filepath.Join(filepath.Dir(dir), basename)
	dirName := filepath.Dir(fullName)

	entry := NewEntry(basename, dirName, true, entries)

	return entry, nil
}

func getDirEntries(ctx context.Context, dir string, level int, maxLevels int) ([]Entry, error) {
	if maxLevels != NoMaxLevel && level+1 > maxLevels {
		return nil, nil
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(dirEntries))

	for _, entry := range dirEntries {
		children := make([]Entry, 0)

		if entry.IsDir() {
			children, err = getDirEntries(ctx, filepath.Join(dir, entry.Name()), level+1, maxLevels)
			if err != nil {
				return nil, fmt.Errorf("error reading directory %s: %w", filepath.Join(dir, entry.Name()), err)
			}
		}

		entries = append(entries, NewEntryFromDirEntry(dir, entry, children))
	}

	return entries, nil
}
