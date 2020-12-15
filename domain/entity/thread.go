package entity

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

//easyjson:json
type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created,omitempty"`
	Forum   string `json:"forum"`
	Id      int64  `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title"`
	Votes   int64  `json:"votes"`
}

//easyjson:json
type Threads []Thread

func GetThreadFromBody(body io.ReadCloser) (*Thread, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return &Thread{}, err
	}
	defer body.Close()

	var f *Thread
	err = json.Unmarshal(data, &f)
	if err != nil {
		return &Thread{}, err
	}
	return f, nil
}
