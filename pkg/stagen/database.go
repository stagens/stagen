package stagen

type Database interface {
	Name() string
	Data() []any
}

type DatabaseImpl struct {
	name   string
	data   []any
	config DatabaseConfig
}

func NewDatabase(
	name string,
	data []any,
	config DatabaseConfig,
) *DatabaseImpl {
	return &DatabaseImpl{
		name:   name,
		data:   data,
		config: config,
	}
}

func (e *DatabaseImpl) Name() string {
	return e.name
}

func (e *DatabaseImpl) Data() []any {
	return e.data
}
