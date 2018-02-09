package main

import (
	"log"
	"net/http"
	"revelforce-admin/cmd/config"
	"revelforce-admin/handlers"
	"revelforce-admin/internal/platform/db"
	"revelforce-admin/internal/platform/session"
	"time"

	"github.com/spf13/viper"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Load Config : %v", err)
	}

	masterDB, err := db.GetConnection()
	if err != nil {
		log.Fatalf("Connect DB : %v", err)
	}
	defer masterDB.Close()

	if err := masterDB.Ping(); err != nil {
		log.Fatalf("DB Ping : %v", err)
	}

	sesh := session.GetSession()

	srv := http.Server{
		Addr:           viper.GetString("addr"),
		Handler:        sesh.Use(handlers.Routes()),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatalln(srv.ListenAndServe())
}
