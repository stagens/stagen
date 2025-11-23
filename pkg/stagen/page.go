package stagen

import (
	"sync"
	"time"
)

type PageFileInfo struct {
	Filename                      string
	BaseFilename                  string
	Path                          string
	PathWithoutWorkDir            string
	PathWithoutWorkDirAndPagesDir string
	FilenameWithoutExtension      string
	FullExtension                 string
	FileExtension                 string
	TemplateExtension             string
	IsTemplate                    bool
	IsMarkdown                    bool
	IsHtml                        bool
	CreatedAt                     time.Time
	ModifiedAt                    time.Time
	AccessedAt                    time.Time
	ChangedAt                     time.Time
}

type Page interface {
	Id() string
	Name() string
	Uri() string
	FileInfo() *PageFileInfo
	Config() PageConfig
	Content() []byte
}

type PageImpl struct {
	id           string
	name         string
	uri          string
	fileInfo     *PageFileInfo
	content      []byte
	dirConfigs   []PageConfig
	config       PageConfig
	mergedConfig PageConfig
	mutex        sync.Mutex
}

func NewPage(
	id string,
	name string,
	uri string,
	fileInfo *PageFileInfo,
	content []byte,
	dirConfigs []PageConfig,
	config PageConfig,
) *PageImpl {
	return &PageImpl{
		id:           id,
		name:         name,
		uri:          uri,
		fileInfo:     fileInfo,
		content:      content,
		dirConfigs:   dirConfigs,
		config:       config,
		mergedConfig: nil,
		mutex:        sync.Mutex{},
	}
}

func (p *PageImpl) Id() string {
	return p.id
}

func (p *PageImpl) Name() string {
	return p.name
}

func (p *PageImpl) Uri() string {
	return p.uri
}

func (p *PageImpl) FileInfo() *PageFileInfo {
	return p.fileInfo
}

func (p *PageImpl) Config() PageConfig {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.mergedConfig != nil {
		return p.mergedConfig
	}

	var pageConfig PageConfig = NewDefaultPageConfig("empty", nil)

	for _, dirConfig := range p.dirConfigs {
		pageConfig = MergePageConfigs(pageConfig, dirConfig)
	}

	pageConfig = MergePageConfigs(pageConfig, p.config)

	p.mergedConfig = pageConfig

	return p.mergedConfig
}

func (p *PageImpl) Content() []byte {
	return p.content
}
