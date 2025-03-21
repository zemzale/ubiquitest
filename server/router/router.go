package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/domain/users"
	"github.com/zemzale/ubiquitest/oapi"
	"github.com/zemzale/ubiquitest/ws"
	"golang.org/x/sync/errgroup"
)

var _ oapi.StrictServerInterface = (*Router)(nil)

type Router struct {
	websocketServer *ws.Server
	taskList        *tasks.List
	tasksStore      *tasks.Store
	tasksCalculate  *tasks.CalculateCost
	usersFindByID   *users.FindByID
	usersUpsert     *users.FindOrCreate

	httpPort string
	mux      *chi.Mux
}

func NewRouter(
	httpPort string,
	taskStore *tasks.Store,
	taskList *tasks.List,
	taskCalculate *tasks.CalculateCost,
	upsertUser *users.FindOrCreate,
	userFindByID *users.FindByID,
	wss *ws.Server,
) *Router {
	return &Router{
		websocketServer: wss,
		taskList:        taskList,
		tasksCalculate:  taskCalculate,
		usersUpsert:     upsertUser,
		tasksStore:      taskStore,
		usersFindByID:   userFindByID,
		mux:             chi.NewRouter(),

		httpPort: httpPort,
	}
}

func (r *Router) Run() error {
	r.setupRoutes()
	r.printDebugRoutes()

	errGroup, ctx := errgroup.WithContext(context.Background())
	errGroup.Go(func() error {
		r.websocketServer.Run(ctx)

		return nil
	})

	errGroup.Go(func() error {
		return http.ListenAndServe(r.httpPort, r.mux)
	})

	return errGroup.Wait()
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
