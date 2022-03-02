package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/MrNeocore/sunrise-api-server/warranty"
	"github.com/google/uuid"
)

var warranties = make([]warranty.Warranty, 0)

func NewWarranty(w http.ResponseWriter, req *http.Request) (warranty.Warranty, error) {
	var warranty warranty.Warranty

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

func getWarranties() []warranty.Warranty {
	return warranties
}

func saveWarranties(warranty warranty.Warranty) {
	warranties = append(warranties, warranty)
}

func serializeWarranty(w http.ResponseWriter, warranty warranty.Warranty) {
	warrantyJson, err := json.MarshalIndent(warranty, "", "    ")

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(warrantyJson))
}

func serializeWarranties(w http.ResponseWriter, warranties []warranty.Warranty) {
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
		saveWarranties(warranty)
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
