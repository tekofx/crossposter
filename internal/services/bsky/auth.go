package bsky

import (
	"bytes"
	"encoding/json"
	"net/http"

	merrors "github.com/tekofx/crossposter/internal/errors"
)

func authenticate() *merrors.MError {

	loginPayload := map[string]string{
		"identifier": BskyClient.Handle,
		"password":   BskyClient.Password,
	}
	body, _ := json.Marshal(loginPayload)
	resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return merrors.New(merrors.BskyAuthRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return merrors.New(merrors.BskyAuthErrorCode, resp.Status)
	}
	var loginData struct {
		AccessJwt string `json:"accessJwt"`
		Did       string `json:"did"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
		return merrors.New(merrors.BskyAuthDecodeErrorCode, err.Error())
	}

	BskyClient.JWT = loginData.AccessJwt
	BskyClient.DID = loginData.Did
	return nil
}
