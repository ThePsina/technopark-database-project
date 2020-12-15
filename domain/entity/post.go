package entity

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"io"
	"io/ioutil"
)

type Post struct {
	Author   string           `json:"author"`
	Created  string           `json:"created"`
	Forum    string           `json:"forum" url:"param"`
	Id       int64            `json:"id"`
	IsEdited bool             `json:"isEdited"`
	Message  string           `json:"message"`
	Parent   int64            `json:"parent"`
	Thread   int64            `json:"thread"`
	Path     pgtype.Int8Array `json:"-"`
}

func GetPostFromBody(body io.ReadCloser) (Post, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return Post{}, err
	}
	defer body.Close()

	var f Post
	err = json.Unmarshal(data, &f)
	if err != nil {
		return Post{}, err
	}
	return f, nil
}
