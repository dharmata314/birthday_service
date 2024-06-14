package handlers

import (
	"birthday-service/api/response"
	errMsg "birthday-service/internal/err"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

func DeleteSub(log *slog.Logger, subRepo Sub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.subs.delete.New"
		log := log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid sub ID"))
			return
		}

		err = subRepo.DeleteSub(r.Context(), id)
		if err != nil {
			log.Error("Failed to delete sub", errMsg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to delete sub"))
			return
		}
		log.Info("subscriptions deleted")
		render.Status(r, http.StatusNoContent)
	}
}
