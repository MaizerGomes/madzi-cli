package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	_ "path/filepath"
	"runtime"
)

const dataDir = "data"
const dataFile = "data/clients.json"

func ensureDataFile() error {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
			return err
		}
	}
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return os.WriteFile(dataFile, []byte("[]"), 0644)
	}
	return nil
}

func saveClient(client Client) error {
	clients, _ := loadClients()
	clients = append(clients, client)

	data, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func loadClients() ([]Client, error) {
	if err := ensureDataFile(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	var clients []Client
	err = json.Unmarshal(data, &clients)
	return clients, err
}
func updateClient(updated Client) error {
	clients, err := loadClients()
	if err != nil {
		return err
	}

	for i, c := range clients {
		if c.Number == updated.Number {
			clients[i] = updated
			break
		}
	}

	data, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, data, 0644)
}

func downloadReceiptPDF(receipt Receipt, openFile bool) error {
	urlGenerateFile := fmt.Sprintf("https://api.madzi.co.mz/api/v1/postpaid/receipt?receipt=%s", receipt.GetNumber())
	url := fmt.Sprintf("https://api.madzi.co.mz/storage/receipts/%s.pdf", receipt.GetNumber())

	resp, err := http.Get(urlGenerateFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resp, err = http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = os.MkdirAll("data/receipts", 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join("data/receipts", receipt.GetFilename()))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	if openFile {
		if err := OpenPDF(f.Name()); err != nil {
			return fmt.Errorf("erro ao abrir o ficheiro PDF: %w", err)
		}
	}
	return err
}
func downloadInvoicePDF(invoice Invoice, openFile bool) error {
	urlGenerateFile := fmt.Sprintf("https://api.madzi.co.mz/api/v1/postpaid/invoice?invoice=%s", invoice.GetNumber())
	url := fmt.Sprintf("https://api.madzi.co.mz/storage/invoices/%s.pdf", invoice.GetNumber())
	resp, err := http.Get(urlGenerateFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resp, err = http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = os.MkdirAll("data/invoices", 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join("data/invoices", invoice.GetFilename()))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	if openFile {
		if err := OpenPDF(f.Name()); err != nil {
			return fmt.Errorf("erro ao abrir o ficheiro PDF: %w", err)
		}
	}

	return err
}

// OpenPDF abre um ficheiro PDF com a aplicação por defeito do sistema operativo
func OpenPDF(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// No Windows usa-se "start", mas tem que ser executado via "cmd"
		cmd = exec.Command("cmd", "/C", "start", "", filePath)
	case "darwin":
		// macOS
		cmd = exec.Command("open", filePath)
	case "linux":
		// Em muitos distros Linux
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("sistema operativo não suportado: %s", runtime.GOOS)
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("erro ao abrir o ficheiro PDF: %w", err)
	}

	return nil
}
