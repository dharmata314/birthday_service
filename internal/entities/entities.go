package entities

import "time"

type User struct {
	ID        int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type Subscription struct {
	ID         int
	UserID     int
	EmployeeID int
}

type Employee struct {
	ID       int       `json:"employee_id"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
}
