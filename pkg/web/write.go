package web

import (
	"encoding/json"
	"net/http"
)

/*
Вызывается в случае появления ошибки, пишет msg в логи.
Статус ответа и сообщение достает из msg.
Возвращает 500 в случае неудачной записи в w.
*/
func WriteError(w http.ResponseWriter, msg *logMsg) {
	msg.Error()
	http.Error(w, msg.Message, msg.Status)
}

/*
Выполняет сериализацию data и пишет в w.
В случаае появления ошибки вызывает writeError().
*/
func WriteData(w http.ResponseWriter, msg *logMsg, data any) {
	response, err := json.Marshal(data)
	if err != nil {
		msg.Message = err.Error()
		msg.Status = http.StatusInternalServerError
		WriteError(w, msg)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	_, err = w.Write(response)
	if err != nil {
		msg.Message = err.Error()
		msg.Status = http.StatusInternalServerError
		WriteError(w, msg)
		return
	}
	msg.Info()
}
