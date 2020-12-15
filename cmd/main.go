package main

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"net/http"
	"tech-db-project/application"
	"tech-db-project/infrasctructure/persistence"
	"tech-db-project/interfaces"
	"time"
)

func main() {
	router := mux.NewRouter()

	dbConf := pgx.ConnConfig{
		User:                 "farcoad",
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

	interfaces.NewUserHandler(router, uUC)
	interfaces.NewForumHandler(router, fUC)
	interfaces.NewThreadHandler(router, thUC)
	interfaces.NewPostHandler(router, pUC, uUC, thUC, fUC)

	server := &http.Server{
		Addr:         "127.0.0.1:5000",
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logrus.Fatal(server.ListenAndServe())
}
