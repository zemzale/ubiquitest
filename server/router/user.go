package router

import (
	"context"

	"github.com/samber/lo"
	"github.com/zemzale/ubiquitest/oapi"
)

func (r *Router) PostLogin(
	ctx context.Context, request oapi.PostLoginRequestObject,
) (oapi.PostLoginResponseObject, error) {
	user, err := r.upsertUser.Run(request.Body.Username)
	if err != nil {
		return oapi.PostLogin500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	return oapi.PostLogin200JSONResponse{Id: user.ID, Username: user.Username}, nil
}

func (r *Router) GetUserId(ctx context.Context, request oapi.GetUserIdRequestObject) (oapi.GetUserIdResponseObject, error) {
	// TODO: Move this do domain
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
