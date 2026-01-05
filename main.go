package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"workout_ledger/database"
)

// @title           Workout Ledger API
// @version         1.0
// @description     Simple REST API for Workout Ledger.
// @BasePath        /
func main() {
	registerRoutes()
	initDatabase()
	addr := ":8080"
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		//TIP <p>To start your debugging session, right-click your code in the editor and select the Debug option.</p> <p>We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>.</p>
		fmt.Println("i =", 100/i)
	}

	srv := &http.Server{Addr: addr, Handler: nil}

	// Graceful shutdown
	shutdownErrCh := make(chan error, 1)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Printf("shutdown signal received; shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			shutdownErrCh <- err
			return
		}

		if err := database.Close(); err != nil {
			shutdownErrCh <- err
			return
		}

		shutdownErrCh <- nil
	}()

	log.Printf("starting HTTP server on %s (GET /hello, GET /health, /swagger/index.html)", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	if err := <-shutdownErrCh; err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func initDatabase() {
	ctx := context.Background()
	postgres, err := database.Open(ctx, database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	database.SetDB(postgres)
}
