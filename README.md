# go-mailer

A Go application for sending bulk emails from CSV files using SMTP. It uses a worker pool pattern to process recipients concurrently and includes retry logic with exponential backoff.

## Overview

The application reads recipient data from a CSV file, processes them through multiple worker goroutines, and sends personalized emails using HTML templates. Failed sends are retried with configurable backoff.

## Features

- Reads recipients from CSV files
- Concurrent email sending using worker pools
- HTML email templates with subject and body
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

The CSV file should have a header row and two columns:

```csv
Name, Email
John Doe, john@example.com
Jane Smith, jane@example.com
```

The first row is treated as a header and skipped. Each subsequent row should contain:
1. Name (column 0)
2. Email address (column 1)

Example file: `data/csv/recipients.csv`

## Email Templates

Email templates use Go's `html/template` package. Templates must define two sections:

- `subject`: The email subject line
- `html`: The HTML email body

Example template (`internal/templates/test_mail.tmpl`):

```html
{{define "subject"}}Test Email for {{.Name}}{{end}}

{{define "html"}}
<!doctype html>
<html>
  <body style="font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif;">
    <p>Hello <strong>{{.Name}}</strong>,</p>
    <p>This is a test email sent to you as part of our system check.</p>
    <p>Best regards,<br />Go Mailer Team</p>
  </body>
</html>
{{end}}
```

Template variables:
- `.Name`: Recipient name from CSV

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

2. **Producer**: A goroutine reads the CSV file and sends recipient data to a channel, skipping the header row.

3. **Workers**: By default, 10 worker goroutines consume from the channel. Each worker:
   - Loads the email template
   - Processes each recipient from the channel
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

