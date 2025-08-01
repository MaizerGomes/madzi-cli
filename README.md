# 💧 madzi-cli ![GitHub](https://img.shields.io/github/license/MaizerGomes/madzi-cli) ![GitHub stars](https://img.shields.io/github/stars/MaizerGomes/madzi-cli?style=social)


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
#### MacOS or Linux system with the following prerequisites:
- Homebrew (package manager): [Install Homebrew](https://brew.sh)

#### Windows:
- Download from Releases section or build from source using Go.

### 🛠️ Installation via Homebrew

```sh
brew tap MaizerGomes/homebrew-madzi-cli
brew install madzi-cli
```
🖥️ Running the CLI
```sh
madzi-cli
```
On first run, you’ll be prompted to configure required environment variables (email). These will be saved to a .env file in the current directory.


📸 Screenshots

Coming soon: Visual walkthrough of CLI features and usage.

⸻

🧪 Development
1.	Clone the repo:
```sh
git clone https://github.com/MaizerGomes/madzi-cli.git
cd madzi-cli
```

2. Run the CLI:
```sh
go run .
```

3. Build the CLI:
```sh
go build -o madzi-cli
./madzi-cli
```

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

Maizer Gomes ![GitHub](https://img.shields.io/github/followers/MaizerGomes?style=social)


github.com/MaizerGomes

⸻

📝 License

This project is licensed under the MIT License. See the LICENSE file for details.
