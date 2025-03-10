package router

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/oapi"
	"github.com/zemzale/ubiquitest/ws"
)

var _ oapi.StrictServerInterface = (*Router)(nil)

type Router struct {
	db   *sqlx.DB
	ws   *ws.Server
	list *tasks.List
}

func NewRouter(db *sqlx.DB) *Router {
	return &Router{db: db, ws: ws.NewServer(), list: tasks.NewList(db)}
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
		AllowedOrigins: []string{"https://ubiquitest.netlify.app"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:          true,
	}).Handler)
	oapi.HandlerFromMux(oapi.NewStrictHandler(r, nil), mux)
	mux.HandleFunc("/ws/todos", r.WsTodos)
}

func printDebugRoutes(mux *chi.Mux) {
	for _, route := range mux.Routes() {
		fmt.Println(route.Pattern)
	}
}

func (r *Router) PostTodos(
	ctx context.Context, request oapi.PostTodosRequestObject,
) (oapi.PostTodosResponseObject, error) {
	result, err := r.db.Exec("INSERT INTO todos (id, title ) VALUES (?, ?)", request.Body.Id, request.Body.Title)
	if err != nil {
		return oapi.PostTodos500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return oapi.PostTodos400JSONResponse{Error: lo.ToPtr("item already exists")}, nil
	}

	return oapi.PostTodos201Response{}, nil
}

func (r *Router) PostLogin(
	ctx context.Context, request oapi.PostLoginRequestObject,
) (oapi.PostLoginResponseObject, error) {
	// TODO: Refactor this to be just normal
	{
		result := r.db.QueryRow("SELECT id FROM users WHERE username = ?", request.Body.Username)
		var userID uint
		if err := result.Scan(&userID); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return oapi.PostLogin500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
			}
		}
		if userID != 0 {
			return oapi.PostLogin200JSONResponse{Id: userID, Username: request.Body.Username}, nil
		}
	}

	{
		result := r.db.QueryRow("INSERT INTO users (username) VALUES (?) RETURNING id", request.Body.Username)

		var userID uint
		if err := result.Scan(&userID); err != nil {
			return oapi.PostLogin500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
		}

		return oapi.PostLogin200JSONResponse{Id: userID, Username: request.Body.Username}, nil
	}
}

func (r *Router) GetTodos(
	ctx context.Context, request oapi.GetTodosRequestObject,
) (oapi.GetTodosResponseObject, error) {
	taskList, err := r.list.Run()
	if err != nil {
		return oapi.GetTodos500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	return oapi.GetTodos200JSONResponse(
		lo.Map(taskList, func(t tasks.Task, _ int) oapi.Todo {
			return oapi.Todo{Id: t.ID, Title: t.Title, CreatedBy: t.CreatedBy}
		}),
	), nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Router) WsTodos(writer http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("user")
	if username == "" {
		log.Println("No user name provided")
		return
	}

	log.Println("User name:", username)

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	r.ws.TakeConnection(username, conn)
}
