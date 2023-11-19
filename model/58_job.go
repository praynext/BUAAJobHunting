package model

import "time"

type TC58Job struct {
	ID          int       `json:"id" db:"id"`
	JobName     string    `json:"job_name" db:"job_name"`
	JobArea     string    `json:"job_area" db:"job_area"`
	Salary      string    `json:"salary" db:"salary"`
	JobWel      string    `json:"job_wel" db:"job_wel"`
	CompanyName string    `json:"company_name" db:"company_name"`
	JobNeed     string    `json:"job_need" db:"job_need"`
	JobURL      string    `json:"job_url" db:"job_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Tokens      string    `json:"tokens" db:"tokens"`
}
