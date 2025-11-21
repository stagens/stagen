package stagen

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/pixality-inc/golang-core/logger"
)

var (
	ErrInit          = errors.New("init")
	ErrNoName        = errors.New("no name")
	ErrLoadExtension = errors.New("extension load")
	ErrLoadDatabase  = errors.New("database load")
	ErrLoadAggDict   = errors.New("agg dict load")
	ErrLoadPage      = errors.New("page load")
	ErrLoadTheme     = errors.New("theme load")
)

var (
	databaseFilenameRegexp   = regexp.MustCompile(`^(.*)\.(yml|yaml)$`)
	pageIgnoreFilenameRegexp = regexp.MustCompile(`(^\.|^(.*)\.(yml|yaml)$)`)
	templateExtensionRegexp  = regexp.MustCompile(`(\.tmpl)`)
	configFilenames          = []string{
		"config.yml",
		"config.yaml",
	}
)

type Stagen interface {
	Build(ctx context.Context) error
}

type Impl struct {
	log         logger.Loggable
	config      Config
	siteConfig  SiteConfig
	initialized bool
	pages       []Page
	mutex       sync.Mutex
}

func New(
	config Config,
	siteConfig SiteConfig,
) *Impl {
	return &Impl{
		log:         logger.NewLoggableImplWithService("stagen"),
		config:      config,
		siteConfig:  siteConfig,
		initialized: false,
		pages:       make([]Page, 0),
		mutex:       sync.Mutex{},
	}
}

// nolint:unused
func (s *Impl) templatesDir() string {
	dir := s.config.TemplatesDir()
	if dir == "" {
		return filepath.Join(s.workDir(), "templates")
	}

	return dir
}
