package helper

import (
	"bytes"
	"encoding/json"
	"finalp2/models"
	"net/http"
)

type Invoice struct {
	ID         string `json:"id"`
	InvoiceUrl string `json:"invoice_url"`
}

func CreateInvoice(product models.Rental, customer models.User, books []models.Book) (*Invoice, error) {
	apiKey := "xnd_development_K6nwB5ulAFPVXsOOKjan9yZUhstq81qZ5yz3CT5gY1vx3Aof8VwuffIcyHjfm"
	apiUrl := "https://api.xendit.co/v2/invoices"

	items := []map[string]interface{}{}

	// Iterate over the books and add each to the items array
	for _, book := range books {
		item := map[string]interface{}{
			"name":     book.Title, 
			"price":    book.Price,
			"quantity": 1, 
		}
		items = append(items, item)
	}

	bodyRequest := map[string]interface{}{
		"external_id":      "1",
		"amount":           product.TotalPrice,
		"description":      "Dummy Invoice RMT003",
		"invoice_duration": 86400,
		"customer": map[string]interface{}{
			"name":  customer.FirstName,
			"email": customer.Email,
		},
		"currency": "IDR",
		"items": items,
	}

	reqBody, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var resInvoice Invoice
	if err := json.NewDecoder(response.Body).Decode(&resInvoice); err != nil {
		return nil, err
	}

	return &resInvoice, nil

}
