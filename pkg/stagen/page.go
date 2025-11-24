package stagen

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/djherbis/times"
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

func NewPageFileInfo(
	pageFilename string,
	workDir string,
	pagesDir string,
	stat times.Timespec,
) *PageFileInfo {
	workDir, _ = strings.CutPrefix(workDir, "./")
	pagesDir, _ = strings.CutPrefix(pagesDir, "./")

	pageDir := filepath.Dir(pageFilename)
	pageDirWithoutWorkDir := strings.TrimPrefix(pageDir, workDir+"/")
	fullPageFilenameWithoutExt, fullExt := removeFileExtension(pageFilename)
	pageFilenameWithoutExt := filepath.Base(fullPageFilenameWithoutExt)

	fileExt := strings.TrimPrefix(pageFilename, fullPageFilenameWithoutExt)
	templateExt := filepath.Ext(fullExt)
	pageDirWithoutWorkDirAndPagesDir, _ := strings.CutPrefix(pageDir, pagesDir+"/")
	pageDirWithoutWorkDirAndPagesDir, _ = strings.CutPrefix(pageDirWithoutWorkDirAndPagesDir, pagesDir)

	if !strings.HasPrefix(fileExt, templateExt) {
		fileExt = strings.TrimSuffix(fileExt, templateExt)
	} else {
		templateExt = ""
	}

	isTemplate := templateExtensionRegexp.MatchString(templateExt)
	isMarkdown := markdownExtensionRegexp.MatchString(fileExt)
	isHtml := htmlExtensionRegexp.MatchString(fileExt)

	pageFileInfo := &PageFileInfo{
		Filename:                      pageFilename,
		BaseFilename:                  filepath.Base(pageFilename),
		Path:                          pageDir,
		PathWithoutWorkDir:            pageDirWithoutWorkDir,
		PathWithoutWorkDirAndPagesDir: pageDirWithoutWorkDirAndPagesDir,
		FilenameWithoutExtension:      pageFilenameWithoutExt,
		FullExtension:                 fullExt,
		FileExtension:                 fileExt,
		TemplateExtension:             templateExt,
		IsTemplate:                    isTemplate,
		IsMarkdown:                    isMarkdown,
		IsHtml:                        isHtml,
		CreatedAt:                     stat.BirthTime(),
		ModifiedAt:                    stat.ModTime(),
		AccessedAt:                    stat.AccessTime(),
		ChangedAt:                     stat.ChangeTime(),
	}

	return pageFileInfo
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
