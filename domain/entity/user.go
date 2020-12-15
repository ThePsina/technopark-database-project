package entity

import (
	"io"
	"io/ioutil"
)

//easyjson:json
type User struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}

//easyjson:json
type Users []User

func GetUserFromBody(body io.ReadCloser) (*User, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return &User{}, err
	}
	defer body.Close()

	var f User
	err = f.UnmarshalJSON(data)
	if err != nil {
		return &User{}, err
	}
	return &f, nil
}
