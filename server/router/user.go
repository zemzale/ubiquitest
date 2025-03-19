package router

import (
	"context"

	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/oapi"
)

func (r *Router) PostLogin(
	ctx context.Context, request oapi.PostLoginRequestObject,
) (oapi.PostLoginResponseObject, error) {
	user, err := r.usersUpsert.Run(request.Body.Username)
	if err != nil {
		return oapi.PostLogin500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	return oapi.PostLogin200JSONResponse{Id: user.ID, Username: user.Username}, nil
}

func (r *Router) GetUserId(ctx context.Context, request oapi.GetUserIdRequestObject) (oapi.GetUserIdResponseObject, error) {
	user, err := r.usersFindByID.Run(request.Id)
	if err != nil {
		return oapi.GetUserId500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	return oapi.GetUserId200JSONResponse{Id: user.ID, Username: user.Username}, nil
}
