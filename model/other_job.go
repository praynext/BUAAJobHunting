package model

import "time"

type OtherJob struct {
	ID          int       `json:"id" db:"id"`
	JobSRC      string    `json:"job_src" db:"job_src"`
	JobName     string    `json:"job_name" db:"job_name"`
	JobArea     string    `json:"job_area" db:"job_area"`
	Salary      string    `json:"salary" db:"salary"`
	CompanyName string    `json:"company_name" db:"company_name"`
	JobNeed     string    `json:"job_need" db:"job_need"`
	JobDesc     string    `json:"job_desc" db:"job_desc"`
	JobURL      string    `json:"job_url" db:"job_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsFull      bool      `json:"is_full" db:"is_full"`
	Tokens      string    `json:"tokens" db:"tokens"`
}
