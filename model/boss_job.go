package model

import "time"

type BossJob struct {
	ID             int       `json:"id" db:"id"`
	JobName        string    `json:"job_name" db:"job_name"`
	JobArea        string    `json:"job_area" db:"job_area"`
	Salary         string    `json:"salary" db:"salary"`
	TagList        string    `json:"tag_list" db:"tag_list"`
	HRInfo         string    `json:"hr_info" db:"hr_info"`
	CompanyLogo    string    `json:"company_logo" db:"company_logo"`
	CompanyName    string    `json:"company_name" db:"company_name"`
	CompanyTagList string    `json:"company_tag_list" db:"company_tag_list"`
	CompanyURL     string    `json:"company_url" db:"company_url"`
	JobNeed        string    `json:"job_need" db:"job_need"`
	JobDesc        string    `json:"job_desc" db:"job_desc"`
	JobURL         string    `json:"job_url" db:"job_url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	Tokens         string    `json:"tokens" db:"tokens"`
}
