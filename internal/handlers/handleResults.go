package handlers

import (
	"context"
	"fmt"
	"mattermost-voting-bot/internal/settings"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

// handleResults - хэндлер для подсчета и отображения результатов голосования
func handleResults(app *settings.Application, post *model.Post) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.Logger.Info().
		Str("userID", post.UserId).
		Str("message", post.Message).
		Msg("Получена команда /poll results")

	parts := strings.Fields(post.Message)
	if len(parts) < 3 {
		app.Logger.Warn().
			Str("userID", post.UserId).
			Msg("Неверный формат команды /poll results")

		SendMsg(app, "Корректный формат запроса: /poll results <pollID>", post.Id)
		return
	}

	pollID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		app.Logger.Warn().
			Str("userID", post.UserId).
			Str("rawPollID", parts[2]).
			Msg("Не удалось распарсить pollID")

		SendMsg(app, "Некорректный pollID", post.Id)
		return
	}

	poll, err := app.Storage.GetPoll(ctx, pollID)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Uint64("pollID", pollID).
			Msg("Не удалось найти голосование в БД")

		SendMsg(app, fmt.Sprintf("Не удалось найти голосование с ID: %d", pollID), post.Id)
		return
	}

	votes, err := app.Storage.GetVotesByPoll(ctx, pollID)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Uint64("pollID", pollID).
			Msg("Не удалось получить голоса из БД")

		SendMsg(app, "Не удалось получить голоса из БД", post.Id)
		return
	}

	// Подсчитываем голоса
	counts := make(map[string]uint64)
	for _, v := range votes {
		counts[v.Option]++
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Результаты голосования ID=%d\n", poll.ID))
	sb.WriteString(fmt.Sprintf("Вопрос: %s\n", poll.Question))
	for _, opt := range poll.Options {
		sb.WriteString(fmt.Sprintf("  %s: %d голос(ов)\n", opt, counts[opt]))
	}
	if poll.Active {
		sb.WriteString("Статус: Открыто\n")
	} else {
		sb.WriteString("Статус: Закрыто\n")
	}

	app.Logger.Info().
		Uint64("pollID", pollID).
		Msg("Результаты голосования получены успешно")

	SendMsg(app, sb.String(), post.Id)
}
