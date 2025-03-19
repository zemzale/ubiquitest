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

	mux *chi.Mux
}

func NewRouter(db *sqlx.DB) *Router {
	taskRepo := storage.NewTaskRepository(db)

	return &Router{
		db:           db,
		ws:           ws.NewServer(db),
		list:         tasks.NewList(db, taskRepo),
		upsertUser:   users.NewFindOrCreate(db),
		storeTask:    tasks.NewStore(taskRepo, storage.NewUserRepository(db)),
		userFindByID: users.NewFindById(db),
		mux:          chi.NewRouter(),
	}
}

func (r *Router) Run() error {
	r.setupRoutes()
	r.printDebugRoutes()
	port := cmp.Or(os.Getenv("HTTP_PORT"), ":8080")
	if err := http.ListenAndServe(port, r.mux); err != nil {
		return err
	}

	return nil
}

func (r *Router) setupRoutes() {
	r.mux.Use(middleware.Logger)
	r.mux.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"https://ubiquitest.netlify.app", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler)
	oapi.HandlerFromMux(oapi.NewStrictHandler(r, nil), r.mux)
	r.mux.HandleFunc("/ws/tasks", r.WsTasks)
}

func (r *Router) printDebugRoutes() {
	for _, route := range r.mux.Routes() {
		fmt.Println(route.Pattern)
	}
}
