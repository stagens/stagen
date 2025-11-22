package stagen

type Extension interface {
	Name() string
	Path() string
	Config() SiteExtensionConfig
}

type ExtensionImpl struct {
	name   string
	path   string
	config SiteExtensionConfig
}

func NewExtension(
	name string,
	path string,
	config SiteExtensionConfig,
) *ExtensionImpl {
	return &ExtensionImpl{
		name:   name,
		path:   path,
		config: config,
	}
}

func (e *ExtensionImpl) Name() string {
	return e.name
}

func (e *ExtensionImpl) Path() string {
	return e.path
}

func (e *ExtensionImpl) Config() SiteExtensionConfig {
	return e.config
}
