package domain

type Admin struct {
	ID       int    `json:"id" form:"id"`
	UUID     string `json:"uuid" form:"uuid"`
	Name     string `json:"name" form:"name" binding:"required" validate:"required,min=3,max=32"`
	Email    string `json:"email" form:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" form:"password" binding:"required" validate:"required,password"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" form:"password" binding:"required" validate:"required,min=6"`
}
