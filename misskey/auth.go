package misskey

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

func NewMiAuthUrl(host, sessionId string, config MiAuthConfig) url.URL {
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

func MiAuthCheck(host, sessionId string) (MiAuthCheckResponse, error) {
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
		return MiAuthCheckResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return MiAuthCheckResponse{}, err
	}

	var tokenResponse MiAuthCheckResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return MiAuthCheckResponse{}, err
	}

	return tokenResponse, nil
}
