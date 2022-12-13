package api

import (
	"encoding/json"
	"net/http"

	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/logger"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (u *User) ToService() *service.User {
	return &service.User{
		ID:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func UserFromService(u *service.User) User {
	return User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func UsersFromService(usrs []service.User) []User {
	res := make([]User, 0, len(usrs))

	for i := range usrs {
		res = append(res, UserFromService(&usrs[i]))
	}

	return res
}
func (p *GetUserParams) ToService() *service.ListUsersParams {
	return &service.ListUsersParams{
		Limit:  p.Limit,
		Offset: p.Offset,
	}
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request, params GetUserParams) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	srvcUsrs, err := c.srvc.ListUsers(ctx, params.ToService())
	if err != nil {
		log.Error("cannot list users", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	WriteSuccessful(ctx, w, UsersFromService(srvcUsrs))
}
func (c *Controller) PostUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	var usr User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		log.Info("wrong user data", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newID, err := c.srvc.CreateUser(ctx, usr.ToService())
	if err != nil {
		log.Error("cannot create user", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	WriteSuccessful(ctx, w, CreatedItem{Id: newID})
}
func (c *Controller) GetUserId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	log := logger.GetLogger(ctx).With(zap.Field{
		Key:    "user_id",
		Type:   zapcore.StringType,
		String: id,
	})

	srvcUsr, err := c.srvc.GetUserByID(ctx, id)
	if err != nil {
		log.Error("cannot get user", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	WriteSuccessful(ctx, w, UserFromService(&srvcUsr))
}
func (c *Controller) PutUserId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	log := logger.GetLogger(ctx).With(zap.Field{
		Key:    "user_id",
		Type:   zapcore.StringType,
		String: id,
	})

	var usr User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		log.Info("wrong user data", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	srvcUsr, err := c.srvc.UpdateUserByID(ctx, id, usr.ToService())
	if err != nil {
		log.Error("cannot update user", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	WriteSuccessful(ctx, w, UserFromService(&srvcUsr))
}
func (c *Controller) DeleteUserId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	log := logger.GetLogger(ctx).With(zap.Field{
		Key:    "user_id",
		Type:   zapcore.StringType,
		String: id,
	})

	err := c.srvc.DeleteUserByID(ctx, id)
	if err != nil {
		log.Error("cannot delete user", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	WriteSuccessful(ctx, w, CreatedItem{Id: id})
}
