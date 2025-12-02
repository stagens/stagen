package stagen

import (
	"github.com/pixality-inc/golang-core/http"
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

type SiteConfigTemplateImportYaml struct {
	NameValue string `yaml:"name"`
}

func (c *SiteConfigTemplateImportYaml) Name() string {
	return c.NameValue
}

type SiteConfigTemplateIncludeYaml struct {
	NameValue string `yaml:"name"`
}

func (c *SiteConfigTemplateIncludeYaml) Name() string {
	return c.NameValue
}

type SiteConfigTemplateExtraYaml struct {
	UrlValue     string         `yaml:"url"`
	OptionsValue map[string]any `yaml:"options"`
}

func (c *SiteConfigTemplateExtraYaml) Url() string {
	return c.UrlValue
}

func (c *SiteConfigTemplateExtraYaml) Options() map[string]any {
	return c.OptionsValue
}

type SiteConfigTemplateYaml struct {
	ThemeValue         string                                      `env:"THEME"          env-default:"default"  yaml:"theme"`
	DefaultLayoutValue string                                      `env:"DEFAULT_LAYOUT" env-default:"_default" yaml:"default_layout"`
	VariablesValue     map[string]any                              `yaml:"variables"`
	ImportsValue       map[string][]*SiteConfigTemplateImportYaml  `yaml:"imports"`
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

func (c *SiteConfigTemplateYaml) Imports() map[string][]SiteConfigTemplateImport {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateImportYaml, SiteConfigTemplateImport](c.ImportsValue)
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

type ExtensionConfigAuthorYaml struct {
	NameValue    string `yaml:"name"`
	EmailValue   string `yaml:"email"`
	WebsiteValue string `yaml:"website"`
}

func (c *ExtensionConfigAuthorYaml) Name() string {
	return c.NameValue
}

func (c *ExtensionConfigAuthorYaml) Email() string {
	return c.EmailValue
}

func (c *ExtensionConfigAuthorYaml) Website() string {
	return c.WebsiteValue
}

type ExtensionConfigYaml struct {
	NameValue       string                                      `yaml:"name"`
	TitleValue      string                                      `yaml:"title"`
	AuthorValue     ExtensionConfigAuthorYaml                   `yaml:"author"`
	VariablesValue  map[string]any                              `yaml:"variables"`
	ImportsValue    map[string][]*SiteConfigTemplateImportYaml  `yaml:"imports"`
	IncludesValue   map[string][]*SiteConfigTemplateIncludeYaml `yaml:"includes"`
	ExtrasValue     map[string][]*SiteConfigTemplateExtraYaml   `yaml:"extras"`
	AggDictsValue   []*SiteAggDictConfigYaml                    `yaml:"agg_dicts"`
	GeneratorsValue []*SiteGeneratorConfigYaml                  `yaml:"generators"`
}

func (c *ExtensionConfigYaml) Name() string {
	return c.NameValue
}

func (c *ExtensionConfigYaml) Title() string {
	return c.TitleValue
}

func (c *ExtensionConfigYaml) Author() ExtensionAuthor {
	return &c.AuthorValue
}

func (c *ExtensionConfigYaml) Variables() map[string]any {
	return c.VariablesValue
}

func (c *ExtensionConfigYaml) Imports() map[string][]SiteConfigTemplateImport {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateImportYaml, SiteConfigTemplateImport](c.ImportsValue)
}

func (c *ExtensionConfigYaml) Includes() map[string][]SiteConfigTemplateInclude {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue)
}

func (c *ExtensionConfigYaml) Extras() map[string][]SiteConfigTemplateExtra {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue)
}

func (c *ExtensionConfigYaml) AggDicts() []SiteAggDictConfig {
	return util.SliceOfRefsToInterfaces[SiteAggDictConfigYaml, SiteAggDictConfig](c.AggDictsValue)
}

func (c *ExtensionConfigYaml) Generators() []SiteGeneratorConfig {
	return util.SliceOfRefsToInterfaces[SiteGeneratorConfigYaml, SiteGeneratorConfig](c.GeneratorsValue)
}

func (c *ExtensionConfigYaml) ToPageConfig() PageConfig {
	return NewPageConfig(
		"extension:"+c.Name(),
		"",
		"",
		"",
		false,
		false,
		false,
		c.Variables(),
		c.Imports(),
		c.Includes(),
		c.Extras(),
	)
}

type SiteAggDictConfigYaml struct {
	NameValue string   `yaml:"name"`
	KeysValue []string `yaml:"keys"`
}

func (c *SiteAggDictConfigYaml) Name() string {
	return c.NameValue
}

func (c *SiteAggDictConfigYaml) Keys() []string {
	return c.KeysValue
}

type SiteGeneratorConfigSourceYaml struct {
	TypeValue GeneratorSourceType `yaml:"type"`
	NameValue string              `yaml:"name"`
}

func (c *SiteGeneratorConfigSourceYaml) Type() GeneratorSourceType {
	return c.TypeValue
}

func (c *SiteGeneratorConfigSourceYaml) Name() string {
	return c.NameValue
}

