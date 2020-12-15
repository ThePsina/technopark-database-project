package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"tech-db-project/application"
	"tech-db-project/domain/entity"
	"tech-db-project/infrasctructure/tools"
)

type ForumHandler struct {
	forumApp *application.ForumApp
}

func NewForumHandler(router *mux.Router, forumApp *application.ForumApp) {
	fh := &ForumHandler{forumApp}

	router.HandleFunc("/api/forum/create", fh.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", fh.GetForumInfo).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/threads", fh.GetForumThreads).
		Queries("desc", "limit", "since").Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/users", fh.GetForumUsers).Methods(http.MethodGet)
}

func (fh *ForumHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
	f, err := entity.GetForumFromBody(r.Body)
	if err != nil {
		tools.HandleError(err)
	}

	if err := fh.forumApp.CreateForum(f); err != nil {
		switch err {
		case tools.UserNotExist:
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "User not found"})
			w.Write(res)
			tools.HandleError(err)
			return
		case tools.ForumExist:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			logrus.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	res, err := json.Marshal(&f)
	w.Write(res)
	tools.HandleError(err)
}

func (fh *ForumHandler) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	f := &entity.Forum{}

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	if err := fh.forumApp.GetForum(f); err != nil {
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "User not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&f)
	w.Write(res)
	tools.HandleError(err)
}

func (fh *ForumHandler) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	f := &entity.Forum{}

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	ths, err := fh.forumApp.GetForumThreads(f, r.FormValue("desc"), r.FormValue("limit"), r.FormValue("since"))
	if err != nil {
		switch err {
		case tools.ForumNotExist:
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "forum not found"})
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			tools.HandleError(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&ths)
	w.Write(res)
	tools.HandleError(err)
}

func (fh *ForumHandler) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	f := &entity.Forum{}

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	users, err := fh.forumApp.GetForumUsers(f, r.FormValue("desc"), r.FormValue("limit"), r.FormValue("since"))
	if err != nil {
		switch err {
		case tools.ForumNotExist:
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "forum not found"})
			tools.HandleError(err)
			w.Write(res)
			tools.HandleError(err)
			return
		default:
			tools.HandleError(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&users)
	w.Write(res)
	tools.HandleError(err)
}
