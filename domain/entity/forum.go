package entity

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

//easyjson:json
type Forum struct {
	Posts   int64  `json:"posts"`
	Slug    string `json:"slug"`
	Threads int32  `json:"threads"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

func GetForumFromBody(body io.ReadCloser) (*Forum, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return &Forum{}, err
	}
	defer body.Close()

	var f *Forum
	err = json.Unmarshal(data, &f)
	if err != nil {
		return &Forum{}, err
	}
	return f, nil
}
