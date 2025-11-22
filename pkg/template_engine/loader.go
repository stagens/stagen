package template_engine

import (
	"context"
	"errors"
)

var (
	ErrLoadTypeNotFound = errors.New("load type not found")
	ErrTemplateNotFound = errors.New("template not found")
)

type LoadType string

const (
	LoadTypeLayout  LoadType = "extends"
	LoadTypeExtends LoadType = "extends"
	LoadTypeInclude LoadType = "include"
)

type Loader interface {
	Load(ctx context.Context, loadType LoadType, path string) (string, error)
}
