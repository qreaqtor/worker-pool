package web

import (
	"io"
	"net/http"
)

// Выполняет проверку заголовка Content-type и читает body.
func ReadRequestBody(r *http.Request, contentType string) ([]byte, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnknownPayload
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = r.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
