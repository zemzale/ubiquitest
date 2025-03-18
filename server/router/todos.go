package router

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/oapi"
)

func (r *Router) PostTodos(
	ctx context.Context, request oapi.PostTodosRequestObject,
) (oapi.PostTodosResponseObject, error) {
	// TODO: Move this to domain
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
