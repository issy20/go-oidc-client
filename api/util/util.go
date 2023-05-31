package util

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type Secrets struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func ReadJson() (*Secrets, error) {
	var secrets = &Secrets{}
	data, err := os.ReadFile("env.json")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	err = json.Unmarshal(data, secrets)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return secrets, nil
}

func AddBearer(accessToken string) string {
	authorizationHeader := "Bearer " + accessToken
	return authorizationHeader
}

func CreateURL(baseURL string, queryStr string, params string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("URL parsing error: %w", err)
	}
	query := u.Query()
	query.Set(queryStr, params)

	u.RawQuery = query.Encode()

	return u.String(), err
}
