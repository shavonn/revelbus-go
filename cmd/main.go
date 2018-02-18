package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"revelforce/cmd/web"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/session"
	"sync"
	"time"

	"github.com/spf13/viper"
)

func main() {
	err := web.LoadConfig()
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
		Handler:        sesh.Use(web.Routes()),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		log.Printf("startup : Listening %s", viper.GetString("addr"))
		log.Printf("shutdown : Listener closed : %v", srv.ListenAndServe())
		wg.Done()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)

	<-osSignals

	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown did not complete in %v : %v", timeout, err)

		if err := srv.Close(); err != nil {
			log.Printf("shutdown : Error killing server : %v", err)
		}
	}

	wg.Wait()
	log.Println("main : Completed")
}
