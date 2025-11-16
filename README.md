# go-mailer

A Go application for sending bulk emails from CSV files using SMTP. It uses a worker pool pattern to process recipients concurrently and includes retry logic with exponential backoff.

## Overview

The application reads recipient data from a CSV file, processes them through multiple worker goroutines, and sends personalized emails using HTML templates. Failed sends are retried with configurable backoff.

## Features

- Reads recipients from CSV files with dynamic column support
- Concurrent email sending using worker pools
- HTML email templates with subject and body
- Dynamic template variables from CSV columns
- Retry logic with exponential backoff
- Environment-based configuration
- Structured logging

## Requirements

- Go 1.25.1 or later
- SMTP server credentials
- CSV file with recipient data

## Installation

Clone the repository:

```bash
git clone https://github.com/hereisSwapnil/go-mailer.git
cd go-mailer
```

Install dependencies:

```bash
go mod download
```

## Configuration

The application loads configuration from YAML files based on the `ENV` environment variable. By default, it uses `config/local.yaml`.

### Configuration File Structure

Create a configuration file at `config/local.yaml` (or `config/{ENV}.yaml`):

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

### Configuration Fields

**SMTP:**
- `host`: SMTP server hostname
- `port`: SMTP server port (typically 587 for TLS)
- `username`: SMTP authentication username (usually your email)
- `password`: SMTP authentication password or app password
- `from`: Email address to send from

**Templates:**
- `test_email_template`: Path to the email template file

**Data:**
- `csv_file_path`: Path to the CSV file containing recipients

**Retry:**
- `max_retries`: Maximum number of retry attempts (0-10)
- `initial_delay_seconds`: Initial delay before first retry (seconds)
- `backoff_multiplier`: Multiplier for exponential backoff (must be >= 1)
- `max_delay_seconds`: Maximum delay between retries (seconds)

## CSV Format

The CSV file should have a header row with at least two required columns: `Name` and `Email`. Additional columns are supported and will be automatically available in email templates.

**Required columns:**
- `Name`: Recipient's name
- `Email`: Recipient's email address

**Optional columns:**
- Any additional columns (e.g., `Coupon`, `Discount`, `ExpiryDate`) will be automatically available in templates

Example CSV file:

```csv
Name, Email, Coupon
John Doe, john@example.com, JOHN1000
Jane Smith, jane@example.com, JANE2000
Bob Johnson, bob@example.com, BOB3000
```

**Notes:**
- The first row is treated as a header and is skipped during processing
- Column names are case-insensitive (e.g., "Name", "name", "NAME" are all valid)
- Additional columns beyond Name and Email are stored and accessible in templates with the first letter capitalized (e.g., "Coupon" column → `{{.Coupon}}` in template)
- Column order doesn't matter; the application automatically maps columns by header name

Example file: `data/csv/recipients.csv`

## Email Templates

Email templates use Go's `html/template` package. Templates must define two sections:

- `subject`: The email subject line
- `html`: The HTML email body

Example template (`internal/templates/test_mail.tmpl`):

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

**Template variables:**
- `.Name`: Recipient name from CSV (required)
- `.Email`: Recipient email address from CSV (required)
- `.Coupon`: Any additional CSV column (e.g., if CSV has "Coupon" column)
- `.{ColumnName}`: All CSV columns are automatically available in templates with the first letter capitalized

**Note:** Any column in your CSV file (beyond Name and Email) can be accessed in templates using `{{.ColumnName}}` where `ColumnName` matches the CSV header with the first letter capitalized. For example:
- CSV column "Coupon" → `{{.Coupon}}` in template
- CSV column "Discount" → `{{.Discount}}` in template
- CSV column "ExpiryDate" → `{{.ExpiryDate}}` in template

## Usage

Set the environment variable (optional, defaults to "local"):

```bash
export ENV=local
```

Run the application:

```bash
go run cmd/go-mailer/main.go
```

Or build and run:

```bash
go build -o go-mailer cmd/go-mailer/main.go
./go-mailer
```

## How It Works

1. **Configuration Loading**: The application loads configuration from `config/{ENV}.yaml` and validates all required fields.

2. **Producer**: A goroutine reads the CSV file and sends recipient data to a channel. It:
   - Parses the header row to map column names dynamically
   - Extracts required fields (Name, Email) and all additional columns
   - Stores additional columns for template access
   - Skips the header row when processing data

3. **Workers**: By default, 10 worker goroutines consume from the channel. Each worker:
   - Loads the email template once
   - Processes each recipient from the channel
   - Builds template data with Name, Email, and all additional CSV columns (with capitalized first letters)
   - Renders the template with recipient data
   - Sends the email via SMTP
   - Retries on failure with exponential backoff

4. **Retry Logic**: If an email send fails, the worker waits and retries up to `max_retries` times. The delay increases exponentially: `initial_delay * (backoff_multiplier ^ attempt)`, capped at `max_delay_seconds`.

5. **Completion**: The application waits for all workers to finish processing all recipients.

## Project Structure

```
go-mailer/
├── cmd/
│   └── go-mailer/
│       └── main.go          # Application entry point
├── config/
│   └── local.yaml           # Configuration file
├── data/
│   └── csv/
│       └── recipients.csv   # Recipient data
├── internal/
│   ├── config/
│   │   └── config.go       # Configuration loading and validation
│   ├── consumer/
│   │   └── consumer.go     # Email worker implementation
│   ├── producer/
│   │   └── producer.go     # CSV reader and channel producer
│   ├── templates/
│   │   └── test_mail.tmpl  # Email template
│   └── types/
│       └── types.go        # Type definitions
├── go.mod
└── README.md
```

## Architecture

[![image.png](https://i.postimg.cc/439Lxc8H/image.png)](https://postimg.cc/hJKbCXRK)