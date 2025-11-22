package stagen

import "github.com/pixality-inc/golang-core/json"

type Config interface {
	Env() string
	WorkDir() string
	BuildDir() string
	DatabasesDir() string
	ExtensionsDir() string
	ThemesDir() string
	TemplatesDir() string
	PagesDir() string
}

type SiteConfigAuthor interface {
	Name() string
	Email() string
	Website() string
}

type SiteConfigLogo interface {
	Url() string
}

type SiteConfigCopyright interface {
	Year() int
	Title() string
	Rights() string
}

// SiteExtensionConfig
//
//nolint:iface
type SiteExtensionConfig interface {
	Name() string
}

// SiteAggDictConfig
//
//nolint:iface
type SiteAggDictConfig interface {
	Name() string
}

// SiteGeneratorConfig
//
//nolint:iface
type SiteGeneratorConfig interface {
	Name() string
}

type DatabaseConfig interface {
	Name() string
	Data() []json.Object
}

type ThemeAuthor interface {
	Name() string
	Email() string
	Website() string
}

type ThemeConfig interface {
	Name() string
	Title() string
	Author() ThemeAuthor
	DefaultLayout() string
	Variables() map[string]any
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
}

type DirConfig interface {
	Title() string
	Layout() string
	Variables() map[string]any
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
}

type PageConfig interface {
	Theme() string
	Layout() string
	Title() string
	IsHidden() bool
	IsDraft() bool
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
}

// SiteConfigTemplateInclude
//
//nolint:iface
type SiteConfigTemplateInclude interface {
	Name() string
}

// SiteConfigTemplateExtra
//
//nolint:iface
type SiteConfigTemplateExtra interface {
	Name() string
}

type SiteConfigTemplate interface {
	Theme() string
	Layout() string
	DefaultLayout() string
	Variables() map[string]any
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
}

type SiteConfig interface {
	BaseUrl() string
	Name() string
	Author() SiteConfigAuthor
	Logo() SiteConfigLogo
	Copyright() SiteConfigCopyright
	Extensions() []SiteExtensionConfig
	AggDicts() []SiteAggDictConfig
	Generators() []SiteGeneratorConfig
	Template() SiteConfigTemplate
}

type PageConfigImpl struct {
	theme    string
	layout   string
	title    string
	isHidden bool
	isDraft  bool
	includes map[string][]SiteConfigTemplateInclude
	extras   map[string][]SiteConfigTemplateExtra
}

func NewPageConfigImpl() *PageConfigImpl {
	return &PageConfigImpl{
		theme:    "",
		layout:   "",
		title:    "",
		isHidden: false,
		isDraft:  false,
		includes: make(map[string][]SiteConfigTemplateInclude),
		extras:   make(map[string][]SiteConfigTemplateExtra),
	}
}

func (p *PageConfigImpl) Theme() string {
	return p.theme
}

func (p *PageConfigImpl) Layout() string {
	return p.layout
}

func (p *PageConfigImpl) Title() string {
	return p.title
}

func (p *PageConfigImpl) IsHidden() bool {
	return p.isHidden
}

func (p *PageConfigImpl) IsDraft() bool {
	return p.isDraft
}

func (p *PageConfigImpl) Includes() map[string][]SiteConfigTemplateInclude {
	return p.includes
}

func (p *PageConfigImpl) Extras() map[string][]SiteConfigTemplateExtra {
	return p.extras
}
