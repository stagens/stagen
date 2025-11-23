package stagen

import "github.com/pixality-inc/golang-core/json"

type Database interface {
	Name() string
	Data() []json.Object
}

type DatabaseImpl struct {
	name   string
	data   []json.Object
	config DatabaseConfig
}

func NewDatabase(
	name string,
	data []json.Object,
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

func (e *DatabaseImpl) Data() []json.Object {
	return e.data
}
