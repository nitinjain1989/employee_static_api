package dto

type MobileRequest struct {
	Type   string `json:"type" example:"home"`
	Number string `json:"number" example:"9876543210"`
}

type CreateEmployeeRequest struct {
	ID          string          `json:"id"`
	Name        string          `json:"name" binding:"required" example:"Nitin Jain"`
	Email       string          `json:"email" binding:"required,email" example:"nitin@test.com"`
	Designation string          `json:"designation" example:"iOS Engineer"`
	Department  string          `json:"department" example:"Engineering"`
	City        string          `json:"city" example:"Noida"`
	IsActive    bool            `json:"is_active"`
	Country     string          `json:"country" example:"India"`
	ImgURL      string          `json:"img_url" example:"https://image.com/profile.jpg"`
	JoiningDate string          `json:"joining_date" example:"2024-01-01"`
	Mobiles     []MobileRequest `json:"mobiles"`
}

type UpdateEmployeeRequest struct {
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Designation string          `json:"designation"`
	Department  string          `json:"department"`
	IsActive    bool            `json:"is_active"`
	ImgURL      string          `json:"img_url"`
	City        string          `json:"city"`
	Country     string          `json:"country"`
	JoiningDate string          `json:"joining_date"`
	Version     int             `json:"version"` // 🔥 required
	Mobiles     []MobileRequest `json:"mobiles"`
}
