package models

type MobileRequest struct {
	Type   string `json:"type" example:"home"`
	Number string `json:"number" example:"9876543210"`
}

type CreateEmployeeRequest struct {
	Name        string          `json:"name" binding:"required" example:"Nitin Jain"`
	Email       string          `json:"email" binding:"required,email" example:"nitin@test.com"`
	Designation string          `json:"designation" example:"iOS Engineer"`
	Department  string          `json:"department" example:"Engineering"`
	City        string          `json:"city" example:"Noida"`
	Country     string          `json:"country" example:"India"`
	ImgURL      string          `json:"img_url" example:"https://image.com/profile.jpg"`
	JoiningDate string          `json:"joining_date" example:"2024-01-01"`
	Mobiles     []MobileRequest `json:"mobiles"`
}
