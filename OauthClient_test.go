package kneu

import (
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOauthClient_GetOauthUrl(t *testing.T) {
	client := OauthClient{
		ClientId: 100,
	}

	expectedUrl := AuthBaseUri + "/oauth?response_type=code&client_id=100&redirect_uri=http%3A%2F%2Fself%2Fredirect.html&state=state88"
	redirectUrl := client.GetOauthUrl("http://self/redirect.html", "state88")

	assert.Equal(t, expectedUrl, redirectUrl)
}

func TestOauthClient_GetOauthToken(t *testing.T) {
	clientId := 100
	clientSecret := "test_secret"
	redirectUri := "http://localhost/redirect_uri"

	expectedPost := "client_id=100&client_secret=test_secret&code=test_code&redirect_uri=http%3A%2F%2Flocalhost%2Fredirect_uri&grant_type=authorization_code"

	defer gock.Off()

	t.Run("success", func(t *testing.T) {
		gock.New(AuthBaseUri).
			Post("/oauth/token").
			MatchType("url").
			BodyString(expectedPost).
			Reply(200).
			JSON(`{
				"access_token": "00000eb1ed29a47b4c38f9700d49AA00",
				"token_type":   "Bearer",
				"expires_in":   7200,
				"user_id":      999
			}`)

		client := OauthClient{
			ClientId:     clientId,
			ClientSecret: clientSecret,
		}

		tokenResponse, err := client.GetOauthToken(redirectUri, "test_code")

		assert.NoError(t, err)
		assert.Equal(t, "00000eb1ed29a47b4c38f9700d49AA00", tokenResponse.AccessToken)
		assert.Equal(t, "Bearer", tokenResponse.TokenType)
		assert.Equal(t, 7200, tokenResponse.ExpiresIn)
		assert.Equal(t, 999, tokenResponse.UserId)
	})

	t.Run("http_error", func(t *testing.T) {
		gock.New(AuthBaseUri).
			Post("/oauth/token").
			MatchType("url").
			BodyString(expectedPost).
			Reply(500)

		client := OauthClient{
			ClientId:     clientId,
			ClientSecret: clientSecret,
		}

		_, err := client.GetOauthToken(redirectUri, "test_code")

		assert.Error(t, err)
		assert.Equal(t, "Receive http code: 500", err.Error())
	})

	t.Run("api_error", func(t *testing.T) {
		gock.New(AuthBaseUri).
			Post("/oauth/token").
			MatchType("url").
			BodyString(expectedPost).
			Reply(404).
			JSON(`{
				"error": "Fake",
				"error_description": "Test error description"
			}`)

		client := OauthClient{
			ClientId:     clientId,
			ClientSecret: clientSecret,
		}

		_, err := client.GetOauthToken(redirectUri, "test_code")

		assert.Error(t, err)
		assert.Equal(t, "Fake: Test error description", err.Error())
	})
}
