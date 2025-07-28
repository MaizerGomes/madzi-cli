package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	loadEnv()

	for {
		fmt.Println("\n🔹 Madzi CLI")
		fmt.Println("1) Select existing client")
		fmt.Println("2) Register new client")
		fmt.Println("3) Update client details")
		fmt.Println("4) Process invoices and send receipts for all clients")
		fmt.Println("5) Setup email configuration")
		fmt.Println("6) Exit")
		fmt.Print("Choose an option: ")

		choice := readLine()

		switch choice {
		case "1":
			client := selectClient()
			if client != nil {
				clientMenu(*client)
			}
		case "2":
			client := registerClient()
			if client != nil {
				clientMenu(*client)
			}
		case "3":
			fmt.Println("Update client details")

			client := selectClient()
			if client != nil {
				updated, done := fetchClientDetails(*client)
				if done {
					return
				}

				err := updateClient(*updated)
				if err != nil {
					return
				}

				fmt.Println("✅ Client updated successfully.")
			}
		case "4":
			ProcessInvoicesAndSendReceipts()
		case "5":
			createEnvFile()
		case "6":
			fmt.Println("👋 Goodbye.")
			os.Exit(0)
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		createEnvFile()
	}
}

func clientMenu(client Client) {
	for {
		fmt.Printf("\n🔸 Current Client: %s (%s - %s)\n", client.Name, client.Number, client.ClientDetails.PartnerName)
		fmt.Printf("🔸 Address: %s \n", client.ClientDetails.Address)
		fmt.Printf("🔸 Total Debt: %s MZM\n", client.ClientDetails.TotalDebt)
		fmt.Println("1) View unpaid invoices")
		fmt.Println("2) View receipts")
		fmt.Println("3) Back to main menu")
		fmt.Print("Choose an option: ")

		choice := readLine()

		switch choice {
		case "1":
			viewUnpaidInvoices(client)
		case "2":
			viewReceipts(client)
		case "3":
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}
func viewUnpaidInvoices(client Client) {
	fmt.Println("\n🔎 Fetching unpaid invoices...")

	invoices, err := fetchUnpaidInvoices(client.Number)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if len(invoices) == 0 {
		fmt.Println("✅ No unpaid invoices found.")
		return
	}

	fmt.Printf("📄 %d unpaid invoice(s) found:\n", len(invoices))
	for i, invoice := range invoices {
		fmt.Println("----------------------------------------")
		fmt.Printf("Option         : %d\n", i+1)
		fmt.Printf("Invoice Number : %s\n", invoice.InvoiceNumber)
		fmt.Printf("Date           : %s\n", invoice.Date[:10])
		fmt.Printf("Amount         : %s MZN\n", invoice.TotalAmount)
		fmt.Printf("Address        : %s\n", invoice.Address)
		fmt.Printf("Partner        : %s\n", invoice.PartnerName)
	}

	// Save the latest invoice number for change detection
	if len(invoices) > 0 {
		client.LastInvoice = invoices[0].InvoiceNumber
		err := updateClient(client)
		if err != nil {
			return
		}
	}

	index := promptInt("Select Invoice: ", 1, len(invoices)) - 1
	selected := invoices[index]
	fmt.Println("1. Pay")
	fmt.Println("2. Download")
	fmt.Println("3. Send via Email")
	fmt.Println("4. Back to menu")
	choice := promptInt("Choose option: ", 1, 4)

	switch choice {
	case 1:

		fmt.Printf("\nYou selected Invoice #%s | Total %s MZM\n", selected.GetNumber(), selected.TotalAmount)
		fmt.Printf("Client Phone Number: %s\n", client.Phone)

		confirm := prompt("Proceed with M-Pesa payment? [Y/n]: ")
		if confirm != "y" && confirm != "" {
			fmt.Println("❌ Cancelled.")
			return
		}

		if initiateInvoicePayment(client, selected) {
			return
		}
	case 2:
		err := downloadInvoicePDF(selected, true)
		if err != nil {
			fmt.Println("Erro ao fazer download:", err)
		} else {
			fmt.Println("Download completo em data/invoices/")
		}
	case 3:
		email := client.Email
		if email == "" {
			email = prompt("Email do cliente: ")
		}
		err := downloadInvoicePDF(selected, false)
		if err != nil {
			fmt.Println("Erro ao preparar ficheiro:", err)
			return
		}
		err = SendEmailWithAttachment(email, "Factura de água", "Segue em anexo a sua factura.", selected, filepath.Join("data/invoices", selected.InvoiceNumber+".pdf"))
		if err != nil {
			fmt.Println("Erro ao enviar email:", err)
		} else {
			fmt.Println("Email enviado com sucesso para", email)
		}
	case 4:
		return
	}
}

func initiateInvoicePayment(client Client, selected Invoice) bool {
	fmt.Println("📡 Initiating payment request...")

	payload := PaymentRequest{
		InvoiceNumber: selected.GetNumber(),
		PaymentMethod: "mpesa",
		PaymentFields: PaymentFields{
			PhoneNumber: client.Phone,
			Amount:      selected.TotalAmount,
		},
	}

	restyClient := resty.New()
	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("https://api.madzi.co.mz/api/v1/postpaid/payment")

	fmt.Println("✅ Payment request sent. Please confirm on your phone.")

	if err != nil {
		fmt.Println("❌ Error during payment request:", err)
		return true
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Printf("❌ Payment failed (%d): %s\n", resp.StatusCode(), resp.Body())
		return true
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return true
	}

	prettyPrintJSON(result)
	return false
}

func prettyPrintJSON(result map[string]interface{}) {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("❌ Error formatting JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func viewReceipts(client Client) {
	fmt.Println("\n🔎 Fetching receipts...")

	receipts, err := fetchAllReceipts(client.Number, 5, nil)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if len(receipts) == 0 {
		fmt.Println("No receipts found.")
		return
	}

	for i, r := range receipts {
		fmt.Printf("\n📄 Receipt: %s (# %d)\n", r.Number, i+1)
		fmt.Printf("   ➤ Date: %s\n", r.Date)
		fmt.Printf("   ➤ Address: %s\n", strings.TrimSpace(r.Address))
	}

	index := promptInt("Select receipt: ", 1, len(receipts)) - 1
	selected := receipts[index]
	fmt.Println("1. Download")
	fmt.Println("2. Email")
	choice := promptInt("Choose option: ", 1, 2)

	if choice == 1 {
		err := downloadReceiptPDF(selected, true)
		if err != nil {
			fmt.Println("Erro ao fazer download:", err)
		} else {
			fmt.Println("Download completo em data/receipts/")
		}
	} else {
		email := client.Email
		if email == "" {
			email = prompt("Email do cliente: ")
		}
		err := downloadReceiptPDF(selected, false)
		if err != nil {
			fmt.Println("Erro ao preparar ficheiro:", err)
			return
		}
		err = SendEmailWithAttachment(email, "Recibo de pagamento", "Segue em anexo o seu recibo de pagamento.", selected, filepath.Join("data/receipts", selected.Number+".pdf"))
		if err != nil {
			fmt.Println("Erro ao enviar email:", err)
		} else {
			fmt.Println("Email enviado com sucesso para", email)
		}
	}
}

func prompt(s string) string {
	fmt.Print(s)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func promptInt(s string, i int, i2 int) int {
	var choice int
	for {
		fmt.Print(s)
		_, err := fmt.Scanf("%d", &choice)
		if err != nil || choice < i || choice > i2 {
			fmt.Printf("Please choose a number between %d e %d.\n", i, i2)
			continue
		}
		break
	}
	return choice

}

func ProcessInvoicesAndSendReceipts() {
	fmt.Println("🔄 Iniciando o processamento automático de facturas e envio de recibos...")

	clients, err := loadClients()
	if err != nil {
		fmt.Println("❌ Erro ao carregar os clientes:", err)
		return
	}

	// 1. Processar facturas não pagas e iniciar pagamento
	for _, client := range clients {
		fmt.Printf("\n👤 Cliente: %s [%s]\n", client.Name, client.Number)

		invoices, err := fetchUnpaidInvoices(client.Number)
		if err != nil {
			fmt.Printf("  ⚠️ Erro ao buscar facturas: %v\n", err)
			continue
		}

		if len(invoices) == 0 {
			fmt.Println("  ✅ Sem facturas por pagar.")
			continue
		}

		for _, invoice := range invoices {
			fmt.Printf("  💳 Iniciando pagamento da factura #%s...\n", invoice.GetNumber())

			if initiateInvoicePayment(client, invoice) {
				fmt.Printf("    ❌ Falha no pagamento da factura #%s: %v\n", invoice.GetNumber())
			} else {
				fmt.Printf("    ✅ Pagamento iniciado para factura #%s\n", invoice.GetNumber())
			}
		}
	}

	// 2. Buscar e enviar os recibos do dia
	fmt.Println("\n📨 A enviar recibos emitidos hoje por email...")

	for _, client := range clients {
		// Obter data actual no formato "YYYY-MM-DD"
		today := time.Now().Format("2006-01-02")
		receipts, err := fetchAllReceipts(client.Number, 5, &today)

		if err != nil {
			fmt.Printf("  ⚠️ Erro ao buscar recibos para %s: %v\n", client.Name, err)
			continue
		}

		if len(receipts) == 0 {
			fmt.Printf("  ℹ️ Sem recibos hoje para %s\n", client.Name)
			continue
		}

		for _, receipt := range receipts {
			err = SendEmailWithAttachment(client.Email, "Recibo de pagamento", "Segue em anexo o seu recibo de pagamento.", receipt, filepath.Join("data/receipts", receipt.GetFilename()))
			if err != nil {
				fmt.Println("Erro ao enviar email:", err)
			} else {
				fmt.Println("Email enviado com sucesso para", client.Email)
			}
		}
	}

	fmt.Println("\n🎉 Processamento concluído.")
}

var requiredEnvVars = []string{
	"SMTP_HOST",
	"SMTP_PORT",
	"SMTP_USERNAME",
	"SMTP_PASSWORD",
	"EMAIL_SENDER_NAME",
	"EMAIL_SENDER_NAME",
	"EMAIL_SENDER_ADDRESS",
}

var defaultEnvVars = map[string]string{
	"SMTP_HOST":            "smtp.elasticemail.com",
	"SMTP_PORT":            "2525",
	"SMTP_USERNAME":        "adem@ekutiva.co.mz",
	"SMTP_PASSWORD":        "21A44FB0E64A72FF7B40B50421DD115B9BB1",
	"EMAIL_SENDER_NAME":    "'Madzi CLI'",
	"EMAIL_SENDER_ADDRESS": "adem@ekutiva.co.mz",
}

// createEnvFile prompts user and writes .env file
func createEnvFile() {
	file, err := os.Create(".env")
	if err != nil {
		log.Fatalf("❌ Failed to create .env file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("❌ Failed to close .env file: %v", err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Email configuration for Madzi CLI\n")
	choice := prompt("Use default values? (y/n):")

	if choice == "y" || choice == "" {
		fmt.Println("Using default email configuration values.")
	} else {

		fmt.Println("Please enter the following environment variables:")
	}

	for key, def := range defaultEnvVars {
		input := ""
		if choice == "n" {
			fmt.Printf("%s: ", key)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
		} else {
			input = def
		}

		line := fmt.Sprintf("%s=%s\n", key, input)
		_, err := writer.WriteString(line)
		if err != nil {
			return
		}
	}

	err = writer.Flush()
	if err != nil {
		return
	}
	fmt.Println(".env file created successfully.")

	loadEnv()
}
