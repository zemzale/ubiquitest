package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/domain/users"
	"github.com/zemzale/ubiquitest/oapi"
	"github.com/zemzale/ubiquitest/ws"
)

var _ oapi.StrictServerInterface = (*Router)(nil)

type Router struct {
	ws           *ws.Server
	list         *tasks.List
	upsertUser   *users.FindOrCreate
	storeTask    *tasks.Store
	userFindByID *users.FindByID

	httpPort string
	mux      *chi.Mux
}

func NewRouter(
	db *sqlx.DB,
	httpPort string,
	taskStore *tasks.Store,
	taskList *tasks.List,
	upsertUser *users.FindOrCreate,
	userFindByID *users.FindByID,
) *Router {
	return &Router{
		ws:           ws.NewServer(db, taskStore, tasks.NewUpdate(db)),
		list:         taskList,
		upsertUser:   upsertUser,
		storeTask:    taskStore,
		userFindByID: userFindByID,
		mux:          chi.NewRouter(),

		httpPort: httpPort,
	}
}

func (r *Router) Run() error {
	r.setupRoutes()
	r.printDebugRoutes()

	if err := http.ListenAndServe(r.httpPort, r.mux); err != nil {
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
