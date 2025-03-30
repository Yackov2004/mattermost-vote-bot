package storage

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
)

// GetVotesByPoll возвращает голоса для голосования
func (s *Storage) GetVotesByPoll(ctx context.Context, pollID uint64) ([]Vote, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result []Vote
	err := s.conn.SelectTyped("poll_votes", "poll_id", 0, 100000, tarantool.IterEq, []interface{}{pollID}, &result)
	if err != nil {
		log.Error().Err(err).Msg("GetVotesByPoll error")
		return nil, err
	}
	return result, nil
}
