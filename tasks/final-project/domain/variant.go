package domain

type Variant struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	VariantName string `json:"variant_name" validate:"required,min=3,max=50"`
	Quantity    int    `json:"quantity" validate:"required,min=0"`
	ProductID   int    `json:"product_id" validate:"required"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
