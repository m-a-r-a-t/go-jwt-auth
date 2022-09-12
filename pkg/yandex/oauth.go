package yandex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type YandexOAuthResponseData struct {
	Id           string   `json:"id,omitempty"`
	Login        string   `json:"login,omitempty"`
	ClientId     string   `json:"client_id,omitempty"`
	DisplayName  string   `json:"display_name,omitempty"`
	RealName     string   `json:"real_name,omitempty"`
	FirstName    string   `json:"first_name,omitempty"`
	LastName     string   `json:"last_name,omitempty"`
	Sex          string   `json:"sex,omitempty"`
	DefaultEmail string   `json:"default_email,omitempty"`
	Emails       []string `json:"emails,omitempty"`
	Birthday     string   `json:"birthday,omitempty"`
	Psuid        string   `json:"psuid,omitempty"`
}

/*
Функция получения данных с яндекса апи по OAuth токену
*/
func GetUserFromYandexByOauthToken(oauth_token string) (*YandexOAuthResponseData, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://login.yandex.ru/info", nil)
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", oauth_token))

	// fmt.Println(req.URL.String())

	resp, _ := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("yandex error", err)
	}

	data := YandexOAuthResponseData{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &data, nil
}
