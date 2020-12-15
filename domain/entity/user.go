package entity

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type User struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
	Nickname string `json:"nickname"`
}

func GetUserFromBody(body io.ReadCloser) (*User, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return &User{}, err
	}
	defer body.Close()

	var f *User
	err = json.Unmarshal(data, &f)
	if err != nil {
		return &User{}, err
	}
	return f, nil
}
