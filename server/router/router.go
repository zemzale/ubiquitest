package router

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/domain/users"
	"github.com/zemzale/ubiquitest/oapi"
	"github.com/zemzale/ubiquitest/ws"
)

var _ oapi.StrictServerInterface = (*Router)(nil)

type Router struct {
	db         *sqlx.DB
	ws         *ws.Server
	list       *tasks.List
	upsertUser *users.FindOrCreate
}

func NewRouter(db *sqlx.DB) *Router {
	return &Router{db: db, ws: ws.NewServer(db), list: tasks.NewList(db), upsertUser: users.NewFindOrCreate(db)}
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
	query := "INSERT INTO todos (id, title, created_by, completed ) VALUES (?, ?, ?, ?)"
	args := []any{request.Body.Id, request.Body.Title, request.Body.CreatedBy, request.Body.Completed}
	if request.Body.ParentId != nil {
		query = "INSERT INTO todos (id, title, created_by, completed, parent_id ) VALUES (?, ?, ?, ?, ?)"
		args = append(args, request.Body.ParentId)
	}

	result, err := r.db.Exec(query, args...)
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
	user, err := r.upsertUser.Run(request.Body.Username)
	if err != nil {
		return oapi.PostLogin500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	return oapi.PostLogin200JSONResponse{Id: user.ID, Username: user.Username}, nil
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
			return oapi.Todo{
				Id:        t.ID,
				Title:     t.Title,
				CreatedBy: t.CreatedBy,
				Completed: t.Completed,
				ParentId: func() *uuid.UUID {
					if t.ParentID == uuid.Nil {
						return nil
					}

					return &t.ParentID
				}(),
			}
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

func (r *Router) GetUserId(ctx context.Context, request oapi.GetUserIdRequestObject) (oapi.GetUserIdResponseObject, error) {
	var user struct {
		Id       uint   `db:"id"`
		Username string `db:"username"`
	}
	err := r.db.Get(&user, "SELECT * FROM users where id=?", request.Id)
	if err != nil {
		return oapi.GetUserId500JSONResponse{Error: lo.ToPtr("Internal server error")}, nil
	}
	return oapi.GetUserId200JSONResponse{Id: user.Id, Username: user.Username}, nil
}
