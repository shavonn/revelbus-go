package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"revelbus/cmd/web"
	"revelbus/pkg/database"
	"revelbus/pkg/sessions"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	err := web.LoadConfig()
	if err != nil {
		log.Fatalf("Load Config : %v", err)
	}

	log.Println("main : Started : Initialize MySql")
	masterDB, err := database.GetConnection()
	if err != nil {
		log.Fatalf("startup : Connect to DB : %v", err)
	}
	defer masterDB.Close()

	if err := masterDB.Ping(); err != nil {
		log.Fatalf("DB Ping : %v", err)
	}

	sesh := sessions.GetSession()

	srv := http.Server{
		Addr:           viper.GetString("addr"),
		Handler:        sesh.Use(web.Routes()),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("startup : Listening %s", viper.GetString("addr"))
		serverErrors <- srv.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-osSignals:
		log.Println("main : Start shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v : %v", (5 * time.Second), err)
			if err := srv.Close(); err != nil {
				log.Fatalf("Could not stop http server: %v", err)
			}
		}
	}

	log.Println("main : Completed")
}
