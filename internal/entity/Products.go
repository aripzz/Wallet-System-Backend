package entity

type Products struct {
	ID       uint64  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type CreateProducts struct {
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Quantity int     `json:"quantity" validate:"required"`
}

type UpdateProducts struct {
	ID       uint64   `json:"id" validate:"required"`
	Name     *string  `json:"name,omitempty"`
	Price    *float64 `json:"price,omitempty"`
	Quantity *int     `json:"quantity,omitempty"`
}
