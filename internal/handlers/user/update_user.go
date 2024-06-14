package handlers

import (
	"birthday-service/api/response"
	"birthday-service/internal/auth"
	errMsg "birthday-service/internal/err"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type RequestUpdateUser struct {
	ID       int    `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewUpdateUserHandler(userRepo User, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			logger.Error("Invalid user ID")
			render.JSON(w, r, response.Error("Invalid user ID"))
			return
		}

		var req RequestUpdateUser

		err = json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			logger.Error("Invalid request", errMsg.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		user, err := userRepo.FindUserById(r.Context(), userID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			logger.Error("Failed to find user")
			render.JSON(w, r, response.Error("user not found"))
			return
		}

		user.ID = req.ID
		user.Email = req.Email
		user.Password, _ = auth.HashPassword(req.Password)

		err = userRepo.UpdateUser(r.Context(), &user)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			logger.Error("Failed to update user")
			render.JSON(w, r, response.Error("Failed to update user"))
			return
		}

		render.JSON(w, r, response.OK())

	}
}
