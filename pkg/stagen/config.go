package stagen

import (
	"maps"

	"github.com/pixality-inc/golang-core/json"
)

type Config interface {
	Env() string
	WorkDir() string
	BuildDir() string
	DatabasesDir() string
	ExtensionsDir() string
	ThemesDir() string
	TemplatesDir() string
	PagesDir() string
	PublicDir() string
}

//nolint:iface
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
	ToPageConfig() PageConfig
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

//nolint:iface
type ExtensionAuthor interface {
	Name() string
	Email() string
	Website() string
}

//nolint:iface
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
	ToPageConfig() PageConfig
}

//nolint:iface
type DirConfig interface {
	Theme() string
	Layout() string
	Title() string
	IsHidden() bool
	IsDraft() bool
	Variables() map[string]any
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
}

//nolint:iface
type PageConfig interface {
	Theme() string
	Layout() string
	Title() string
	IsHidden() bool
	IsDraft() bool
	Variables() map[string]any
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
	theme     string
	layout    string
	title     string
	isHidden  bool
	isDraft   bool
	variables map[string]any
	includes  map[string][]SiteConfigTemplateInclude
	extras    map[string][]SiteConfigTemplateExtra
}

func NewDefaultPageConfig(variables map[string]any) *PageConfigImpl {
	return NewPageConfig(
		"",
		"",
		"",
		false,
		false,
		variables,
		nil,
		nil,
	)
}

func NewPageConfig(
	theme string,
	layout string,
	title string,
	isHidden bool,
	isDraft bool,
	variables map[string]any,
	includes map[string][]SiteConfigTemplateInclude,
	extras map[string][]SiteConfigTemplateExtra,
) *PageConfigImpl {
	if variables == nil {
		variables = make(map[string]any)
	}

	if includes == nil {
		includes = make(map[string][]SiteConfigTemplateInclude)
	}

	if extras == nil {
		extras = make(map[string][]SiteConfigTemplateExtra)
	}

	return &PageConfigImpl{
		theme:     theme,
		layout:    layout,
		title:     title,
		isHidden:  isHidden,
		isDraft:   isDraft,
		variables: variables,
		includes:  includes,
		extras:    extras,
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

func (p *PageConfigImpl) Variables() map[string]any {
	return p.variables
}

func (p *PageConfigImpl) Includes() map[string][]SiteConfigTemplateInclude {
	return p.includes
}

func (p *PageConfigImpl) Extras() map[string][]SiteConfigTemplateExtra {
	return p.extras
}

func MergePageConfigs(cfg1 PageConfig, cfg2 PageConfig) PageConfig {
	theme := cfg1.Theme()
	if cfg2.Theme() != "" {
		theme = cfg2.Theme()
	}

	layout := cfg1.Layout()
	if cfg2.Layout() != "" {
		layout = cfg2.Layout()
	}

	title := cfg1.Title()
	if cfg2.Title() != "" {
		title = cfg2.Title()
	}

	isHidden := cfg1.IsHidden()
	if cfg2.IsHidden() {
		isHidden = true
	}

	isDraft := cfg2.IsDraft()
	if cfg2.IsDraft() {
		isDraft = true
	}

	variables := cfg1.Variables()
	maps.Copy(variables, cfg2.Variables())

	includes := cfg1.Includes()
	maps.Copy(includes, cfg2.Includes())

	extras := cfg1.Extras()
	maps.Copy(extras, cfg2.Extras())

	return NewPageConfig(
		theme,
		layout,
		title,
		isHidden,
		isDraft,
		variables,
		includes,
		extras,
	)
}
