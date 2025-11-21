package filetree

import (
	"os"
	"path/filepath"
)

type Entry interface {
	Name() string
	Path() string
	IsDir() bool
	Children() []Entry
}

type EntryImpl struct {
	name     string
	path     string
	isDir    bool
	children []Entry
}

func NewEntry(name string, path string, isDir bool, children []Entry) Entry {
	return &EntryImpl{
		name:     name,
		path:     path,
		isDir:    isDir,
		children: children,
	}
}

func NewEntryFromDirEntry(
	path string,
	dirEntry os.DirEntry,
	children []Entry,
) Entry {
	basename := filepath.Base(dirEntry.Name())
	fullName := filepath.Join(path, basename)
	dirName := filepath.Dir(fullName)

	return NewEntry(
		basename,
		dirName,
		dirEntry.IsDir(),
		children,
	)
}

func (e *EntryImpl) Name() string {
	return e.name
}

func (e *EntryImpl) Path() string {
	return e.path
}

func (e *EntryImpl) IsDir() bool {
	return e.isDir
}

func (e *EntryImpl) Children() []Entry {
	return e.children
}
