package storage

import (
	"context"

	"github.com/rs/zerolog/log"
)

// CreateVote создает голос
func (s *Storage) CreateVote(ctx context.Context, v *Vote) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	tuple := []interface{}{
		v.PollID,
		v.UserID,
		v.Option,
	}
	_, err := s.conn.Insert("poll_votes", tuple)
	if err != nil {
		log.Error().Err(err).Msg("CreateVote error")
	}
	return err
}
