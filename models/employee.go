package models

import "time"

type Mobile struct {
	ID         string `json:"id,omitempty"`
	EmployeeID string `json:"employee_id,omitempty"`
	Type       string `json:"type"` // home, office, other
	Number     string `json:"number"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int        `json:"version,omitempty"`
}

type Employee struct {
	ID          string
	Name        string
	Designation string
	Department  string
	IsActive    bool
	ImgURL      string
	Email       string
	City        string
	Country     string
	JoiningDate *time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Version     int

	Mobiles []Mobile
}
