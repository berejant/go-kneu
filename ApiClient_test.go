package kneu

import (
	"fmt"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiClient_GetUserMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedUserMe := UserMeResponse{
			Id:           999,
			Email:        "test@example.com",
			Name:         "Петренко Петр Петрович",
			LastName:     "Петренко",
			FirstName:    "Петр",
			MiddleName:   "Петрович",
			Type:         "student",
			StudentId:    123,
			GroupId:      50,
			Sex:          "male",
			TeacherId:    0,
			DepartmentId: 0,
		}

		client := ApiClient{
			accessToken: "test-access-token",
		}

		gock.New(AuthBaseUri).
			Get("/api/user/me").
			MatchHeader("Authorization", "Bearer test-access-token").
			Reply(200).
			JSON(`{
			"id": 999,
			"email": "test@example.com",
			"name": "Петренко Петр Петрович",
			"last_name": "Петренко",
			"first_name": "Петр",
			"middle_name": "Петрович",
			"type": "student",
			"student_id": 123,
			"group_id": 50,
			"sex": "male"
		}`)

		userMe, err := client.GetUserMe()

		assert.NoError(t, err)
		assert.Equal(t, expectedUserMe, userMe)
	})

	t.Run("http_error", func(t *testing.T) {
		client := ApiClient{
			accessToken: "test-access-token",
		}

		gock.New(AuthBaseUri).
			Get("/api/user/me").
			MatchHeader("Authorization", "Bearer test-access-token").
			Reply(500)

		userMe, err := client.GetUserMe()
		fmt.Println(err.Error())
		assert.Error(t, err)
		assert.Equal(t, "Receive http code: 500", err.Error())
		assert.Empty(t, userMe)
	})

	t.Run("api_error", func(t *testing.T) {
		client := ApiClient{
			accessToken: "test-access-token",
		}

		gock.New(AuthBaseUri).
			Get("/api/user/me").
			MatchHeader("Authorization", "Bearer test-access-token").
			Reply(402).
			JSON(`{
				"error": "Test error description"
			}`)

		userMe, err := client.GetUserMe()
		fmt.Println(err.Error())
		assert.Error(t, err)
		assert.Equal(t, "API error: Test error description", err.Error())
		assert.Empty(t, userMe)
	})

}
