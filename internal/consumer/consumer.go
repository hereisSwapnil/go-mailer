package consumer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/hereisSwapnil/go-mailer/internal/types"
)

func EmailWorker(workerId int, emailChannel chan types.Recipient) error {

	auth := smtp.PlainAuth("", "8bp1alexgt1@gmail.com", "sewd omuo dyhc tfjh", "smtp.gmail.com")

	// Parse the email template only once
	t, err := template.ParseFiles("internal/templates/test_mail.tmpl")
	if err != nil {
		return err
	}

	for recipient := range emailChannel {
		fmt.Printf("Worker %d: Processing recipient %s\n", workerId, recipient.Name)

		// Prepare recipient specific data
		data := map[string]interface{}{
			"Name":    recipient.Name,
		}

		// Render subject template
		var subjectBuffer bytes.Buffer
		err := t.ExecuteTemplate(&subjectBuffer, "subject", data)
		if err != nil {
			return fmt.Errorf("subject template execution failed for recipient %s: %w", recipient.Name, err)
		}

		// Render HTML body template
		var bodyBuffer bytes.Buffer
		err = t.ExecuteTemplate(&bodyBuffer, "html", data)
		if err != nil {
			return fmt.Errorf("html template execution failed for recipient %s: %w", recipient.Name, err)
		}

		// Compose email headers and body
		subject := fmt.Sprintf("Subject: %s\r\n", subjectBuffer.String())
		mime := "MIME-Version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
		message := []byte(subject + mime + bodyBuffer.String())

		// Send to the actual recipient
		to := []string{"8bp1alexgt1@gmail.com"}
		err = smtp.SendMail("smtp.gmail.com:587", auth, "8bp1alexgt1@gmail.com", to, message)
		if err != nil {
			return fmt.Errorf("failed to send to %s: %w", recipient.Email, err)
		}
	}

	return nil
}