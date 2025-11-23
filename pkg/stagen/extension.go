package stagen

type Extension interface {
	Name() string
	Path() string
	Config() ExtensionConfig
}

type ExtensionImpl struct {
	name       string
	path       string
	siteConfig SiteExtensionConfig
	config     ExtensionConfig
}

func NewExtension(
	name string,
	path string,
	siteConfig SiteExtensionConfig,
	config ExtensionConfig,
) *ExtensionImpl {
	return &ExtensionImpl{
		name:       name,
		path:       path,
		siteConfig: siteConfig,
		config:     config,
	}
}

func (e *ExtensionImpl) Name() string {
	return e.name
}

func (e *ExtensionImpl) Path() string {
	return e.path
}

func (e *ExtensionImpl) Config() ExtensionConfig {
	return e.config
}
