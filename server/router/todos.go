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
	parnetID := uuid.Nil
	if request.Body.ParentId != nil {
		parnetID = *request.Body.ParentId
	}

	err := r.storeTask.Run(tasks.Task{
		ID:        request.Body.Id,
		Title:     request.Body.Title,
		CreatedBy: request.Body.CreatedBy,
		Completed: false,
		ParentID:  parnetID,
	})
	if err != nil {
		return oapi.PostTodos500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
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
