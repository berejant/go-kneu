package kneu

import (
	"errors"
	"net/http"
	"strconv"
)

type ApiClientInterface interface {
	GetUserMe() (response UserMeResponse, err error)
}

type ApiClient struct {
	AccessToken string
	BaseUri     string
}

type apiErrorResponse struct {
	Error string `json:"error"`
}

type UserMeResponse struct {
	Id         int    `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Type       string `json:"type"`

	StudentId int    `json:"student_id"`
	GroupId   int    `json:"group_id"`
	Sex       string `json:"sex"`

	TeacherId    int `json:"teacher_id"`
	DepartmentId int `json:"department_id"`
}

func (client *ApiClient) doRequest(requestUri string, responseInterface any) error {
	var response *http.Response

	if client.BaseUri == "" {
		client.BaseUri = AuthBaseUri
	}

	request, err := http.NewRequest(http.MethodGet, client.BaseUri+"/api/"+requestUri, nil)
	request.Header.Set("Authorization", "Bearer "+client.AccessToken)

	if err == nil {
		response, err = http.DefaultClient.Do(request)

		if err == nil && response.StatusCode != 200 {
			errorResponse := apiErrorResponse{}
			err = unmarshalResponse(response, &errorResponse)
			if err == nil {
				err = errors.New("API error: " + errorResponse.Error)
			} else {
				err = errors.New("Receive http code: " + strconv.Itoa(response.StatusCode))
			}
		}
	}

	if err == nil {
		err = unmarshalResponse(response, &responseInterface)
	}

	return err
}

func (client *ApiClient) GetUserMe() (response UserMeResponse, err error) {
	err = client.doRequest("user/me", &response)
	return
}
