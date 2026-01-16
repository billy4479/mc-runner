package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/billy4479/mc-runner/internal/config"
	"github.com/rs/zerolog/log"
)

func sendTelegramMessage(cfg *config.Config, message string) {
	if cfg.EnvConfig.TgToken == "" || cfg.EnvConfig.TgChatIds == "" {
		return
	}

	chatIds := strings.Split(cfg.EnvConfig.TgChatIds, ";")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, chatId := range chatIds {
		chatId = strings.TrimSpace(chatId)
		if chatId == "" {
			continue
		}

		go func(cid string) {
			payload := map[string]string{
				"chat_id": cid,
				"text":    message,
			}
			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal telegram payload")
				return
			}

			url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.EnvConfig.TgToken)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
			if err != nil {
				log.Error().Err(err).Msg("failed to create telegram request")
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Error().Err(err).Msg("failed to send telegram message")
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Error().Int("status_code", resp.StatusCode).Msg("telegram api returned non-200 status")
			}
		}(chatId)
	}
}
