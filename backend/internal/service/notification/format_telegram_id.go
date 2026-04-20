package notification

import "strconv"

func formatTelegramID(telegramID int64) string {
	return strconv.FormatInt(telegramID, 10)
}
