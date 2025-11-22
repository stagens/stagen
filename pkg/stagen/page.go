package stagen

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
}

type Page interface {
	Id() string
	Name() string
	Uri() string
	FileInfo() *PageFileInfo
	Title() string
	Theme() string
	Layout() string
	IsHidden() bool
	IsDraft() bool
}

type PageImpl struct {
	id         string
	name       string
	uri        string
	fileInfo   *PageFileInfo
	content    []byte
	variables  map[string]any
	dirConfigs []DirConfig
	config     PageConfig
}

func NewPage(
	id string,
	name string,
	uri string,
	fileInfo *PageFileInfo,
	content []byte,
	variables map[string]any,
	dirConfigs []DirConfig,
	config PageConfig,
) *PageImpl {
	return &PageImpl{
		id:         id,
		name:       name,
		uri:        uri,
		fileInfo:   fileInfo,
		content:    content,
		variables:  variables,
		dirConfigs: dirConfigs,
		config:     config,
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

func (p *PageImpl) Title() string {
	return p.config.Title()
}

func (p *PageImpl) Theme() string {
	return p.config.Theme()
}

func (p *PageImpl) Layout() string {
	return p.config.Layout()
}

func (p *PageImpl) IsHidden() bool {
	return p.config.IsHidden()
}

func (p *PageImpl) IsDraft() bool {
	return p.config.IsDraft()
}
