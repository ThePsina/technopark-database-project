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

func NewThreadHandler(router *mux.Router, threadApp *application.ThreadApp) {
	thH := &ThreadHandler{threadApp}

	router.HandleFunc("/api/forum/{forum}/create", thH.CreateThread)
	router.HandleFunc("/api/thread/{slug}/details", thH.GetThreadInfo)
	router.HandleFunc("/api/thread/{slug}/details", thH.UpdateThread)
	router.HandleFunc("/api/thread/{slug}/vote", thH.CreateVote)
	router.HandleFunc("/api/thread/{slug}/posts", thH.GetThreadPosts).Queries("desc", "sort", "limit", "since")
}

func (thH *ThreadHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	th, err := entity.GetThreadFromBody(r.Body)
	tools.HandleError(err)
	vars := mux.Vars(r)

	th.Forum = vars["forum"]
	if err := thH.threadApp.CreateThread(th); err != nil {
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
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "thread not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

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
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "user not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}
	if err == tools.ThreadNotExist {
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "thread not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

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
		w.WriteHeader(http.StatusNotFound)
		res, err := json.Marshal(&tools.Message{Message: "thread not found"})
		tools.HandleError(err)
		w.Write(res)
		return
	}

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
		r.FormValue("desc"),
		r.FormValue("sort"),
		r.FormValue("limit"),
		r.FormValue("since"),
	)

	if err != nil {
		if err == tools.ThreadNotExist {
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "thread not found"})
			tools.HandleError(err)
			w.Write(res)
			return
		}
		tools.HandleError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&posts)
	tools.HandleError(err)
	w.Write(res)
	return
}
