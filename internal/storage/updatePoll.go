package storage

import (
	"context"

	"github.com/rs/zerolog/log"
)

// UpdatePoll перезаписывает существующее голосование в Tarantool
func (s *Storage) UpdatePoll(ctx context.Context, poll *Poll) error {

	tPoll := []interface{}{
		poll.ID,
		poll.Question,
		poll.Options,
		poll.Active,
		poll.OwnerID,
	}

	_, err := s.conn.Replace("polls", tPoll)
	if err != nil {
		log.Error().Err(err).Msg("UpdatePoll error")
	}
	return err
}
