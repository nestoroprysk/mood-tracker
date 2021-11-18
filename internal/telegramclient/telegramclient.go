package telegramclient

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nestoroprysk/mood-tracker/internal/util"
)

// TelegramClient is an interface for sending text to chat.
type TelegramClient interface {
	SendMessage(text string) (Response, error)
	SendPNG(name string, png io.Reader) (Response, error)
}

// Poster posts an HTTP request.
type Poster interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
	Do(*http.Request) (resp *http.Response, err error)
}

var _ TelegramClient = &telegramClient{}

type telegramClient struct {
	token  string
	chatID string
	client Poster
}

// Config defines the Telegram client.
type Config struct {
	// Token is a Telegram bot token.
	Token string `validate:"required"`
}

// NewTelegramClient creates a Telegram client.
func New(conf Config, chatID int, client Poster) TelegramClient {
	return &telegramClient{
		token:  conf.Token,
		chatID: strconv.Itoa(chatID),
		client: client,
	}
}

// Send sends text to chat.
func (tc telegramClient) SendMessage(text string) (Response, error) {
	response, err := tc.client.PostForm(
		"https://api.telegram.org/bot"+tc.token+"/sendMessage",
		url.Values{
			"chat_id":    {tc.chatID},
			"text":       {util.FormatCode(text)},
			"parse_mode": {"markdown"},
		},
	)
	if err != nil {
		err := fmt.Errorf("error when posting text to the chat %q: %w", tc.chatID, err)
		return Response{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("expecting status code %d for the Telegram response; got %d", http.StatusOK, response.StatusCode)
	}

	result, err := ParseResponse(response.Body)
	if err != nil {
		return Response{}, err
	}

	if !result.Ok {
		return Response{}, fmt.Errorf("expecting ok; got %+v", result)
	}

	if result.ErrorCode != 0 {
		return Response{}, fmt.Errorf("expecting zero exit code; got %+v", result)
	}

	return result, nil
}

func (tc telegramClient) SendPNG(name string, png io.Reader) (Response, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	file, err := w.CreateFormFile("photo", name)
	if err != nil {
		return Response{}, err
	}

	if _, err := io.Copy(file, png); err != nil {
		return Response{}, err
	}

	if err := w.WriteField("chat_id", tc.chatID); err != nil {
		return Response{}, err
	}

	w.Close()

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+tc.token+"/sendPhoto", &b)
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := tc.client.Do(req)
	if err != nil {
		err := fmt.Errorf("error when posting text to the chat %q: %w", tc.chatID, err)
		return Response{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("expecting status code %d for the Telegram response; got %d", http.StatusOK, res.StatusCode)
	}

	result, err := ParseResponse(res.Body)
	if err != nil {
		return Response{}, err
	}

	if !result.Ok {
		return Response{}, fmt.Errorf("expecting ok; got %+v", result)
	}

	if result.ErrorCode != 0 {
		return Response{}, fmt.Errorf("expecting zero exit code; got %+v", result)
	}

	return result, nil
}
