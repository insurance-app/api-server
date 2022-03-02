package warranty

import (
	"github.com/MrNeocore/sunrise-api-server/date"
	"github.com/google/uuid"
)

type Product struct {
	ID    uuid.UUID `json:"id"`
	Brand string    `json:"brand"`
	Model string    `json:"model"`
	Color string    `json:"color"`
	Price float32   `json:"price"`
}

type ContractType string

const (
	StandardContract ContractType = "standard"
	PremiumContract  ContractType = "premium"
)

type Contract struct {
	ID        uuid.UUID    `json:"id"`
	Type      ContractType `json:"type"`
	StartDate date.Date    `json:"start_date"`
	EndDate   date.Date    `json:"end_date"`
}

type Warranty struct {
	ID       uuid.UUID `json:"id"`
	Product  `json:"product"`
	Contract `json:"contract"`
}
