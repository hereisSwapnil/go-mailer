# go-mailer

A Go application for sending bulk emails from CSV files using SMTP with concurrent processing and retry logic.

## Features

- Bulk email sending from CSV files
- Concurrent processing with worker pools
- HTML email templates with dynamic variables
- Retry logic with exponential backoff
- Environment-based configuration

## 🎥 Demo

<!-- Add your demo video here -->

*Click the image above to watch the demo video*

## 📸 Screenshots

[![Clean-Shot-2025-11-17-at-02-49-09-2x.png](https://i.postimg.cc/dVJQYKhP/Clean-Shot-2025-11-17-at-02-49-09-2x.png)](https://postimg.cc/qzZHsScj)

[![image.png](https://i.postimg.cc/PrjTWkzM/image.png)](https://postimg.cc/14JkS2dg)

[![Clean-Shot-2025-11-17-at-02-50-41-2x.png](https://i.postimg.cc/DwCh34Vp/Clean-Shot-2025-11-17-at-02-50-41-2x.png)](https://postimg.cc/nXQgvzhq)

[![Clean-Shot-2025-11-17-at-02-50-59-2x.png](https://i.postimg.cc/mkNxXtNk/Clean-Shot-2025-11-17-at-02-50-59-2x.png)](https://postimg.cc/y3drW6GC)

## Requirements

- Go 1.25.1 or later
- SMTP server credentials
- CSV file with recipient data

## Architecture

[![image.png](https://i.postimg.cc/439Lxc8H/image.png)](https://postimg.cc/hJKbCXRK)

## Installation

```bash
git clone https://github.com/hereisSwapnil/go-mailer.git
cd go-mailer
go mod download
```

## Configuration

Create `config/local.yaml`:

```yaml
smtp:
  host: "smtp.gmail.com"
  port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from: "your-email@gmail.com"

templates:
  test_email_template: "internal/templates/test_mail.tmpl"

data:
  csv_file_path: "data/csv/recipients.csv"

retry:
  max_retries: 3
  initial_delay_seconds: 1
  backoff_multiplier: 2.0
  max_delay_seconds: 30
```

## CSV Format

Required columns: `Name`, `Email`

Optional columns: Any additional columns (e.g., `Coupon`, `Discount`) are available in templates.

Example:
```csv
Name, Email, Coupon
John Doe, john@example.com, JOHN1000
Jane Smith, jane@example.com, JANE2000
```

## Email Templates

Templates use Go's `html/template` package with two sections: `subject` and `html`.

Example (`internal/templates/test_mail.tmpl`):
```html
{{define "subject"}}Thank you for joining, {{.Name}}. Here is your coupon{{end}}

{{define "html"}}
<!doctype html>
<html>
  <body style="font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif;">
    <p>Hello <strong>{{.Name}}</strong>,</p>
    <p>Thank you for joining us. We are glad to welcome you.</p>
    <p>Your coupon code is:</p>
    <p style="font-size: 20px; font-weight: bold;">{{.Coupon}}</p>
    <p>Enjoy your shopping.</p>
    <p>Warm regards,<br />Go Mailer Team</p>
  </body>
</html>
{{end}}
```

Template variables: All CSV columns are available (e.g., `{{.Name}}`, `{{.Email}}`, `{{.Coupon}}`).

## Usage

```bash
export ENV=local  # Optional, defaults to "local"
go run cmd/go-mailer/main.go
```

Or build and run:
```bash
go build -o go-mailer cmd/go-mailer/main.go
./go-mailer
```

## How It Works

1. Loads configuration from `config/{ENV}.yaml`
2. Producer reads CSV and sends recipients to a channel
3. Worker pool (10 workers) processes recipients concurrently
4. Each worker renders templates and sends emails via SMTP
5. Failed sends are retried with exponential backoff

## Project Structure

```
go-mailer/
├── cmd/go-mailer/main.go
├── config/local.yaml
├── data/csv/recipients.csv
├── internal/
│   ├── config/config.go
│   ├── consumer/consumer.go
│   ├── producer/producer.go
│   ├── templates/test_mail.tmpl
│   └── types/types.go
├── go.mod
└── README.md
```