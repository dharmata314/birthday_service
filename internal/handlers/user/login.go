package handlers

import (
	"birthday-service/api/response"
	"birthday-service/internal/auth"
	errMsg "birthday-service/internal/err"
	"birthday-service/jwt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type ResponseAuthUser struct {
	response.Response
	ID    int    `json:"user_id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func LoginFunc(log *slog.Logger, userRepository User, jwt *jwt.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestUser
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", errMsg.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		user, err := userRepository.FindUserByEmail(r.Context(), req.Email)
		if err != nil {
			log.Error("User not found with email")
			render.JSON(w, r, response.Error("Invalid email"))
			return
		}

		errAuth := auth.ComparePasswordHash(req.Password, user.Password)
		if errAuth != nil {
			log.Error("Invalid password")
			render.JSON(w, r, response.Error("Invalid password"))
			return
		}
		token, err := jwt.GenerateToken(user.Email, time.Second*600)
		if err != nil {
			log.Error("failed to authoriza")
			return
		}

		log.Info("User authenticated")
		responseAuthOK(w, r, req.Email, user.ID, token)
	}
}

func responseAuthOK(w http.ResponseWriter, r *http.Request, email string, userID int, token string) {
	render.JSON(w, r, ResponseAuthUser{Response: response.OK(),
		Email: email, ID: userID, Token: token})
}
