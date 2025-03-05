package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zemzale/ubiquitest/oapi"
)

var _ oapi.StrictServerInterface = (*Router)(nil)

type Router struct{}

func Run() error {
	mux := chi.NewRouter()
	r := &Router{}

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
	return nil, errors.New("not implemented")
}