type SiteGeneratorConfigTemplateYaml struct {
	NameValue string `yaml:"name"`
}

func (c *SiteGeneratorConfigTemplateYaml) Name() string {
	return c.NameValue
}

type SiteGeneratorConfigOutputYaml struct {
	DirValue              string `yaml:"dir"`
	FilenameTemplateValue string `yaml:"filename_template"`
}

func (c *SiteGeneratorConfigOutputYaml) Dir() string {
	return c.DirValue
}

func (c *SiteGeneratorConfigOutputYaml) FilenameTemplate() string {
	return c.FilenameTemplateValue
}

type SiteGeneratorConfigYaml struct {
	NameValue     string                          `yaml:"name"`
	SourceValue   SiteGeneratorConfigSourceYaml   `yaml:"source"`
	TemplateValue SiteGeneratorConfigTemplateYaml `yaml:"template"`
	OutputValue   SiteGeneratorConfigOutputYaml   `yaml:"output"`
	DataValue     []json.Object                   `yaml:"data"`
}

func (c *SiteGeneratorConfigYaml) Name() string {
	return c.NameValue
}

func (c *SiteGeneratorConfigYaml) Source() SiteGeneratorConfigSource {
	return &c.SourceValue
}

func (c *SiteGeneratorConfigYaml) Template() SiteGeneratorConfigTemplate {
	return &c.TemplateValue
}

func (c *SiteGeneratorConfigYaml) Output() SiteGeneratorConfigOutput {
	return &c.OutputValue
}

func (c *SiteGeneratorConfigYaml) Data() []json.Object {
	return c.DataValue
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
	ImportsValue       map[string][]*SiteConfigTemplateImportYaml  `yaml:"imports"`
	IncludesValue      map[string][]*SiteConfigTemplateIncludeYaml `yaml:"includes"`
	ExtrasValue        map[string][]*SiteConfigTemplateExtraYaml   `yaml:"extras"`
	AggDictsValue      []*SiteAggDictConfigYaml                    `yaml:"agg_dicts"`
	GeneratorsValue    []*SiteGeneratorConfigYaml                  `yaml:"generators"`
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

func (c *ThemeConfigYaml) Imports() map[string][]SiteConfigTemplateImport {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateImportYaml, SiteConfigTemplateImport](c.ImportsValue)
}

func (c *ThemeConfigYaml) Includes() map[string][]SiteConfigTemplateInclude {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue)
}

func (c *ThemeConfigYaml) Extras() map[string][]SiteConfigTemplateExtra {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue)
}

func (c *ThemeConfigYaml) AggDicts() []SiteAggDictConfig {
	return util.SliceOfRefsToInterfaces[SiteAggDictConfigYaml, SiteAggDictConfig](c.AggDictsValue)
}

func (c *ThemeConfigYaml) Generators() []SiteGeneratorConfig {
	return util.SliceOfRefsToInterfaces[SiteGeneratorConfigYaml, SiteGeneratorConfig](c.GeneratorsValue)
}

func (c *ThemeConfigYaml) ToPageConfig() PageConfig {
	return NewPageConfig(
		"theme:"+c.Name(),
		c.Name(),
		c.DefaultLayout(),
		"",
		false,
		false,
		false,
		c.Variables(),
		c.Imports(),
		c.Includes(),
		c.Extras(),
	)
}

type ConfigSettingsYaml struct {
	UseUriHtmlFileExtensionValue bool `env:"USE_URI_HTML_FILE_EXTENSION" env-default:"true" yaml:"use_uri_html_file_extension"`
}

func (c *ConfigSettingsYaml) UseUriHtmlFileExtension() bool {
	return c.UseUriHtmlFileExtensionValue
}

type ConfigDirsYaml struct {
	WorkValue       string `env:"WORK"       env-default:"." yaml:"work"`
	BuildValue      string `env:"BUILD"      env-default:""  yaml:"build"`
	DatabasesValue  string `env:"DATABASES"  env-default:""  yaml:"databases"`
	ExtensionsValue string `env:"EXTENSIONS" env-default:""  yaml:"extensions"`
	ThemesValue     string `env:"THEMES"     env-default:""  yaml:"themes"`
	TemplatesValue  string `env:"TEMPLATES"  env-default:""  yaml:"templates"`
	PagesValue      string `env:"PAGES"      env-default:""  yaml:"pages"`
	PublicValue     string `env:"PUBLIC"     env-default:""  yaml:"public"`
}

func (c *ConfigDirsYaml) Work() string {
	return c.WorkValue
}

func (c *ConfigDirsYaml) Build() string {
	return c.BuildValue
}

func (c *ConfigDirsYaml) Databases() string {
	return c.DatabasesValue
}

func (c *ConfigDirsYaml) Extensions() string {
	return c.ExtensionsValue
}

func (c *ConfigDirsYaml) Themes() string {
	return c.ThemesValue
}

func (c *ConfigDirsYaml) Pages() string {
	return c.PagesValue
}

func (c *ConfigDirsYaml) Public() string {
	return c.PublicValue
}

func (c *ConfigDirsYaml) Templates() string {
	return c.TemplatesValue
}

