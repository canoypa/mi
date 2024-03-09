package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

func NewSessionId() string {
	return uuid.New().String()
}

type AuthConfig struct {
	Name       string
	Icon       string
	Callback   string
	Permission []string
}

func NewAuthUrl(host, sessionId string, config AuthConfig) url.URL {
	url := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   fmt.Sprintf("/miauth/%s", sessionId),
	}

	q := url.Query()
	if config.Name != "" {
		q.Set("name", config.Name)
	}
	if config.Icon != "" {
		q.Set("icon", config.Icon)
	}
	if config.Callback != "" {
		q.Set("callback", config.Callback)
	}
	if len(config.Permission) > 0 {
		q.Set("permission", strings.Join(config.Permission, ","))
	}

	url.RawQuery = q.Encode()

	return url
}

type TokenResponse struct {
	Token string `json:"token"`
}

func FetchToken(host, sessionId string) (TokenResponse, error) {
	url := url.URL{
		Scheme: "https",
		Host:   host,
		Path: fmt.Sprintf(
			"/api/miauth/%s/check",
			sessionId,
		),
	}

	req := http.Request{
		Method: "POST",
		URL:    &url,
	}

	res, err := http.DefaultClient.Do(&req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TokenResponse{}, err
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return TokenResponse{}, err
	}

	return tokenResponse, nil
}
