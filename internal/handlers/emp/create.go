package handlers

import (
	"birthday-service/api/response"
	"birthday-service/internal/entities"
	errMsg "birthday-service/internal/err"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Employee interface {
	CreateEmployee(ctx context.Context, employee *entities.Employee) error
	DeleteEmpById(ctx context.Context, id int) error
	GetAllEmp(ctx context.Context) ([]entities.Employee, error)
	GetUpcomingBirthdays(ctx context.Context) ([]entities.Employee, error)
}

type RequestEmp struct {
	Name     string     `json:"name" validate:"required"`
	Birthday CustomDate `json:"birthday" validate:"required"`
}

type ResponseEmp struct {
	response.Response
	ID       int       `json:"emp_id"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
}

type CustomDate time.Time

const customDateFormat = "02.01.2006"

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	dateStr := string(data[1 : len(data)-1])
	parsedTime, err := time.Parse(customDateFormat, dateStr)
	if err != nil {
		return err
	}
	*cd = CustomDate(parsedTime)
	return nil
}

func (cd CustomDate) ToTime() time.Time {
	return time.Time(cd)
}

func New(log *slog.Logger, empRepository Employee) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.createEmp.New"
		log = log.With(
			slog.Any("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestEmp
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
		emp := entities.Employee{Name: req.Name, Birthday: req.Birthday.ToTime()}
		err = empRepository.CreateEmployee(r.Context(), &emp)
		if err != nil {
			log.Error("Failed to create employee", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to create employee"))
			return
		}
		log.Info("employee added")
		responseOK(w, r, req.Name, req.Birthday.ToTime(), emp.ID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, name string, birthday time.Time, emp_id int) {
	render.JSON(w, r, ResponseEmp{
		response.OK(),
		emp_id,
		name,
		birthday,
	})
}
