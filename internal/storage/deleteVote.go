package storage

import (
	"context"

	"github.com/rs/zerolog/log"
)

// DeleteVote удаляет голос
func (s *Storage) DeleteVote(ctx context.Context, pollID uint64, userID, option string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := s.conn.Delete("poll_votes", "poll_id_user_id_option", []interface{}{pollID, userID, option})
	if err != nil {
		log.Error().Err(err).Msg("DeleteVote error")
	}
	return err
}
