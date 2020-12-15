package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"tech-db-project/application"
	"tech-db-project/domain/entity"
	"tech-db-project/infrasctructure/tools"
)

type ThreadHandler struct {
	threadApp *application.ThreadApp
}

func NewThreadHandler(threadApp *application.ThreadApp) *ThreadHandler {
	return &ThreadHandler{threadApp}
}

func (thH *ThreadHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	th, err := entity.GetThreadFromBody(r.Body)
	tools.HandleError(err)
	vars := mux.Vars(r)

	th.Forum = vars["forum"]
	if err := thH.threadApp.CreateThread(th); err != nil {
		w.Header().Add("Content-Type", "application/json")
		if err == tools.ThreadExist {
			w.WriteHeader(http.StatusConflict)
			res, err := json.Marshal(&th)
			tools.HandleError(err)
			w.Write(res)
			return
		}
		if err == tools.UserNotExist {
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "user not exist"})
			tools.HandleError(err)
			w.Write(res)
			return
		}
		tools.HandleError(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res, err := json.Marshal(&th)
	tools.HandleError(err)
	w.Write(res)
}

func (thH *ThreadHandler) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	th, err := entity.GetThreadFromBody(r.Body)
	tools.HandleError(err)
	vars := mux.Vars(r)

	th.Slug = vars["slug"]
	if err := thH.threadApp.GetThreadInfo(th); err == tools.ThreadNotExist {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "thread not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&th)
	tools.HandleError(err)
	w.Write(res)
}

func (thH *ThreadHandler) CreateVote(w http.ResponseWriter, r *http.Request) {
	th := &entity.Thread{}
	vote, err := entity.GetVoteFromBody(r.Body)
	tools.HandleError(err)
	vars := mux.Vars(r)

	th.Slug = vars["slug"]
	err = thH.threadApp.CreateVote(th, vote)
	if err == tools.UserNotExist {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "user not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}
	if err == tools.ThreadNotExist {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "thread not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&th)
	tools.HandleError(err)
	w.Write(res)
	return
}

func (thH *ThreadHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	th, err := entity.GetThreadFromBody(r.Body)
	tools.HandleError(err)

	vars := mux.Vars(r)
	th.Slug = vars["slug"]

	err = thH.threadApp.UpdateThread(th)
	if err == tools.ThreadNotExist {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "thread not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&th)
	tools.HandleError(err)
	w.Write(res)
	return
}

func (thH *ThreadHandler) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	th, err := entity.GetThreadFromBody(r.Body)
	tools.HandleError(err)

	vars := mux.Vars(r)
	th.Slug = vars["slug"]

	posts, err := thH.threadApp.GetThreadPosts(
		th,
		r.URL.Query().Get("desc"),
		r.URL.Query().Get("sort"),
		r.URL.Query().Get("limit"),
		r.URL.Query().Get("since"),
	)

	if err != nil {
		if err == tools.ThreadNotExist {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "thread not found"})
			tools.HandleError(err)
			w.Write(res)
			return
		}
		tools.HandleError(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&posts)
	tools.HandleError(err)
	w.Write(res)
	return
}
