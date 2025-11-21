package stagen

type PageFileInfo struct {
	Filename                 string
	BaseFilename             string
	Path                     string
	PathWithoutWorkDir       string
	FilenameWithoutExtension string
	FullExtension            string
	FileExtension            string
	TemplateExtension        string
	IsTemplate               bool
}

type Page interface {
	Name() string
	FileInfo() *PageFileInfo
}

type PageImpl struct {
	name       string
	fileInfo   *PageFileInfo
	content    []byte
	variables  map[string]any
	dirConfigs []DirConfig
	config     PageConfig
}

func NewPage(
	name string,
	fileInfo *PageFileInfo,
	content []byte,
	variables map[string]any,
	dirConfigs []DirConfig,
	config PageConfig,
) *PageImpl {
	return &PageImpl{
		name:       name,
		fileInfo:   fileInfo,
		content:    content,
		variables:  variables,
		dirConfigs: dirConfigs,
		config:     config,
	}
}

func (p *PageImpl) Name() string {
	return p.name
}

func (p *PageImpl) FileInfo() *PageFileInfo {
	return p.fileInfo
}
