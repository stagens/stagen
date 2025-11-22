package stagen

import (
	"github.com/pixality-inc/golang-core/json"

	"stagen/pkg/util"
)

type SiteConfigAuthorYaml struct {
	NameValue    string `env:"NAME"    yaml:"name"`
	EmailValue   string `env:"EMAIL"   yaml:"email"`
	WebsiteValue string `env:"WEBSITE" yaml:"website"`
}

func (c *SiteConfigAuthorYaml) Name() string {
	return c.NameValue
}

func (c *SiteConfigAuthorYaml) Email() string {
	return c.EmailValue
}

func (c *SiteConfigAuthorYaml) Website() string {
	return c.WebsiteValue
}

type SiteConfigLogoYaml struct {
	UrlValue string `env:"URL" yaml:"url"`
}

func (c *SiteConfigLogoYaml) Url() string {
	return c.UrlValue
}

type SiteConfigCopyrightYaml struct {
	YearValue   int    `env:"YEAR"   yaml:"year"`
	TitleValue  string `env:"TITLE"  yaml:"title"`
	RightsValue string `env:"RIGHTS" yaml:"rights"`
}

func (c *SiteConfigCopyrightYaml) Year() int {
	return c.YearValue
}

func (c *SiteConfigCopyrightYaml) Title() string {
	return c.TitleValue
}

func (c *SiteConfigCopyrightYaml) Rights() string {
	return c.RightsValue
}

type SiteConfigTemplateIncludeYaml struct {
	NameValue string `yaml:"name"`
}

type SiteConfigTemplateExtraYaml struct {
	NameValue string `yaml:"name"`
}

type SiteConfigTemplateYaml struct {
	ThemeValue         string                                      `env:"THEME"          env-default:"default"  yaml:"theme"`
	DefaultLayoutValue string                                      `env:"DEFAULT_LAYOUT" env-default:"_default" yaml:"default_layout"`
	VariablesValue     map[string]any                              `yaml:"variables"`
	IncludesValue      map[string][]*SiteConfigTemplateIncludeYaml `yaml:"includes"`
	ExtrasValue        map[string][]*SiteConfigTemplateExtraYaml   `yaml:"extras"`
}

func (c *SiteConfigTemplateYaml) Theme() string {
	return c.ThemeValue
}

func (c *SiteConfigTemplateYaml) DefaultLayout() string {
	return c.DefaultLayoutValue
}

func (c *SiteConfigTemplateYaml) Variables() map[string]any {
	return c.VariablesValue
}

func (c *SiteConfigTemplateYaml) Includes() map[string][]SiteConfigTemplateInclude {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue)
}

func (c *SiteConfigTemplateYaml) Extras() map[string][]SiteConfigTemplateExtra {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue)
}

type SiteExtensionConfigYaml struct {
	NameValue string `yaml:"name"`
}

func (c *SiteExtensionConfigYaml) Name() string {
	return c.NameValue
}

type SiteAggDictConfigYaml struct {
	NameValue string `yaml:"name"`
}

func (c *SiteAggDictConfigYaml) Name() string {
	return c.NameValue
}

type SiteGeneratorConfigYaml struct {
	NameValue string `yaml:"name"`
}

func (c *SiteGeneratorConfigYaml) Name() string {
	return c.NameValue
}

type DatabaseConfigYaml struct {
	NameValue string        `yaml:"name"`
	DataValue []json.Object `yaml:"data"`
}

func (c *DatabaseConfigYaml) Name() string {
	return c.NameValue
}

func (c *DatabaseConfigYaml) Data() []json.Object {
	return c.DataValue
}

type ThemeConfigAuthorYaml struct {
	NameValue    string `yaml:"name"`
	EmailValue   string `yaml:"email"`
	WebsiteValue string `yaml:"website"`
}

func (c *ThemeConfigAuthorYaml) Name() string {
	return c.NameValue
}

func (c *ThemeConfigAuthorYaml) Email() string {
	return c.EmailValue
}

func (c *ThemeConfigAuthorYaml) Website() string {
	return c.WebsiteValue
}

