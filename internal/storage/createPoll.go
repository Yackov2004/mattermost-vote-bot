package storage

import (
	"context"

	"github.com/rs/zerolog/log"
)

// CreatePoll создаёт запись о голосовании в Tarantool
func (s *Storage) CreatePoll(ctx context.Context, poll *Poll) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	s.lastID.Add(1)
	newID := s.lastID.Load()

	tPoll := []interface{}{
		newID,
		poll.Question,
		poll.Options,
		poll.Active,
		poll.OwnerID,
	}

	_, err := s.conn.Insert("polls", tPoll)
	if err != nil {
		log.Error().Err(err).Msg("CreatePoll error")
		return 0, err
	}

	return newID, nil
}
