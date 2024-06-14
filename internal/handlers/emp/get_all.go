package handlers

import (
	"birthday-service/api/response"
	"birthday-service/internal/entities"
	errMsg "birthday-service/internal/err"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ResponseEmpList struct {
	response.Response
	Employees []entities.Employee `json:"employees"`
}

func ListAllEmployees(log *slog.Logger, empRepository Employee) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.listAllEmployees"
		log = log.With(
			slog.Any("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		employees, err := empRepository.GetAllEmp(r.Context())
		if err != nil {
			log.Error("Failed to retrieve employees", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to retrieve employees"))
			return
		}
		log.Info("employees retrieved", slog.Any("employees", employees))

		responseOKgetEmp(w, r, employees)
	}
}

func responseOKgetEmp(w http.ResponseWriter, r *http.Request, employees []entities.Employee) {
	render.JSON(w, r, ResponseEmpList{
		Response:  response.OK(),
		Employees: employees,
	})
}
