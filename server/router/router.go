package router

import (
	"cmp"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/domain/users"
	"github.com/zemzale/ubiquitest/oapi"
	"github.com/zemzale/ubiquitest/storage"
	"github.com/zemzale/ubiquitest/ws"
)

var _ oapi.StrictServerInterface = (*Router)(nil)

type Router struct {
	db           *sqlx.DB
	ws           *ws.Server
	list         *tasks.List
	upsertUser   *users.FindOrCreate
	storeTask    *tasks.Store
	userFindByID *users.FindByID
}

func NewRouter(db *sqlx.DB) *Router {
	return &Router{
		db:           db,
		ws:           ws.NewServer(db),
		list:         tasks.NewList(db),
		upsertUser:   users.NewFindOrCreate(db),
		storeTask:    tasks.NewStore(db, storage.NewTaskRepository(db), storage.NewUserRepository(db)),
		userFindByID: users.NewFindById(db),
	}
}

func Run(db *sqlx.DB) error {
	mux := chi.NewRouter()
	r := NewRouter(db)

	setupRoutes(r, mux)
	printDebugRoutes(mux)

	port := cmp.Or(os.Getenv("HTTP_PORT"), ":8080")
	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}

	return nil
}

func setupRoutes(r *Router, mux *chi.Mux) {
	mux.Use(middleware.Logger)
	mux.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"https://ubiquitest.netlify.app", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler)
	oapi.HandlerFromMux(oapi.NewStrictHandler(r, nil), mux)
	mux.HandleFunc("/ws/tasks", r.WsTasks)
}

func printDebugRoutes(mux *chi.Mux) {
	for _, route := range mux.Routes() {
		fmt.Println(route.Pattern)
	}
}
