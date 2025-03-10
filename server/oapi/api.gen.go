// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package oapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Error defines model for Error.
type Error struct {
	// Error The error message
	Error *string `json:"error,omitempty"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	// Id The ID of the user
	Id uint `json:"id"`

	// Username The username of the logged in user
	Username string `json:"username"`
}

// Todo defines model for Todo.
type Todo struct {
	// CreatedBy The user id of the user who create the todo item
	CreatedBy uint `json:"created_by"`

	// Id The ID of the todo item
	Id openapi_types.UUID `json:"id"`

	// Title The title of the todo item
	Title string `json:"title"`
}

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	// Username The username to login with
	Username string `json:"username"`
}

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// PostTodosJSONRequestBody defines body for PostTodos for application/json ContentType.
type PostTodosJSONRequestBody = Todo

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Login the user with the given username
	// (POST /login)
	PostLogin(w http.ResponseWriter, r *http.Request)
	// Get all todo items
	// (GET /todos)
	GetTodos(w http.ResponseWriter, r *http.Request)
	// Create a new todo item
	// (POST /todos)
	PostTodos(w http.ResponseWriter, r *http.Request)
	// Get user by id
	// (GET /user/{id})
	GetUserId(w http.ResponseWriter, r *http.Request, id uint)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Login the user with the given username
// (POST /login)
func (_ Unimplemented) PostLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get all todo items
// (GET /todos)
func (_ Unimplemented) GetTodos(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create a new todo item
// (POST /todos)
func (_ Unimplemented) PostTodos(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get user by id
// (GET /user/{id})
func (_ Unimplemented) GetUserId(w http.ResponseWriter, r *http.Request, id uint) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostLogin operation middleware
func (siw *ServerInterfaceWrapper) PostLogin(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostLogin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetTodos operation middleware
func (siw *ServerInterfaceWrapper) GetTodos(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTodos(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostTodos operation middleware
func (siw *ServerInterfaceWrapper) PostTodos(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostTodos(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetUserId operation middleware
func (siw *ServerInterfaceWrapper) GetUserId(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id uint

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUserId(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/login", wrapper.PostLogin)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/todos", wrapper.GetTodos)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/todos", wrapper.PostTodos)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/user/{id}", wrapper.GetUserId)
	})

	return r
}

type PostLoginRequestObject struct {
	Body *PostLoginJSONRequestBody
}

type PostLoginResponseObject interface {
	VisitPostLoginResponse(w http.ResponseWriter) error
}

type PostLogin200JSONResponse LoginResponse

func (response PostLogin200JSONResponse) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostLogin400JSONResponse Error

func (response PostLogin400JSONResponse) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostLogin401JSONResponse Error

func (response PostLogin401JSONResponse) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type PostLogin500JSONResponse Error

func (response PostLogin500JSONResponse) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetTodosRequestObject struct {
}

type GetTodosResponseObject interface {
	VisitGetTodosResponse(w http.ResponseWriter) error
}

type GetTodos200JSONResponse []Todo

func (response GetTodos200JSONResponse) VisitGetTodosResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetTodos500JSONResponse Error

func (response GetTodos500JSONResponse) VisitGetTodosResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type PostTodosRequestObject struct {
	Body *PostTodosJSONRequestBody
}

type PostTodosResponseObject interface {
	VisitPostTodosResponse(w http.ResponseWriter) error
}

type PostTodos201Response struct {
}

func (response PostTodos201Response) VisitPostTodosResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type PostTodos400JSONResponse Error

func (response PostTodos400JSONResponse) VisitPostTodosResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostTodos500JSONResponse Error

func (response PostTodos500JSONResponse) VisitPostTodosResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetUserIdRequestObject struct {
	Id uint `json:"id"`
}

type GetUserIdResponseObject interface {
	VisitGetUserIdResponse(w http.ResponseWriter) error
}

type GetUserId200JSONResponse LoginResponse

func (response GetUserId200JSONResponse) VisitGetUserIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetUserId400JSONResponse Error

func (response GetUserId400JSONResponse) VisitGetUserIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetUserId401JSONResponse Error

func (response GetUserId401JSONResponse) VisitGetUserIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetUserId500JSONResponse Error

func (response GetUserId500JSONResponse) VisitGetUserIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Login the user with the given username
	// (POST /login)
	PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error)
	// Get all todo items
	// (GET /todos)
	GetTodos(ctx context.Context, request GetTodosRequestObject) (GetTodosResponseObject, error)
	// Create a new todo item
	// (POST /todos)
	PostTodos(ctx context.Context, request PostTodosRequestObject) (PostTodosResponseObject, error)
	// Get user by id
	// (GET /user/{id})
	GetUserId(ctx context.Context, request GetUserIdRequestObject) (GetUserIdResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// PostLogin operation middleware
func (sh *strictHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	var request PostLoginRequestObject

	var body PostLoginJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostLogin(ctx, request.(PostLoginRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostLogin")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostLoginResponseObject); ok {
		if err := validResponse.VisitPostLoginResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetTodos operation middleware
func (sh *strictHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	var request GetTodosRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetTodos(ctx, request.(GetTodosRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetTodos")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetTodosResponseObject); ok {
		if err := validResponse.VisitGetTodosResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostTodos operation middleware
func (sh *strictHandler) PostTodos(w http.ResponseWriter, r *http.Request) {
	var request PostTodosRequestObject

	var body PostTodosJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostTodos(ctx, request.(PostTodosRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostTodos")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostTodosResponseObject); ok {
		if err := validResponse.VisitPostTodosResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetUserId operation middleware
func (sh *strictHandler) GetUserId(w http.ResponseWriter, r *http.Request, id uint) {
	var request GetUserIdRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetUserId(ctx, request.(GetUserIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetUserId")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetUserIdResponseObject); ok {
		if err := validResponse.VisitGetUserIdResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
