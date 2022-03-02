package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

var warranties = make([]Warranty, 0)

type IsoDate time.Time

func (j *IsoDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = IsoDate(t)
	return nil
}

func (j IsoDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.String())
}

func (j IsoDate) String() string {
	return time.Time(j).Format("2006-01-02")
}

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
	StartDate IsoDate      `json:"start_date"`
	EndDate   IsoDate      `json:"end_date"`
}

type Warranty struct {
	ID       uuid.UUID `json:"id"`
	Product  `json:"product"`
	Contract `json:"contract"`
}

func NewWarranty(w http.ResponseWriter, req *http.Request) (Warranty, error) {
	var warranty Warranty

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&warranty)
	if err != nil {
		fmt.Println(err)
		return warranty, err
	}

	warranty.ID = uuid.New()
	warranty.Product.ID = uuid.New()
	warranty.Contract.ID = uuid.New()

	return warranty, nil
}

func (warranty Warranty) Save() {
	warranties = append(warranties, warranty)
}

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

func serializeWarranties(w http.ResponseWriter, warranties []Warranty) {
	for _, warranty := range warranties {
		serializeWarranty(w, warranty)
	}
}

func Warranties(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		warranty, err := NewWarranty(w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		warranty.Save()
		serializeWarranty(w, warranty)

	case "GET":
		warranties := getWarranties()
		serializeWarranties(w, warranties)

	default:
		http.Error(w, fmt.Sprintf("Method %v not supported.", req.Method), http.StatusNotFound)
		return
	}
}

func Home(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		warranties := getWarranties()
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, warranties)

	default:
		http.Error(w, fmt.Sprintf("Method %v not supported.", req.Method), http.StatusNotFound)
		return
	}
}

func main() {
	http.HandleFunc("/warranties", Warranties)
	http.HandleFunc("/", Home)
	http.ListenAndServe(":8080", nil)
}
