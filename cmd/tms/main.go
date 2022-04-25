package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/parinb/tms/pkg/tms"

	_ "github.com/lib/pq" //postgreSQL driver
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalln(err)
		}
	}()
	address := flag.String("address", ":8000", "http address to listen to")
	dsn := flag.String("dsn", "", "database connection string. Defaults to DSN environment variable value")
	flag.Parse()

	if *dsn == "" {
		envDSN := os.Getenv("DSN")
		dsn = &envDSN
	}

	store, err := tms.NewPosgresStore(*dsn)
	if err != nil {
		return
	}
	defer store.Close()

	h := tms.NewHTTPTaskManager(store)

	r := mux.NewRouter()
	//creates a  users
	r.Handle("/api/v1/user/create", h.CreateUserHandler()).Methods(http.MethodPost)
	//assigns a task to a  users
	r.Handle("/api/v1/user/{id:[1-9]+}/assign/task", h.CreateUserTaskHandler()).Methods(http.MethodPost)
	//gets tasks  of a  user
	r.Handle("/api/v1/user/{id:[1-9]+}/get/tasks", h.GetUserTasksHandler()).Methods(http.MethodGet)
	//gets  a  user
	r.Handle("/api/v1/user/get/{id:[1-9]+}", h.GetUserHandler()).Methods(http.MethodGet)
	//gets  all users
	r.Handle("/api/v1/users/get", h.GetUsersHandler()).Methods(http.MethodGet)
	//deletes a  certain  user
	r.Handle("/api/v1/user/delete/{id:[1-9]+}", h.DeleteUserHandler()).Methods(http.MethodDelete)
	//gets all tasks
	r.Handle("/api/v1/task/get", h.GetTasksHandler()).Methods(http.MethodGet)
	//gets all  a specific task
	r.Handle("/api/v1/task/get/{id:[1-9]+}", h.GetTaskHandler()).Methods(http.MethodGet)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// starts http server
	srv := &http.Server{
		Addr:    *address,
		Handler: r,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("unable to start server: %v", err)
		}
	}()
	log.Printf("server started on port: %s", *address)

	<-ctx.Done()
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("gracefully shutdown server err = %v", err)
		return
	}
	log.Printf("gracefully shutdown server")
}
