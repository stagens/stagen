package stagen

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrDatabaseGeneratorIdNotFound     = errors.New("database generator id not found")
	ErrDatabaseGeneratorIdIsNotAString = errors.New("database generator id is not a string")
)

type DatabaseGeneratorSource struct {
	database Database
}

func NewDatabaseGeneratorSource(
	database Database,
) *DatabaseGeneratorSource {
	return &DatabaseGeneratorSource{
		database: database,
	}
}

func (s *DatabaseGeneratorSource) Entries(_ context.Context) ([]GeneratorSourceEntry, error) {
	databaseData := s.database.Data()

	entries := make([]GeneratorSourceEntry, len(databaseData))

	for index, data := range databaseData {
		idValue, ok := data["id"]
		if !ok {
			return nil, fmt.Errorf("%w: index %d", ErrDatabaseGeneratorIdNotFound, index)
		}

		idString, ok := idValue.(string)
		if !ok {
			return nil, fmt.Errorf("%w: index %d, type %T, value %+v", ErrDatabaseGeneratorIdIsNotAString, index, idValue, idValue)
		}

		entry := NewGeneratorSourceEntry(
			idString,
			data,
		)

		entries[index] = entry
	}

	return entries, nil
}

func (s *DatabaseGeneratorSource) Variables() map[string]any {
	return map[string]any{
		"Database": map[string]any{
			"Name": s.database.Name(),
			"Data": s.database.Data(),
		},
	}
}
