package router

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/oapi"
)

func (r *Router) PostTasks(
	ctx context.Context, request oapi.PostTasksRequestObject,
) (oapi.PostTasksResponseObject, error) {
	parnetID := uuid.Nil
	if request.Body.ParentId != nil {
		parnetID = *request.Body.ParentId
	}

	err := r.tasksStore.Run(tasks.Task{
		ID:        request.Body.Id,
		Title:     request.Body.Title,
		CreatedBy: request.Body.CreatedBy,
		Completed: false,
		ParentID:  parnetID,
		Cost:      lo.FromPtr(request.Body.Cost),
	})
	if err != nil {
		return oapi.PostTasks500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	return oapi.PostTasks201Response{}, nil
}

func (r *Router) GetTasks(
	ctx context.Context, request oapi.GetTasksRequestObject,
) (oapi.GetTasksResponseObject, error) {
	taskList, err := r.taskList.Run()
	if err != nil {
		return oapi.GetTasks500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	return oapi.GetTasks200JSONResponse(
		lo.Map(r.tasksCalculate.Run(taskList), func(t tasks.Task, _ int) oapi.Todo {
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
				Cost: lo.ToPtr(t.Cost),
			}
		}),
	), nil
}
