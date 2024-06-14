package handlers

import (
	"birthday-service/api/response"
	"birthday-service/internal/entities"
	errMsg "birthday-service/internal/err"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Sub interface {
	CreateSub(ctx context.Context, sub *entities.Subscription) error
	DeleteSub(ctx context.Context, id int) error
	GetSubs(ctx context.Context, EmployeeID int) ([]entities.User, error)
}

type RequestSub struct {
	UserID int `json:"user_id"`
	EmpID  int `json:"emp_id"`
}

type ResponseSub struct {
	response.Response
	ID int `json:"id"`
}

func New(log *slog.Logger, subsRepository Sub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.createSub.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestSub
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", errMsg.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", errMsg.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		sub := entities.Subscription{UserID: req.UserID, EmployeeID: req.EmpID}
		err = subsRepository.CreateSub(r.Context(), &sub)
		if err != nil {
			log.Error("Failed to create subscription", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to create subscription"))
			return
		}
		log.Info("user added")
		responseOK(w, r, sub.ID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	render.JSON(w, r, ResponseSub{
		response.OK(),
		id,
	})
}
