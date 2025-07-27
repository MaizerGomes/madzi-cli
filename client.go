package main

type Client struct {
	Name          string        `json:"name"`
	Number        string        `json:"number"` // Contract number
	Phone         string        `json:"phone"`
	Email         string        `json:"email"`
	LastInvoice   string        `json:"last_invoice,omitempty"` // used for change detection
	ClientDetails ClientDetails `json:"client_details"`
}

type ClientDetails struct {
	PartnerName     string `json:"partner_name"`
	ContractAccount string `json:"contract_account"`
	Address         string `json:"address"`
	TotalDebt       string `json:"total_debt"`
}
