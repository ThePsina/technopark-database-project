package entity

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"io"
	"io/ioutil"
)

//easyjson:json
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

//easyjson:json
type Posts []Post

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

func GetPostsFromBody(body io.ReadCloser) ([]*Post, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var f []*Post
	err = json.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func ConvertToPosts(p []*Post) Posts {
	ps := make(Posts, 0, len(p))
	for _, val := range p {
		ps = append(ps, *val)
	}

	return ps
}
