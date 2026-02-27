package bsky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func authenticate() error {
	loginPayload := map[string]string{
		"identifier": BskyClient.Handle,
		"password":   BskyClient.Password,
	}
	body, _ := json.Marshal(loginPayload)
	resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed: %s", resp.Status)
	}
	var loginData struct {
		AccessJwt string `json:"accessJwt"`
		Did       string `json:"did"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
		return fmt.Errorf("auth decode failed: %w", err)
	}

	BskyClient.JWT = loginData.AccessJwt
	BskyClient.DID = loginData.Did
	return nil
}
