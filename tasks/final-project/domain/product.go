package domain

type Product struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name" validate:"required,min=3,max=50"`
	ImageURL  string `json:"image_url" validate:"omitempty,url"`
	AdminID   int    `json:"admin_id" validate:"required"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
