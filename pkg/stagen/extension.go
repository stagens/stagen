package stagen

type Extension interface {
	Index() int
	Name() string
	Path() string
	Config() ExtensionConfig
}

type ExtensionImpl struct {
	index      int
	name       string
	path       string
	siteConfig SiteExtensionConfig
	config     ExtensionConfig
}

func NewExtension(
	index int,
	name string,
	path string,
	siteConfig SiteExtensionConfig,
	config ExtensionConfig,
) *ExtensionImpl {
	return &ExtensionImpl{
		index:      index,
		name:       name,
		path:       path,
		siteConfig: siteConfig,
		config:     config,
	}
}

func (e *ExtensionImpl) Index() int {
	return e.index
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
