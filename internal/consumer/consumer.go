package consumer

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
	"strings"
	"time"

	"github.com/hereisSwapnil/go-mailer/internal/config"
	"github.com/hereisSwapnil/go-mailer/internal/types"
)

func EmailWorker(workerId int, emailChannel chan types.Recipient, cfg *config.Config) error {
	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)

	t, err := template.ParseFiles(cfg.Templates.TestEmailTemplate)
	if err != nil {
		return err
	}

	for recipient := range emailChannel {
		slog.Info("Processing recipient", "workerId", workerId, "email", recipient.Email)

		data := map[string]interface{}{
			"Name":  recipient.Name,
			"Email": recipient.Email,
		}

		// Add any extra fields to the template data with capitalized first letter
		for key, value := range recipient.Extra {
			// Capitalize first letter for template access (e.g., "coupon" -> "Coupon")
			if len(key) > 0 {
				capitalizedKey := strings.ToUpper(key[:1]) + key[1:]
				data[capitalizedKey] = value
			}
		}

		var subjectBuffer bytes.Buffer
		if err := t.ExecuteTemplate(&subjectBuffer, "subject", data); err != nil {
			return fmt.Errorf("Subject template error (%s): %w", recipient.Email, err)
		}

		var bodyBuffer bytes.Buffer
		if err := t.ExecuteTemplate(&bodyBuffer, "html", data); err != nil {
			return fmt.Errorf("HTML template error (%s): %w", recipient.Email, err)
		}

		subject := fmt.Sprintf("Subject: %s\r\n", subjectBuffer.String())
		mime := "MIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
		message := []byte(subject + mime + bodyBuffer.String())

		to := []string{recipient.Email}

		// Retry loop
		delay := time.Duration(cfg.Retry.InitialDelay) * time.Second

		for attempt := 1; attempt <= cfg.Retry.MaxRetries; attempt++ {
			err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port), auth, cfg.SMTP.From, to, message)
			if err == nil {
				slog.Info("Email sent", "workerId", workerId, "email", recipient.Email, "attempt", attempt)
				break
			}

			// Last attempt → return error
			if attempt == cfg.Retry.MaxRetries {
				return fmt.Errorf("Failed after %d retries to send to %s: %w", attempt, recipient.Email, err)
			}

			slog.Warn("Email send failed, retrying...",
				"workerId", workerId,
				"email", recipient.Email,
				"retry", attempt,
				"wait_seconds", delay.Seconds(),
				"err", err.Error(),
			)

			time.Sleep(delay)

			delay = time.Duration(float64(delay) * cfg.Retry.BackoffMultiplier)

			maxDelay := time.Duration(cfg.Retry.MaxDelay) * time.Second
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	return nil
}
