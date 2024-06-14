package database

import (
	"birthday-service/internal/entities"
	errMsg "birthday-service/internal/err"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewEmployeeRepository(db *pgxpool.Pool, log *slog.Logger) *EmployeeRepository {
	return &EmployeeRepository{db, log}
}

func (e *EmployeeRepository) CreateEmployee(ctx context.Context, employee *entities.Employee) error {
	err := e.db.QueryRow(ctx, `INSERT INTO Employees (name, birthday) VALUES ($1, $2) RETURNING id`, employee.Name, employee.Birthday).Scan(&employee.ID)
	if err != nil {
		e.log.Error("failed to create Employee", errMsg.Err(err))
		return err
	}
	return nil
}

func (e *EmployeeRepository) FindEmployeeByName(ctx context.Context, name string) (entities.Employee, error) {
	query, err := e.db.Query(ctx, `SELECT * FROM Employees WHERE name = $1`, name)
	if err != nil {
		e.log.Error("Error querying users", errMsg.Err(err))
		return entities.Employee{}, err
	}
	row := entities.Employee{}
	defer query.Close()
	if !query.Next() {
		e.log.Error("user not found")
		return entities.Employee{}, fmt.Errorf("user not found")
	} else {
		err := query.Scan(&row.ID, &row.Name, &row.Birthday)
		if err != nil {
			e.log.Error("error scanning users", errMsg.Err(err))
			return entities.Employee{}, err
		}
	}
	return row, nil
}

func (e *EmployeeRepository) FindEmployeeById(ctx context.Context, id int) (entities.Employee, error) {
	query, err := e.db.Query(ctx, `SELECT * FROM Employees WHERE id = $1`, id)
	if err != nil {
		e.log.Error("error querying employees", errMsg.Err(err))
		return entities.Employee{}, err
	}
	defer query.Close()
	rowArray := entities.Employee{}
	if !query.Next() {
		e.log.Error("user not found")
		return entities.Employee{}, nil
	} else {
		err := query.Scan(&rowArray.ID, &rowArray.Name, &rowArray.Birthday)
		if err != nil {
			e.log.Error("error scanning employees", errMsg.Err(err))
			return entities.Employee{}, err
		}
	}
	return rowArray, nil
}

func (e *EmployeeRepository) DeleteEmpById(ctx context.Context, id int) error {
	_, err := e.db.Exec(ctx, `DELETE FROM Employees WHERE id = $1`, id)
	if err != nil {
		e.log.Error("failed to delete employee", errMsg.Err(err))
		return err
	}
	return nil
}

func (e *EmployeeRepository) GetAllEmp(ctx context.Context) ([]entities.Employee, error) {
	query, err := e.db.Query(ctx, `SELECT * FROM Employees`)
	if err != nil {
		e.log.Error("Error querying employees", errMsg.Err(err))
		return nil, err
	}
	defer query.Close()

	var employees []entities.Employee
	for query.Next() {
		var employee entities.Employee
		err := query.Scan(&employee.ID, &employee.Name, &employee.Birthday)
		if err != nil {
			e.log.Error("Error scanning employees", errMsg.Err(err))
			return nil, err
		}
		employees = append(employees, employee)
	}

	if err := query.Err(); err != nil {
		e.log.Error("Error iterating over employees", errMsg.Err(err))
		return nil, err
	}

	return employees, nil

}

func (e *EmployeeRepository) GetUpcomingBirthdays(ctx context.Context) ([]entities.Employee, error) {
	var employees []entities.Employee

	query := `SELECT id, name, birthday
			FROM Employees
			WHERE 
    (birthday::date + INTERVAL '1 year' * (EXTRACT(year FROM CURRENT_DATE) - EXTRACT(year FROM birthday::date))) 
    BETWEEN CURRENT_DATE AND (CURRENT_DATE + INTERVAL '7 day');

`

	rows, err := e.db.Query(ctx, query)
	if err != nil {
		e.log.Error("failed to get umcoming birthdays", errMsg.Err(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee entities.Employee
		if err := rows.Scan(&employee.ID, &employee.Name, &employee.Birthday); err != nil {
			e.log.Error("failed to scan employee", errMsg.Err(err))
			return nil, err

		}
		employees = append(employees, employee)
	}

	return employees, nil
}
