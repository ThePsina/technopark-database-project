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

type UserHandler struct {
	userApp *application.UserApp
}

func NewUserHandler(router *mux.Router, userApp *application.UserApp) *UserHandler {
	uh := &UserHandler{userApp}

	router.HandleFunc("/api/service/status", uh.GetStatus)
	router.HandleFunc("/api/service/clear", uh.DeleteAll)
	router.HandleFunc("/api/user/{nickname}/profile", uh.GetUser)
	router.HandleFunc("/api/user/{nickname}/profile", uh.UpdateUser)
	router.HandleFunc("/api/user/{nickname}/create", uh.AddUser)

	return uh
}

func (uh *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	resp, err := entity.GetUserFromBody(r.Body)
	tools.HandleError(err)

	vars := mux.Vars(r)
	resp.Nickname = vars["nickname"]

	if users, err := uh.userApp.CreateUser(resp); err != nil {
		switch err {
		case tools.UserExist:
			w.WriteHeader(http.StatusConflict)
			res, err := json.Marshal(&users)
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			logrus.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	res, err := json.Marshal(&resp)
	tools.HandleError(err)
	w.Write(res)
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	resp := &entity.User{}
	vars := mux.Vars(r)
	resp.Nickname = vars["nickname"]

	if err := uh.userApp.GetUser(resp); err != nil {
		switch err {
		case tools.UserNotExist:
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "user not found"})
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			logrus.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&resp)
	tools.HandleError(err)
	w.Write(res)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	u, err := entity.GetUserFromBody(r.Body)
	tools.HandleError(err)

	vars := mux.Vars(r)
	u.Nickname = vars["nickname"]

	if err := uh.userApp.UpdateUser(u); err != nil {
		switch err {
		case tools.UserNotExist:
			w.WriteHeader(http.StatusNotFound)
			res, err := json.Marshal(&tools.Message{Message: "user doesn't exist"})
			tools.HandleError(err)
			w.Write(res)
			return
		case tools.UserNotUpdated:
			w.WriteHeader(http.StatusConflict)
			res, err := json.Marshal(&tools.Message{Message: "conflict while updating"})
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			logrus.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&u)
	tools.HandleError(err)
	w.Write(res)
}

func (uh *UserHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := uh.userApp.DeleteAll()
	tools.HandleError(err)

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&tools.Message{Message: "all info deleted"})
	tools.HandleError(err)
	w.Write(res)
}

func (uh *UserHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	s := &entity.Status{}
	err := uh.userApp.GetStatus(s)

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&s)
	tools.HandleError(err)
	w.Write(res)
}
