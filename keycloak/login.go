package keycloak

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

const (
	tokenEndpoint   = "%s/auth/realms/%s/protocol/openid-connect/token"
	formContentType = "application/x-www-form-urlencoded"
	loginBody       = "grant_type=client_credentials"
)

// Attempt to login to Keycloak with the provided information.
func (c *KeycloakClient) Login() error {
	url := fmt.Sprintf(tokenEndpoint, c.url, c.realm)

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(loginBody))
	req.Header.Set("Authorization", createBasicAuthorizationHeader(c.id, c.secret))
	req.Header.Set("Content-Type", formContentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("Keycloak login failed: %s (%d)", string(body), resp.StatusCode)
	}

	var t tokenResponse
	err = json.Unmarshal(body, &t)
	if err != nil {
		return err
	}

	c.token = t.AccessToken
	return nil
}

func createBasicAuthorizationHeader(id string, secret string) string {
	input := fmt.Sprintf("%s:%s", id, secret)
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return fmt.Sprintf("Basic %s", encoded)
}
