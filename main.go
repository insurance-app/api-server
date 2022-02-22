package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID    uuid.UUID `json:"id"`
	Brand string    `json:"brand"`
	Model string    `json:"model"`
	Color string    `json:"color"`
	Price string    `json:"price"`
}

type ContractType string

const (
	StandardContract ContractType = "standard"
	PremiumContract  ContractType = "premium"
)

type Contract struct {
	ID        uuid.UUID    `json:"id"`
	Type      ContractType `json:"type"`
	StartDate time.Time    `json:"start_date"`
	EndDate   time.Time    `json:"end_date"`
}

type Warranty struct {
	ID       uuid.UUID `json:"id"`
	Product  `json:"product"`
	Contract `json:"contract"`
}

func NewWarranty(req *http.Request) Warranty {
	warranty := Warranty{
		ID: uuid.New(),
		Product: Product{
			ID:    uuid.New(),
			Brand: "something",
			Model: "something",
			Color: "something",
			Price: "something",
		},
		Contract: Contract{
			ID:        uuid.New(),
			Type:      StandardContract,
			StartDate: time.Time{},
			EndDate:   time.Time{},
		},
	}

	return warranty
}

func (warranty Warranty) Save() {
	warranties = append(warranties, warranty)
}

var warranties = make([]Warranty, 0)

func getWarranties() []Warranty {
	return warranties
}

func serializeWarranty(w http.ResponseWriter, warranty Warranty) {
	warrantyJson, err := json.MarshalIndent(warranty, "", "    ")

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(warrantyJson))
}

func returnWarranties(w http.ResponseWriter, warranties []Warranty) {
	for _, warranty := range warranties {
		serializeWarranty(w, warranty)
	}
}

func Warranties(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		warranty := NewWarranty(req)
		warranty.Save()
		serializeWarranty(w, warranty)

	case "GET":
		warranties := getWarranties()
		returnWarranties(w, warranties)

	default:
		http.Error(w, fmt.Sprintf("Method %v not supported.", req.Method), http.StatusNotFound)
		return
	}
}

func GetWarranties(w http.ResponseWriter, req *http.Request) {

}

func main() {
	http.HandleFunc("/warranties", Warranties)
	http.ListenAndServe(":8080", nil)
}
