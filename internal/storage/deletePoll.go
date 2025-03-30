package storage

import (
	"context"

	"github.com/rs/zerolog/log"
)

// DeletePoll удаляет голосование из Tarantool
func (s *Storage) DeletePoll(ctx context.Context, id uint64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	_, err := s.conn.Delete("polls", "primary", []interface{}{id})
	if err != nil {
		log.Error().Err(err).Msg("DeletePoll error")
	}
	return err
}
