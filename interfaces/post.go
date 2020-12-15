package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"tech-db-project/application"
	"tech-db-project/domain/entity"
	"tech-db-project/infrasctructure/tools"
)

type PostHandler struct {
	postApp   *application.PostApp
	userApp   *application.UserApp
	threadApp *application.ThreadApp
	forumApp  *application.ForumApp
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func NewPostHandler(postApp *application.PostApp, userApp *application.UserApp,
	threadApp *application.ThreadApp, forumApp *application.ForumApp) *PostHandler {
	return &PostHandler{postApp, userApp, threadApp, forumApp}
}

func (ph *PostHandler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	p, _ := entity.GetPostsFromBody(r.Body)

	th := entity.Thread{}
	vars := mux.Vars(r)
	th.Slug = vars["slug"]

	err := ph.postApp.CreatePosts(p, &th)
	if err != nil {
		if err == tools.ParentNotExist {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			res, err := json.Marshal(&tools.Message{Message: "parent conflict"})
			tools.HandleError(err)
			w.Write(res)
			return
		}
		if err == tools.ThreadNotExist || err == tools.UserNotExist {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "user or thread not found"})
			tools.HandleError(err)
			w.Write(res)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		res, err := json.Marshal(&tools.Message{Message: err.Error()})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res, err := entity.ConvertToPosts(p).MarshalJSON()
	tools.HandleError(err)
	w.Write(res)
}

func (ph *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("related")
	related := strings.Split(str, ",")

	p := &entity.Post{}
	th := &entity.Thread{}
	f := &entity.Forum{}
	u := &entity.User{}

	var err error
	vars := mux.Vars(r)

	p.Id, err = strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		logrus.Error("Cannot parse id")
	}
	if err := ph.postApp.GetPost(p); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: err.Error()})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	if Find(related, "user") {
		u.Nickname = p.Author
		err := ph.userApp.GetUser(u)
		tools.HandleError(err)
	} else {
		u = nil
	}

	if Find(related, "thread") {
		th.Slug = strconv.FormatInt(p.Thread, 10)
		err := ph.threadApp.GetThreadInfo(th)
		tools.HandleError(err)
	} else {
		th = nil
	}

	if Find(related, "forum") {
		f.Slug = p.Forum
		err := ph.forumApp.GetForum(f)
		tools.HandleError(err)
	} else {
		f = nil
	}

	ans := entity.Info{Post:   p, Forum:  f, Thread: th, Author: u}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&ans)
	tools.HandleError(err)
	w.Write(res)
}

func (ph *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	p, _ := entity.GetPostFromBody(r.Body)
	vars := mux.Vars(r)
	var err error

	p.Id, err = strconv.ParseInt(vars["id"], 10, 64)
	tools.HandleError(err)
	if err := ph.postApp.UpdatePost(&p); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "post not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := p.MarshalJSON()
	tools.HandleError(err)
	w.Write(res)
}
