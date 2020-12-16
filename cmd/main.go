package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"tech-db-project/interfaces"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":5000",
		Handler:      interfaces.RegisterHandlers(),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logrus.Fatal(server.ListenAndServe())
}
