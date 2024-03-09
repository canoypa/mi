package misskey

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func NotesCreate(host string, reqBody NotesCreateRequestBody) (data NotesCreateResponse, err error) {
	url := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "api/notes/create",
	}

	bodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return data, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return data, errors.New("status code is not 200")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	var response NotesCreateResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return data, err
	}

	return response, nil
}
