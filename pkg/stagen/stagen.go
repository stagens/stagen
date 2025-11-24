package stagen

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/pixality-inc/golang-core/logger"
)

var (
	ErrInit          = errors.New("init")
	ErrNoName        = errors.New("no name")
	ErrLoadExtension = errors.New("extension load")
	ErrLoadDatabase  = errors.New("database load")
	ErrLoadAggDict   = errors.New("agg dict load")
	ErrLoadGenerator = errors.New("generator load")
	ErrLoadPage      = errors.New("page load")
)

var (
	databaseFilenameRegexp   = regexp.MustCompile(`^(.*)\.(yml|yaml)$`)
	pageIgnoreFilenameRegexp = regexp.MustCompile(`(^\.|^(.*)\.(yml|yaml)$)`)
	templateExtensionRegexp  = regexp.MustCompile(`(\.tmpl)`)
	markdownExtensionRegexp  = regexp.MustCompile(`(\.md)`)
	htmlExtensionRegexp      = regexp.MustCompile(`(\.html|\.htm)`)
	configFilenames          = []string{
		"config.yml",
		"config.yaml",
	}
)

type Stagen interface {
	Build(ctx context.Context) error
}

type Impl struct {
	log          logger.Loggable
	config       Config
	siteConfig   SiteConfig
	buildTime    time.Time
	initialized  bool
	extensions   map[string]Extension
	databases    map[string]Database
	aggDictsData map[string]map[string]map[string][]Page
	generators   map[string]any
	pages        map[string]Page
	themes       map[string]Theme
	createdDirs  map[string]struct{}
	mutex        sync.Mutex
}

func New(
	config Config,
	siteConfig SiteConfig,
) *Impl {
	return &Impl{
		log:          logger.NewLoggableImplWithService("stagen"),
		config:       config,
		siteConfig:   siteConfig,
		buildTime:    time.Now(),
		initialized:  false,
		extensions:   make(map[string]Extension),
		databases:    make(map[string]Database),
		aggDictsData: make(map[string]map[string]map[string][]Page),
		generators:   make(map[string]any),
		pages:        make(map[string]Page),
		themes:       make(map[string]Theme),
		createdDirs:  make(map[string]struct{}),
		mutex:        sync.Mutex{},
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
