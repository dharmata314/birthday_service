package handlers

import (
	"birthday-service/api/response"
	"birthday-service/internal/auth"
	"birthday-service/internal/entities"
	errMsg "birthday-service/internal/err"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type User interface {
	CreateUser(ctx context.Context, user *entities.User) error
	FindUserByEmail(ctx context.Context, email string) (entities.User, error)
	FindUserById(ctx context.Context, id int) (entities.User, error)
	DeleteUserById(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, user *entities.User) error
}

type RequestUser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResponseUser struct {
	response.Response
	ID    int    `json:"user_id"`
	Email string `json:"email"`
}

func New(log *slog.Logger, userRepository User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.createUser.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestUser
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
		hashPass, _ := auth.HashPassword(req.Password)
		user := entities.User{Email: req.Email, Password: hashPass}
		err = userRepository.CreateUser(r.Context(), &user)
		if err != nil {
			log.Error("Failed to create user", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to create user"))
			return
		}
		log.Info("user added")
		responseOK(w, r, req.Email, user.ID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, email string, userID int) {
	render.JSON(w, r, ResponseUser{
		response.OK(),
		userID,
		email,
	})
}
