package interfaces

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"net/http"
	"tech-db-project/application"
	"tech-db-project/infrasctructure/persistence"
)

func RegisterHandlers() http.Handler {
	router := mux.NewRouter()

	dbConf := pgx.ConnConfig{
		User:                 "thepsina",
		Database:             "forum",
		Password:             "postgres",
		PreferSimpleProtocol: false,
	}

	dbPoolConf := pgx.ConnPoolConfig{
		ConnConfig:     dbConf,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	dbConn, err := pgx.NewConnPool(dbPoolConf)
	if err != nil {
		logrus.Fatal(err)
	}

	uRep := persistence.NewUserRepo(dbConn)
	fRep := persistence.NewForumDB(dbConn)
	thRep := persistence.NewThreadDB(dbConn)
	pRep := persistence.NewPostDB(dbConn)

	err = uRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}
	err = fRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}
	err = thRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}
	err = pRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}

	uUC := application.NewUserApp(uRep)
	fUC := application.NewForumApp(fRep, uRep)
	thUC := application.NewThreadApp(thRep ,fRep)
	pUC := application.NewPostApp(pRep, thRep)

	uh := NewUserHandler(uUC)
	fh := NewForumHandler(fUC)
	thH := NewThreadHandler(thUC)
	ph := NewPostHandler(pUC, uUC, thUC, fUC)

	router.HandleFunc("/api/forum/create", fh.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", fh.GetForumInfo).Methods(http.MethodGet)

	router.HandleFunc("/api/forum/{slug}/users", fh.GetForumUsers).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/threads", fh.GetForumThreads).Methods(http.MethodGet)

	router.HandleFunc("/api/thread/{slug}/create", ph.CreatePosts).Methods(http.MethodPost)
	router.HandleFunc("/api/post/{id}/details", ph.GetPost).Methods(http.MethodGet)
	router.HandleFunc("/api/post/{id}/details", ph.UpdatePost).Methods(http.MethodPost)

	router.HandleFunc("/api/forum/{forum}/create", thH.CreateThread).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug}/details", thH.GetThreadInfo).Methods(http.MethodGet)
	router.HandleFunc("/api/thread/{slug}/details", thH.UpdateThread).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug}/vote", thH.CreateVote).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug}/posts", thH.GetThreadPosts).Methods(http.MethodGet)

	router.HandleFunc("/api/service/status", uh.GetStatus).Methods(http.MethodGet)
	router.HandleFunc("/api/service/clear", uh.DeleteAll).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", uh.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/api/user/{nickname}/profile", uh.UpdateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/create", uh.AddUser).Methods(http.MethodPost)

	return router
}
