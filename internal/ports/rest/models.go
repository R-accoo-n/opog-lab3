package rest

type CategoryPayload struct {
	Name string  `json:"name"`
	Tax  float64 `json:"tax"`
}

type CreateProductRequest struct {
	Name     string          `json:"name"`
	Category CategoryPayload `json:"category"`
	Price    float64         `json:"price"`
}

type ProductResponse struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Category   CategoryPayload `json:"category"`
	Price      float64         `json:"price"`
	FinalPrice float64         `json:"final_price"`
}
