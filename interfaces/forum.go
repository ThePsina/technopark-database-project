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

func NewForumHandler(forumApp *application.ForumApp) *ForumHandler {
	return &ForumHandler{forumApp}
}

func (fh *ForumHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
	f, err := entity.GetForumFromBody(r.Body)
	if err != nil {
		tools.HandleError(err)
	}

	if err := fh.forumApp.CreateForum(f); err != nil {
		switch err {
		case tools.UserNotExist:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "User not found"})
			w.Write(res)
			tools.HandleError(err)
			return
		case tools.ForumExist:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			res, err := f.MarshalJSON()
			w.Write(res)
			tools.HandleError(err)
			return
		default:
			logrus.Error(err)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res, err := json.Marshal(&f)
	w.Write(res)
	tools.HandleError(err)
}

func (fh *ForumHandler) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	f, err := entity.GetForumFromBody(r.Body)
	tools.HandleError(err)

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	if err := fh.forumApp.GetForum(f); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "User not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&f)
	w.Write(res)
	tools.HandleError(err)
}

func (fh *ForumHandler) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	f, err := entity.GetForumFromBody(r.Body)
	tools.HandleError(err)

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	ths, err := fh.forumApp.GetForumThreads(f, r.URL.Query().Get("desc"), r.URL.Query().Get("limit"), r.URL.Query().Get("since"))
	if err != nil {
		switch err {
		case tools.ForumNotExist:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			res, err := tools.Message{Message: "forum not found"}.MarshalJSON()
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			tools.HandleError(err)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := ths.MarshalJSON()
	tools.HandleError(err)
	w.Write(res)
}

func (fh *ForumHandler) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	f := &entity.Forum{}

	vars := mux.Vars(r)
	f.Slug = vars["slug"]

	users, err := fh.forumApp.GetForumUsers(f, r.FormValue("desc"), r.FormValue("limit"), r.FormValue("since"))
	if err != nil {
		switch err {
		case tools.ForumNotExist:
			w.Header().Add("Content-Type", "application/json")
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&users)
	w.Write(res)
	tools.HandleError(err)
}
