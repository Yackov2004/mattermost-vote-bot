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

// handleDelete - хэндлер удаления голосования
func handleDelete(application *settings.Application, post *model.Post) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application.Logger.Info().
		Str("userID", post.UserId).
		Str("message", post.Message).
		Msg("Получена команда /poll delete")

	parts := strings.Fields(post.Message)
	if len(parts) < 3 {
		application.Logger.Warn().
			Str("userID", post.UserId).
			Msg("Неверный формат команды /poll delete")

		SendMsg(application, "Корректный формат запроса: /poll delete <pollID>", post.Id)
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
			Msg("Не удалось найти голосование в БД")

		SendMsg(application, fmt.Sprintf("Голосование не найдено: %d", pollID), post.Id)
		return
	}

	// Проверяем, что удаляет создатель
	if poll.OwnerID != post.UserId {
		application.Logger.Warn().
			Str("ownerID", poll.OwnerID).
			Str("requestUserID", post.UserId).
			Uint64("pollID", pollID).
			Msg("Пользователь не является владельцем опроса")

		SendMsg(application, "Только создатель опроса может его удалить", post.Id)
		return
	}

	err = application.Storage.DeletePoll(ctx, pollID)
	if err != nil {
		application.Logger.Error().
			Err(err).
			Uint64("pollID", pollID).
			Msg("Не удалось удалить голосование в БД")

		SendMsg(application, fmt.Sprintf("Не удалось удалить голосование %d из БД", pollID), post.Id)
		return
	}

	application.Logger.Info().
		Uint64("pollID", pollID).
		Msg("Голосование успешно удалено")

	SendMsg(application, fmt.Sprintf("Голосование %d удалено", pollID), post.Id)
}
