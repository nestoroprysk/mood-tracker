package telegramclient

import (
	"encoding/json"
	"fmt"
	"io"
)

// Response is a Telegram response.
type Response struct {
	Ok          bool    `json:"ok"`
	Result      Message `json:"result"`
	ErrorCode   int     `json:"error_code"`
	Description string  `json:"description"`
}

// ParseResponse parses the response and closes the body.
func ParseResponse(body io.ReadCloser) (Response, error) {
	defer body.Close()

	var response Response
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		err := fmt.Errorf("could not decode an incoming response: %w", err)
		return Response{}, err
	}

	return response, nil
}