type ThemeConfigYaml struct {
	NameValue          string                                      `yaml:"name"`
	TitleValue         string                                      `yaml:"title"`
	AuthorValue        ThemeConfigAuthorYaml                       `yaml:"author"`
	DefaultLayoutValue string                                      `yaml:"default_layout"`
	VariablesValue     map[string]any                              `yaml:"variables"`
	IncludesValue      map[string][]*SiteConfigTemplateIncludeYaml `yaml:"includes"`
	ExtrasValue        map[string][]*SiteConfigTemplateExtraYaml   `yaml:"extras"`
}

func (c *ThemeConfigYaml) Name() string {
	return c.NameValue
}

func (c *ThemeConfigYaml) Title() string {
	return c.TitleValue
}

func (c *ThemeConfigYaml) Author() ThemeAuthor {
	return &c.AuthorValue
}

func (c *ThemeConfigYaml) DefaultLayout() string {
	return c.DefaultLayoutValue
}

func (c *ThemeConfigYaml) Variables() map[string]any {
	return c.VariablesValue
}

func (c *ThemeConfigYaml) Includes() map[string][]SiteConfigTemplateInclude {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue)
}

func (c *ThemeConfigYaml) Extras() map[string][]SiteConfigTemplateExtra {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue)
}

func (c *ThemeConfigYaml) ToPageConfig() PageConfig {
	return NewPageConfig(
		c.Name(),
		c.DefaultLayout(),
		"",
		false,
		false,
		c.Variables(),
		c.Includes(),
		c.Extras(),
	)
}

type ConfigYaml struct {
	EnvValue           string `env:"ENV"            env-default:"dev" yaml:"env"`
	WorkDirValue       string `env:"WORK_DIR"       env-default:"."   yaml:"work_dir"`
	BuildDirValue      string `env:"BUILD_DIR"      env-default:""    yaml:"build_dir"`
	DatabasesDirValue  string `env:"DATABASES_DIR"  env-default:""    yaml:"databases_dir"`
	ExtensionsDirValue string `env:"EXTENSIONS_DIR" env-default:""    yaml:"extensions_dir"`
	ThemesDirValue     string `env:"THEMES_DIR"     env-default:""    yaml:"themes_dir"`
	TemplatesDirValue  string `env:"TEMPLATES_DIR"  env-default:""    yaml:"templates_dir"`
	PagesDirValue      string `env:"PAGES_DIR"      env-default:""    yaml:"pages_dir"`
	PublicDirValue     string `env:"PUBLIC_DIR"     env-default:""    yaml:"public_dir"`
}

func (c *ConfigYaml) Env() string {
	return c.EnvValue
}

func (c *ConfigYaml) WorkDir() string {
	return c.WorkDirValue
}

func (c *ConfigYaml) BuildDir() string {
	return c.BuildDirValue
}

func (c *ConfigYaml) DatabasesDir() string {
	return c.DatabasesDirValue
}

func (c *ConfigYaml) ExtensionsDir() string {
	return c.ExtensionsDirValue
}

func (c *ConfigYaml) ThemesDir() string {
	return c.ThemesDirValue
}

func (c *ConfigYaml) PagesDir() string {
	return c.PagesDirValue
}

func (c *ConfigYaml) PublicDir() string {
	return c.PublicDirValue
}

func (c *ConfigYaml) TemplatesDir() string {
	return c.TemplatesDirValue
}

type SiteConfigYaml struct {
	BaseUrlValue    string                     `env:"BASE_URL"         env-default:"http://127.0.0.1:8080" yaml:"base_url"`
	NameValue       string                     `env:"NAME"             env-default:"My Cool Website"       yaml:"name"`
	AuthorValue     SiteConfigAuthorYaml       `env-prefix:"AUTHOR"    yaml:"author"`
	LogoValue       SiteConfigLogoYaml         `env-prefix:"LOGO"      yaml:"logo"`
	CopyrightValue  SiteConfigCopyrightYaml    `env-prefix:"COPYRIGHT" yaml:"copyright"`
	ExtensionsValue []*SiteExtensionConfigYaml `yaml:"extensions"`
	AggDictsValue   []*SiteAggDictConfigYaml   `yaml:"agg_dicts"`
	GeneratorsValue []*SiteGeneratorConfigYaml `yaml:"generators"`
	TemplateValue   SiteConfigTemplateYaml     `env-prefix:"TEMPLATE"  yaml:"template"`
}

