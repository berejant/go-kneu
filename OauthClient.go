package kneu

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type OauthClient struct {
	clientId     int
	clientSecret string
}

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	UserId      int    `json:"user_id"`
}

type oauthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (client *OauthClient) GetOauthUrl(redirectUri string, state string) string {
	return AuthBaseUri +
		"/oauth?response_type=code" +
		"&client_id=" + strconv.Itoa(client.clientId) +
		"&redirect_uri=" + url.QueryEscape(redirectUri) +
		"&state=" + url.QueryEscape(state)
}

func (client *OauthClient) GetOauthToken(redirectUri string, code string) (tokenResponse OauthTokenResponse, err error) {
	var response *http.Response

	postData := "client_id=" + strconv.Itoa(client.clientId) +
		"&client_secret=" + url.QueryEscape(client.clientSecret) +
		"&code=" + url.QueryEscape(code) +
		"&redirect_uri=" + url.QueryEscape(redirectUri) +
		"&grant_type=authorization_code"

	if err == nil {
		response, err = http.Post(
			AuthBaseUri+"/oauth/token",
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
