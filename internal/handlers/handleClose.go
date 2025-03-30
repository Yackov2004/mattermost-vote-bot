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

// handleClose - хэндлер закрытия голосования
func handleClose(application *settings.Application, post *model.Post) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application.Logger.Info().
		Str("userID", post.UserId).
		Str("message", post.Message).
		Msg("Получена команда /poll close")

	parts := strings.Fields(post.Message)
	if len(parts) < 3 {
		application.Logger.Warn().
			Str("userID", post.UserId).
			Msg("Неверный формат команды /poll close")

		SendMsg(application, "Корректный формат запроса: /poll close <pollID>", post.Id)
		return
	}

	pollID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		application.Logger.Warn().
			Str("userID", post.UserId).
			Str("rawPollID", parts[2]).
			Msg("Не удалось распарсить pollID")

		SendMsg(application, "Некорректный pollID", post.Id)
		return
	}

	poll, err := application.Storage.GetPoll(ctx, pollID)
	if err != nil {
		application.Logger.Error().
			Err(err).
			Uint64("pollID", pollID).
			Msg("Голосование не найдено в БД")

		SendMsg(application, fmt.Sprintf("Голосование не найдено: %d", pollID), post.Id)
		return
	}

	// Проверяем, что закрывает создатель
	if poll.OwnerID != post.UserId {
		application.Logger.Warn().
			Str("ownerID", poll.OwnerID).
			Str("requestUserID", post.UserId).
			Uint64("pollID", pollID).
			Msg("Пользователь не является владельцем опроса")

		SendMsg(application, "Только создатель опроса может его закрыть", post.Id)
		return
	}

	poll.Active = false
	err = application.Storage.UpdatePoll(ctx, poll)

	if err != nil {
		application.Logger.Error().
			Err(err).
			Uint64("pollID", pollID).
			Msg("Не удалось обновить поле Active в БД")

		SendMsg(application, "Не удалось закрыть голосование", post.Id)
		return
	}

	application.Logger.Info().
		Str("ownerID", poll.OwnerID).
		Uint64("pollID", pollID).
		Msg("Голосование успешно закрыто")

	SendMsg(application, fmt.Sprintf("Голосование %d закрыто", pollID), post.Id)
}