func (c *SiteConfigYaml) BaseUrl() string {
	return c.BaseUrlValue
}

func (c *SiteConfigYaml) Name() string {
	return c.NameValue
}

func (c *SiteConfigYaml) Author() SiteConfigAuthor {
	return &c.AuthorValue
}

func (c *SiteConfigYaml) Logo() SiteConfigLogo {
	return &c.LogoValue
}

func (c *SiteConfigYaml) Copyright() SiteConfigCopyright {
	return &c.CopyrightValue
}

func (c *SiteConfigYaml) Extensions() []SiteExtensionConfig {
	return util.SliceOfRefsToInterfaces[SiteExtensionConfigYaml, SiteExtensionConfig](c.ExtensionsValue)
}

func (c *SiteConfigYaml) AggDicts() []SiteAggDictConfig {
	return util.SliceOfRefsToInterfaces[SiteAggDictConfigYaml, SiteAggDictConfig](c.AggDictsValue)
}

func (c *SiteConfigYaml) Generators() []SiteGeneratorConfig {
	return util.SliceOfRefsToInterfaces[SiteGeneratorConfigYaml, SiteGeneratorConfig](c.GeneratorsValue)
}

func (c *SiteConfigYaml) Template() SiteConfigTemplate {
	return &c.TemplateValue
}

type DirConfigYaml struct {
	ThemeValue     string                                      `yaml:"theme"`
	LayoutValue    string                                      `yaml:"layout"`
	TitleValue     string                                      `yaml:"title"`
	IsHiddenValue  bool                                        `yaml:"is_hidden"`
	isDraftValue   bool                                        `yaml:"is_draft"`
	VariablesValue map[string]any                              `yaml:"variables"`
	IncludesValue  map[string][]*SiteConfigTemplateIncludeYaml `yaml:"includes"`
	ExtrasValue    map[string][]*SiteConfigTemplateExtraYaml   `yaml:"extras"`
}

func (c *DirConfigYaml) Theme() string {
	return c.ThemeValue
}

func (c *DirConfigYaml) Layout() string {
	return c.LayoutValue
}

func (c *DirConfigYaml) Title() string {
	return c.TitleValue
}

func (c *DirConfigYaml) IsHidden() bool {
	return c.IsHiddenValue
}

func (c *DirConfigYaml) IsDraft() bool {
	return c.isDraftValue
}

func (c *DirConfigYaml) Variables() map[string]any {
	return c.VariablesValue
}

func (c *DirConfigYaml) Includes() map[string][]SiteConfigTemplateInclude {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue)
}

func (c *DirConfigYaml) Extras() map[string][]SiteConfigTemplateExtra {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue)
}

type PageConfigYaml struct {
	ThemeValue    string                                      `yaml:"theme"`
	LayoutValue   string                                      `yaml:"layout"`
	TitleValue    string                                      `yaml:"title"`
	IsHiddenValue bool                                        `yaml:"is_hidden"`
	IsDraftValue  bool                                        `yaml:"is_draft_value"`
	IncludesValue map[string][]*SiteConfigTemplateIncludeYaml `yaml:"includes"`
	ExtrasValue   map[string][]*SiteConfigTemplateExtraYaml   `yaml:"extras"`
}

func (c *PageConfigYaml) ToPageConfig(variables map[string]any) PageConfig {
	pageConfig := NewPageConfig(
		c.ThemeValue,
		c.LayoutValue,
		c.TitleValue,
		c.IsHiddenValue,
		c.IsDraftValue,
		variables,
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue),
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue),
	)

	return pageConfig
}
