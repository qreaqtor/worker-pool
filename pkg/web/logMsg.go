package web

import "log/slog"

type logMsg struct {
	URL     string
	Method  string
	Message string
	Status  int
}

// Возвращает структуру, которая пишет логи с помощью logger.
// Остальные поля - информация, которая будет выводиться.
func NewLogMsg(url, method string) *logMsg {
	return &logMsg{
		URL:    url,
		Method: method,
	}
}

func (msg *logMsg) Set(message string, status int) {
	msg.Message = message
	msg.Status = status
}

func (msg *logMsg) Info() {
	slog.Info(msg.Message, getArgs(msg)...)
}

func (msg *logMsg) Error() {
	slog.Error(msg.Message, getArgs(msg)...)
}

func getArgs(msg *logMsg) []any {
	return []any{
		"status", msg.Status,
		"url", msg.URL,
		"method", msg.Method,
	}
}
