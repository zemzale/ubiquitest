package router

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/oapi"
)

var _ oapi.StrictServerInterface = (*Router)(nil)

type Router struct {
	db *sqlx.DB
}

func NewRouter(db *sqlx.DB) *Router {
	return &Router{db: db}
}

func Run(db *sqlx.DB) error {
	mux := chi.NewRouter()
	r := NewRouter(db)

	setupRoutes(r, mux)
	printDebugRoutes(mux)

	if err := http.ListenAndServe(":9999", mux); err != nil {
		return err
	}

	return nil
}

func setupRoutes(r *Router, mux *chi.Mux) {
	oapi.HandlerFromMux(oapi.NewStrictHandler(r, nil), mux)
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
