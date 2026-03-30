package dto

import "time"

type SyncRequest struct {
	Cursor    *Cursor    `json:"cursor,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Employees []Employee `json:"employees,omitempty"`
}

type Cursor struct {
	Seq int64 `json:"seq"`
}

func (c SyncRequest) IsEmpty() bool {
	return len(c.Employees) == 0
}

type Mobile struct {
	ID         string     `json:"id"`
	EmployeeID string     `json:"employee_id"`
	Type       string     `json:"type"`
	Number     string     `json:"number"`
	Version    int        `json:"version"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type Employee struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Designation string     `json:"designation"`
	Department  string     `json:"department"`
	IsActive    bool       `json:"is_active"`
	ImgURL      string     `json:"img_url"`
	Email       string     `json:"email"`
	City        string     `json:"city"`
	Country     string     `json:"country"`
	JoiningDate string     `json:"joining_date"`
	Version     int        `json:"version"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Mobiles     []Mobile   `json:"mobiles,omitempty"`
}
