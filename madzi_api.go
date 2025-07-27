package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const baseURL = "https://api.madzi.co.mz/api/v1/postpaid"

type Invoice struct {
	InvoiceNumber   string `json:"invoice_number"`
	Date            string `json:"date"`
	ContractAccount string `json:"contract_account"`
	PartnerName     string `json:"partner_name"`
	TotalAmount     string `json:"total_amount"`
	Address         string `json:"address"`
	TotalDebt       string `json:"total_debt"`
	Balance         string `json:"balance"`
}

func (i Invoice) GetNumber() string {
	return i.InvoiceNumber
}

func (i Invoice) GetFilename() string {
	return fmt.Sprintf("%s.pdf", i.InvoiceNumber)
}

type InvoiceResponse struct {
	Invoices     []Invoice `json:"invoices"`
	TotalPages   int       `json:"total_pages"`
	CurrentPage  string    `json:"current_page"`
	TotalRecords int       `json:"total_records"`
}

type Receipt struct {
	Number          string `json:"number"`
	Date            string `json:"date"`
	ContractAccount string `json:"contract_account"`
	PartnerName     string `json:"partner_name"`
	Address         string `json:"address"`
}

func (r Receipt) GetNumber() string {
	return r.Number
}

func (r Receipt) GetFilename() string {
	return fmt.Sprintf("%s.pdf", r.Number)
}

type ReceiptResponse struct {
	Receipts     []Receipt `json:"receipts"`
	TotalPages   int       `json:"total_pages"`
	CurrentPage  string    `json:"current_page"`
	TotalRecords int       `json:"total_records"`
}

type PaymentRequest struct {
	InvoiceNumber string        `json:"invoice_number"`
	PaymentMethod string        `json:"payment_method"`
	PaymentFields PaymentFields `json:"payment_fields"`
}

type PaymentFields struct {
	PhoneNumber string `json:"phone_number"`
	Amount      string `json:"amount"`
}

func fetchUnpaidInvoices(contract string) ([]Invoice, error) {
	page := 1
	var all []Invoice

	for {
		url := fmt.Sprintf("https://api.madzi.co.mz/api/v1/postpaid/invoices?page=%d&contract=%s", page, contract)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching invoices:", err)
			break
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		body, _ := io.ReadAll(resp.Body)
		var response InvoiceResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		all = append(all, response.Invoices...)

		if page >= response.TotalPages {
			break
		}
		page++
	}

	return all, nil
}

func fetchAllReceipts(contract string, max int, filterDate *string) ([]Receipt, error) {
	page := 1
	var all []Receipt

	for {
		url := fmt.Sprintf("https://api.madzi.co.mz/api/v1/postpaid/receipts?page=%d&contract=%s", page, contract)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching receipts:", err)
			break
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		body, _ := io.ReadAll(resp.Body)
		var response ReceiptResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, nil
		}

		all = append(all, response.Receipts...)

		if len(all) >= max {
			break
		}

		if page >= response.TotalPages {
			break
		}
		page++
	}

	if filterDate != nil {
		var filtered []Receipt
		for _, r := range all {
			datePart := strings.Split(r.Date, "T")[0]

			if &datePart == filterDate {
				filtered = append(filtered, r)
			}
		}
		all = filtered
	}

	return all, nil
}
