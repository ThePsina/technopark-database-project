package entity

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

//easyjson:json
type Vote struct {
	Vote     int32  `json:"voice"`
	Nickname string `json:"nickname"`
	Thread   int64  `json:"thread,omitempty"`
}

func GetVoteFromBody(body io.ReadCloser) (*Vote, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return &Vote{}, err
	}
	defer body.Close()

	var f *Vote
	err = json.Unmarshal(data, &f)
	if err != nil {
		return &Vote{}, err
	}
	return f, nil
}
