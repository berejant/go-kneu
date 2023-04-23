package kneu

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type OauthClientInterface interface {
	GetOauthUrl(redirectUri string, state string) string
	GetOauthToken(redirectUri string, code string) (tokenResponse OauthTokenResponse, err error)
}

type OauthClient struct {
	ClientId     uint
	ClientSecret string
	BaseUri      string
}

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	UserId      uint   `json:"user_id"`
}

type oauthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (client *OauthClient) GetOauthUrl(redirectUri string, state string) string {
	if client.BaseUri == "" {
		client.BaseUri = AuthBaseUri
	}

	return client.BaseUri +
		"/oauth?response_type=code" +
		"&client_id=" + strconv.FormatUint(uint64(client.ClientId), 10) +
		"&redirect_uri=" + url.QueryEscape(redirectUri) +
		"&state=" + url.QueryEscape(state)
}

func (client *OauthClient) GetOauthToken(redirectUri string, code string) (tokenResponse OauthTokenResponse, err error) {
	var response *http.Response

	if client.BaseUri == "" {
		client.BaseUri = AuthBaseUri
	}

	postData := "client_id=" + strconv.FormatUint(uint64(client.ClientId), 10) +
		"&client_secret=" + url.QueryEscape(client.ClientSecret) +
		"&code=" + url.QueryEscape(code) +
		"&redirect_uri=" + url.QueryEscape(redirectUri) +
		"&grant_type=authorization_code"

	if err == nil {
		response, err = http.Post(
			client.BaseUri+"/oauth/token",
			"application/x-www-form-urlencoded",
			strings.NewReader(postData),
		)

		if err == nil && response.StatusCode != 200 {
			errorResponse := oauthErrorResponse{}
			err = unmarshalResponse(response, &errorResponse)
			if err == nil {
				err = errors.New(errorResponse.Error + ": " + errorResponse.ErrorDescription)
			} else {
				err = errors.New("Receive http code: " + strconv.Itoa(response.StatusCode))
			}
		}
	}

	if err == nil {
		err = unmarshalResponse(response, &tokenResponse)
	}

	return
}
