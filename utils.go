package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func selectClient() *Client {
	clients, err := loadClients()
	if err != nil || len(clients) == 0 {
		fmt.Println("❌ No clients found.")
		return nil
	}

	fmt.Println("\n📋 Select a Client:")
	for i, c := range clients {
		fmt.Printf("%d) %s (%s)\n", i+1, c.Name, c.Number)
	}

	fmt.Print("Enter number: ")
	input := readLine()

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(clients) {
		fmt.Println("❌ Invalid selection.")
		return nil
	}

	return &clients[idx-1]
}

func registerClient() *Client {
	fmt.Println("\n📝 Register a New Client")

	fmt.Print("Name: ")
	name := readLine()

	fmt.Print("Contract Number: ")
	number := readLine()

	fmt.Print("Phone Number: ")
	phone := readLine()

	fmt.Print("Email: ")
	email := readLine()

	client := Client{
		Name:   name,
		Number: number,
		Phone:  phone,
		Email:  email,
	}

	c, done := fetchClientDetails(client)
	if done {
		return c
	}

	if err := saveClient(*c); err != nil {
		fmt.Println("❌ Failed to save client:", err)
		return nil
	}

	fmt.Println("✅ Client registered successfully.")
	return c
}

func fetchClientDetails(client Client) (*Client, bool) {
	invoices, _ := fetchUnpaidInvoices(client.Number)

	if len(invoices) > 0 {
		client.LastInvoice = invoices[0].GetNumber()
		client.ClientDetails = ClientDetails{
			PartnerName:     invoices[0].PartnerName,
			ContractAccount: invoices[0].ContractAccount,
			Address:         invoices[0].Address,
			TotalDebt:       invoices[0].Balance,
		}

		fmt.Printf("Client Name: %s\n", client.ClientDetails.PartnerName)
		fmt.Printf("Client Address: %s\n", client.ClientDetails.Address)
		fmt.Printf("Total Debt: %s MZM\n", client.ClientDetails.TotalDebt)
		fmt.Print("Are these details correct? [y/n]: ")
		confirm := readLine()
		if strings.ToLower(confirm) != "y" {
			fmt.Println("❌ Client registration cancelled.")
			return nil, true
		}
	}
	return &client, false
}
