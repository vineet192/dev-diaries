package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type InfiscalAccessToken struct {
	Token             string `json:"accessToken"`
	ExpiresIn         uint   `json:"expiresIn"`
	AccessTokenMaxTTL uint   `json:"accessTokenMaxTTL"`
	TokenType         string `json:"tokenType"`
}

type SecretResponse struct {
	Secret struct {
		ID            string `json:"_id"`
		Version       int    `json:"version"`
		Workspace     string `json:"workspace"`
		Type          string `json:"type"`
		Environment   string `json:"environment"`
		SecretKey     string `json:"secretKey"`
		SecretValue   string `json:"secretValue"`
		SecretComment string `json:"secretComment"`
	} `json:"secret"`
}

func (t *InfiscalAccessToken) IsExpired() bool {
	_, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	return errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet)
}

var AccessToken InfiscalAccessToken

var InfiscalAuthEndpoint string
var InfiscalSecretsEndpoint string

func refreshAccessToken() (err error) {

	//Build up the request with payload and headers
	payload := url.Values{}
	payload.Add("clientSecret", os.Getenv("INFISCAL_SECRET"))
	payload.Add("clientId", os.Getenv("INFISCAL_CLIENT_ID"))

	req, err := http.NewRequest("POST", InfiscalAuthEndpoint, strings.NewReader(payload.Encode()))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//Perform the POST request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected error with secrets vault status code is %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	json.Unmarshal(raw, &AccessToken)

	return
}

func getInfiscalSecret(key string) (string, error) {
	url, err := url.JoinPath(InfiscalSecretsEndpoint, key)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", url, nil)

	query := req.URL.Query()
	query.Add("environment", os.Getenv("INFISCAL_ENVIRONMENT"))
	query.Add("workspaceId", os.Getenv("INFISCAL_WORKSPACE_ID"))
	req.URL.RawQuery = query.Encode()

	req.Header.Add("Authorization", "Bearer "+AccessToken.Token)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var data SecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.Secret.SecretValue, nil
}

func GetSecret(key string) (value string, err error) {

	InfiscalAuthEndpoint, err = url.JoinPath(os.Getenv("INFISCAL_HOST"), "/api/v1/auth/universal-auth/login")

	if err != nil {
		return "", err
	}

	InfiscalSecretsEndpoint, err = url.JoinPath(os.Getenv("INFISCAL_HOST"), "/api/v3/secrets/raw/")

	if err != nil {
		return "", err
	}

	if AccessToken == (InfiscalAccessToken{}) {
		refreshAccessToken()
	}

	value, err = getInfiscalSecret(key)

	if err != nil {
		return "", err
	}

	return value, nil
}
