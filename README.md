# 💧 madzi-cli

**A simple, interactive CLI to manage your postpaid water invoices with ADRM (Mozambique).**

> Built with Go · Cross-platform · Secure · Blazing fast

---

## 📦 Features

- 🔍 Fetch unpaid invoices by contract number  
- 🧾 Retrieve and download receipts in PDF format  
- 💰 Initiate M-Pesa payments for invoices  
- 📧 Email receipts directly to clients  
- 🧠 Intelligent tracking of invoice changes *(coming soon)*  
- 🛠️ Interactive TUI (text user interface) — No need to remember commands  
- 🔐 Environment variable setup with defaults and `.env` generation  

---

## 🚀 Quick Start

### ✅ Prerequisites

- macOS (Apple Silicon or Intel)
- Homebrew (package manager): [Install Homebrew](https://brew.sh)

### 🛠️ Installation via Homebrew

```sh
brew tap MaizerGomes/homebrew-madzi-cli
brew install madzi-cli

🖥️ Running the CLI

madzi-cli

On first run, you’ll be prompted to configure required environment variables (email). These will be saved to a .env file in the current directory.


📸 Screenshots

Coming soon: Visual walkthrough of CLI features and usage.

⸻

🧪 Development
	1.	Clone the repo:

git clone https://github.com/MaizerGomes/madzi-cli.git
cd madzi-cli

	2.	Build the CLI:

go build -o madzi-cli
./madzi-cli


⸻

🏗 Roadmap
	•	Interactive invoice lookup by contract number
	•	Email receipts with attachments
	•	PDF receipt download
	•	Homebrew distribution
	•	Background invoice checking
	•	System tray notifier (GUI support)
	•	Docker support for server-side usage

⸻

🛡 Security Notes
	•	Your .env file contains sensitive credentials. Make sure it is .gitignored and not committed.
	•	Email passwords are used only for SMTP authentication. We recommend using app passwords where available.

⸻

👤 Author

Maizer Gomes
github.com/MaizerGomes

⸻

📝 License

This project is licensed under the MIT License. See the LICENSE file for details.
