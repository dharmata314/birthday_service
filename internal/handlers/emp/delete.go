package handlers

import (
	"birthday-service/api/response"
	errMsg "birthday-service/internal/err"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func DeleteEmpHandler(log *slog.Logger, empRepo Employee) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.delete.employee"
		log := log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		log.Info("Extracted ID from URL", slog.String("id", idStr))
		if idStr == "" {
			log.Error("ID parameter is missing in the URL")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Missing emp ID"))
			return
		}
		if err != nil {
			log.Error("Invalid employee ID", errMsg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid emp ID"))
			return
		}

		err = empRepo.DeleteEmpById(r.Context(), id)
		if err != nil {
			log.Error("Failed to delete employee", errMsg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to delete employee"))
			return
		}
		log.Info("employee deleted")
		render.Status(r, http.StatusNoContent)
		render.JSON(w, r, response.OK())
	}
}
