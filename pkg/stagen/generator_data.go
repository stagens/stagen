package stagen

import (
	"context"
	"errors"
	"fmt"

	"github.com/pixality-inc/golang-core/json"
)

var (
	ErrDataGeneratorIdNotFound     = errors.New("data generator id not found")
	ErrDataGeneratorIdIsNotAString = errors.New("data generator id is not a string")
)

type DataGeneratorSource struct {
	data []json.Object
}

func NewDataGeneratorSource(
	data []json.Object,
) *DataGeneratorSource {
	return &DataGeneratorSource{
		data: data,
	}
}

func (s *DataGeneratorSource) Entries(_ context.Context) ([]GeneratorSourceEntry, error) {
	entries := make([]GeneratorSourceEntry, len(s.data))

	for index, data := range s.data {
		idValue, ok := data["id"]
		if !ok {
			return nil, fmt.Errorf("%w: index %d", ErrDataGeneratorIdNotFound, index)
		}

		idString, ok := idValue.(string)
		if !ok {
			return nil, fmt.Errorf("%w: index %d, type %T, value %+v", ErrDataGeneratorIdIsNotAString, index, idValue, idValue)
		}

		entry := NewGeneratorSourceEntry(
			idString,
			data,
		)

		entries[index] = entry
	}

	return entries, nil
}

func (s *DataGeneratorSource) Variables() map[string]any {
	return map[string]any{
		"Data": s.data,
	}
}
