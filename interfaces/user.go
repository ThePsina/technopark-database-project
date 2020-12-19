package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"tech-db-project/application"
	"tech-db-project/domain/entity"
	"tech-db-project/infrasctructure/tools"
)

type UserHandler struct {
	userApp *application.UserApp
}

func NewUserHandler(userApp *application.UserApp) *UserHandler {
	return &UserHandler{userApp}
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
			res, err := users.MarshalJSON()
			w.Header().Add("Content-Type", "application/json")
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res, err := resp.MarshalJSON()
	tools.HandleError(err)
	w.Write(res)
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	resp, err := entity.GetUserFromBody(r.Body)
	tools.HandleError(err)
	vars := mux.Vars(r)
	resp.Nickname = vars["nickname"]

	if err := uh.userApp.GetUser(resp); err != nil {
		switch err {
		case tools.UserNotExist:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			res, err := tools.Message{Message: "user not found"}.MarshalJSON()
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := resp.MarshalJSON()
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
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			res, err := tools.Message{Message: "user doesn't exist"}.MarshalJSON()
			tools.HandleError(err)
			w.Write(res)
			return
		case tools.UserNotUpdated:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			res, err := tools.Message{Message: "conflict while updating"}.MarshalJSON()
			tools.HandleError(err)
			w.Write(res)
			return
		default:
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := u.MarshalJSON()
	tools.HandleError(err)
	w.Write(res)
}

func (uh *UserHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := uh.userApp.DeleteAll()
	tools.HandleError(err)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&tools.Message{Message: "all info deleted"})
	tools.HandleError(err)
	w.Write(res)
}

func (uh *UserHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	s := &entity.Status{}
	err := uh.userApp.GetStatus(s)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&s)
	tools.HandleError(err)
	w.Write(res)
}
