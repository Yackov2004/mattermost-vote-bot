package storage

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tarantool/go-tarantool"
)

// GetPoll возвращает голосование по ID, конвертируя данные из Tarantool в структуру Poll
func (s *Storage) GetPoll(ctx context.Context, id uint64) (*Poll, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result []Poll

	err := s.conn.SelectTyped("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{id}, &result)
	if err != nil {
		log.Error().Err(err).Msg("GetPoll error")
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("не найдено голосование с ID=%d", id)
	}

	tPoll := result[0]

	poll := &Poll{
		ID:       uint64(tPoll.ID),
		Question: tPoll.Question,
		Options:  tPoll.Options,
		Active:   tPoll.Active,
		OwnerID:  tPoll.OwnerID,
	}

	return poll, nil
}
