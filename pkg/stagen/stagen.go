package stagen

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/pixality-inc/golang-core/logger"

	"stagen/pkg/git"
)

const Version = "0.2.0"

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
	Init(ctx context.Context, cfg Config, siteConfig SiteConfig) error
	NewProject(ctx context.Context, name string) error
	Build(ctx context.Context) error
}

type Impl struct {
	log          logger.Loggable
	config       Config
	siteConfig   SiteConfig
	git          git.Git
	buildTime    time.Time
	initialized  bool
	extensions   map[string]Extension
	databases    map[string]Database
	aggDicts     map[string]SiteAggDictConfig
	aggDictsData map[string]map[string]map[string][]Page
	generators   map[string]Generator
	pages        map[string]Page
	themes       map[string]Theme
	createdDirs  map[string]struct{}
	mutex        sync.Mutex
}

func New(
	gitTool git.Git,
) *Impl {
	return &Impl{
		log:          logger.NewLoggableImplWithService("stagen"),
		git:          gitTool,
		buildTime:    time.Now(),
		initialized:  false,
		extensions:   make(map[string]Extension),
		databases:    make(map[string]Database),
		aggDicts:     make(map[string]SiteAggDictConfig),
		aggDictsData: make(map[string]map[string]map[string][]Page),
		generators:   make(map[string]Generator),
		pages:        make(map[string]Page),
		themes:       make(map[string]Theme),
		createdDirs:  make(map[string]struct{}),
		mutex:        sync.Mutex{},
	}
}

// nolint:unused
func (s *Impl) templatesDir() string {
	dir := s.config.Dirs().Templates()
	if dir == "" {
		return filepath.Join(s.workDir(), "templates")
	}

	return dir
}