type ConfigYaml struct {
	EnvValue      string             `env:"ENV"              env-default:"dev" yaml:"env"`
	HttpValue     http.ConfigYaml    `env-prefix:"HTTP_"     yaml:"http"`
	SettingsValue ConfigSettingsYaml `env-prefix:"SETTINGS_" yaml:"settings"`
	DirsValue     ConfigDirsYaml     `env-prefix:"DIRS_"     yaml:"dirs"`
}

func (c *ConfigYaml) Env() string {
	return c.EnvValue
}

func (c *ConfigYaml) Http() http.Config {
	return &c.HttpValue
}

func (c *ConfigYaml) Settings() SettingsConfig {
	return &c.SettingsValue
}

func (c *ConfigYaml) Dirs() DirsConfig {
	return &c.DirsValue
}

type SiteConfigYaml struct {
	BaseUrlValue     string                     `env:"BASE_URL"         env-default:"http://127.0.0.1:8080"       yaml:"base_url"`
	NameValue        string                     `env:"NAME"             env-default:"My Cool Website"             yaml:"name"`
	DescriptionValue string                     `env:"DESCRIPTION"      env-default:"My Cool Website Description" yaml:"description"`
	LangValue        string                     `env:"LANG"             env-default:"en"                          yaml:"lang"`
	AuthorValue      SiteConfigAuthorYaml       `env-prefix:"AUTHOR"    yaml:"author"`
	LogoValue        SiteConfigLogoYaml         `env-prefix:"LOGO"      yaml:"logo"`
	CopyrightValue   SiteConfigCopyrightYaml    `env-prefix:"COPYRIGHT" yaml:"copyright"`
	ExtensionsValue  []*SiteExtensionConfigYaml `yaml:"extensions"`
	AggDictsValue    []*SiteAggDictConfigYaml   `yaml:"agg_dicts"`
	GeneratorsValue  []*SiteGeneratorConfigYaml `yaml:"generators"`
	TemplateValue    SiteConfigTemplateYaml     `env-prefix:"TEMPLATE"  yaml:"template"`
}

func (c *SiteConfigYaml) BaseUrl() string {
	return c.BaseUrlValue
}

func (c *SiteConfigYaml) Name() string {
	return c.NameValue
}

func (c *SiteConfigYaml) Description() string {
	return c.DescriptionValue
}

func (c *SiteConfigYaml) Lang() string {
	return c.LangValue
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
	IsDraftValue   bool                                        `yaml:"is_draft"`
	IsSystemValue  bool                                        `yaml:"is_system"`
	VariablesValue map[string]any                              `yaml:"variables"`
	ImportsValue   map[string][]*SiteConfigTemplateImportYaml  `yaml:"imports"`
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
	return c.IsDraftValue
}

func (c *DirConfigYaml) Variables() map[string]any {
	return c.VariablesValue
}

func (c *DirConfigYaml) Imports() map[string][]SiteConfigTemplateImport {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateImportYaml, SiteConfigTemplateImport](c.ImportsValue)
}

func (c *DirConfigYaml) Includes() map[string][]SiteConfigTemplateInclude {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue)
}

func (c *DirConfigYaml) Extras() map[string][]SiteConfigTemplateExtra {
	return util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue)
}

func (c *DirConfigYaml) ToPageConfig(dir string) PageConfig {
	pageConfig := NewPageConfig(
		"dir:"+dir,
		c.ThemeValue,
		c.LayoutValue,
		c.TitleValue,
		c.IsHiddenValue,
		c.IsDraftValue,
		c.IsSystemValue,
		c.VariablesValue,
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateImportYaml, SiteConfigTemplateImport](c.ImportsValue),
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue),
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue),
	)

	return pageConfig
}

type PageConfigYaml struct {
	ThemeValue    string                                      `yaml:"theme"`
	LayoutValue   string                                      `yaml:"layout"`
	TitleValue    string                                      `yaml:"title"`
	IsHiddenValue bool                                        `yaml:"is_hidden"`
	IsDraftValue  bool                                        `yaml:"is_draft"`
	IsSystemValue bool                                        `yaml:"is_system"`
	ImportsValue  map[string][]*SiteConfigTemplateImportYaml  `yaml:"imports"`
	IncludesValue map[string][]*SiteConfigTemplateIncludeYaml `yaml:"includes"`
	ExtrasValue   map[string][]*SiteConfigTemplateExtraYaml   `yaml:"extras"`
}

func (c *PageConfigYaml) ToPageConfig(variables map[string]any) PageConfig {
	pageConfig := NewPageConfig(
		"page",
		c.ThemeValue,
		c.LayoutValue,
		c.TitleValue,
		c.IsHiddenValue,
		c.IsDraftValue,
		c.IsSystemValue,
		variables,
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateImportYaml, SiteConfigTemplateImport](c.ImportsValue),
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateIncludeYaml, SiteConfigTemplateInclude](c.IncludesValue),
		util.MapOfSlicesOfRefsToInterfaces[string, SiteConfigTemplateExtraYaml, SiteConfigTemplateExtra](c.ExtrasValue),
	)

	return pageConfig
}
