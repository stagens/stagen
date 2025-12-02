package stagen

import (
	"fmt"

	"github.com/pixality-inc/golang-core/http"
	"github.com/pixality-inc/golang-core/json"
)

type SettingsConfig interface {
	UseUriHtmlFileExtension() bool
}

type Config interface {
	Env() string
	Http() http.Config
	Settings() SettingsConfig
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
}

// SiteAggDictConfig
//
//nolint:iface
type SiteAggDictConfig interface {
	Name() string
	Keys() []string
}

type GeneratorSourceType string

const (
	GeneratorSourceTypeAggDict  GeneratorSourceType = "agg_dict"
	GeneratorSourceTypeDatabase GeneratorSourceType = "database"
	GeneratorSourceTypeData     GeneratorSourceType = "data"
)

type SiteGeneratorConfigSource interface {
	Type() GeneratorSourceType
	Name() string
}

//nolint:iface
type SiteGeneratorConfigTemplate interface {
	Name() string
}

type SiteGeneratorConfigOutput interface {
	Dir() string
	FilenameTemplate() string
}

// SiteGeneratorConfig
//
//nolint:iface
type SiteGeneratorConfig interface {
	Name() string
	Source() SiteGeneratorConfigSource
	Template() SiteGeneratorConfigTemplate
	Output() SiteGeneratorConfigOutput
	Data() []json.Object
}

//nolint:iface
type DatabaseConfig interface {
	Name() string
	Data() []json.Object
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
	Imports() map[string][]SiteConfigTemplateImport
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
	AggDicts() []SiteAggDictConfig
	Generators() []SiteGeneratorConfig
	ToPageConfig() PageConfig
}

//nolint:iface
type ExtensionAuthor interface {
	Name() string
	Email() string
	Website() string
}

type ExtensionConfig interface {
	Name() string
	Title() string
	Author() ExtensionAuthor
	Variables() map[string]any
	Imports() map[string][]SiteConfigTemplateImport
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
	AggDicts() []SiteAggDictConfig
	Generators() []SiteGeneratorConfig
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
	Imports() map[string][]SiteConfigTemplateImport
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
	ToPageConfig(dir string) PageConfig
}

//nolint:iface
type PageConfig interface {
	ConfigSource() string
	Theme() string
	Layout() string
	Title() string
	IsHidden() bool
	IsDraft() bool
	IsSystem() bool
	Variables() map[string]any
	Imports() map[string][]SiteConfigTemplateImport
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
}

// SiteConfigTemplateImport
//
//nolint:iface
type SiteConfigTemplateImport interface {
	Name() string
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
	Url() string
	Options() map[string]any
}

type SiteConfigTemplate interface {
	Theme() string
	DefaultLayout() string
	Variables() map[string]any
	Imports() map[string][]SiteConfigTemplateImport
	Includes() map[string][]SiteConfigTemplateInclude
	Extras() map[string][]SiteConfigTemplateExtra
}

type SiteConfig interface {
	BaseUrl() string
	Name() string
	Description() string
	Lang() string
	Author() SiteConfigAuthor
	Logo() SiteConfigLogo
	Copyright() SiteConfigCopyright
	Extensions() []SiteExtensionConfig
	AggDicts() []SiteAggDictConfig
	Generators() []SiteGeneratorConfig
	Template() SiteConfigTemplate
}

type PageConfigImpl struct {
	configSource string
	theme        string
	layout       string
	title        string
	isHidden     bool
	isDraft      bool
	isSystem     bool
	variables    map[string]any
	imports      map[string][]SiteConfigTemplateImport
	includes     map[string][]SiteConfigTemplateInclude
	extras       map[string][]SiteConfigTemplateExtra
}

func NewDefaultPageConfig(configSource string, variables map[string]any) *PageConfigImpl {
	return NewPageConfig(
		configSource,
		"",
		"",
		"",
		false,
		false,
		false,
		variables,
		nil,
		nil,
		nil,
	)
}

func NewPageConfig(
	configSource string,
	theme string,
	layout string,
	title string,
	isHidden bool,
	isDraft bool,
	isSystem bool,
	variables map[string]any,
	imports map[string][]SiteConfigTemplateImport,
	includes map[string][]SiteConfigTemplateInclude,
	extras map[string][]SiteConfigTemplateExtra,
) *PageConfigImpl {
	if variables == nil {
		variables = make(map[string]any)
	}

	if imports == nil {
		imports = make(map[string][]SiteConfigTemplateImport)
	}

	if includes == nil {
		includes = make(map[string][]SiteConfigTemplateInclude)
	}

	if extras == nil {
		extras = make(map[string][]SiteConfigTemplateExtra)
	}

	return &PageConfigImpl{
		configSource: configSource,
		theme:        theme,
		layout:       layout,
		title:        title,
		isHidden:     isHidden,
		isDraft:      isDraft,
		isSystem:     isSystem,
		variables:    variables,
		imports:      imports,
		includes:     includes,
		extras:       extras,
	}
}

func (p *PageConfigImpl) ConfigSource() string {
	return p.configSource
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

func (p *PageConfigImpl) IsSystem() bool {
	return p.isSystem
}

func (p *PageConfigImpl) Variables() map[string]any {
	return p.variables
}

func (p *PageConfigImpl) Imports() map[string][]SiteConfigTemplateImport {
	return p.imports
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

	isSystem := cfg2.IsSystem()
	if cfg2.IsSystem() {
		isSystem = true
	}

	variables := cfg1.Variables()
	for k, v := range cfg2.Variables() {
		//nolint:modernize // @todo
		variables[k] = v
	}

	imports := cfg1.Imports()
	for k, v := range cfg2.Imports() {
		imports[k] = append(imports[k], v...)
	}

	includes := cfg1.Includes()
	for k, v := range cfg2.Includes() {
		includes[k] = append(includes[k], v...)
	}

	extras := cfg1.Extras()
	for k, v := range cfg2.Extras() {
		extras[k] = append(extras[k], v...)
	}

	return NewPageConfig(
		fmt.Sprintf("(merged %s :: %s)", cfg1.ConfigSource(), cfg2.ConfigSource()),
		theme,
		layout,
		title,
		isHidden,
		isDraft,
		isSystem,
		variables,
		imports,
		includes,
		extras,
	)
}
