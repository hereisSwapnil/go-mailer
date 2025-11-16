package producer

import (
	"encoding/csv"
	"log/slog"
	"os"
	"strings"

	"github.com/hereisSwapnil/go-mailer/internal/types"
)

func LoadRecipientsUsingCsv(filePath string, emailChannel chan types.Recipient) error {
	slog.Info("Loading recipients from CSV file", "filePath", filePath)
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		slog.Error("CSV file does not exist", "filePath", filePath)
		os.Exit(1)
	}
	if err != nil {
		slog.Error("Failed to access CSV file", "error", err, "filePath", filePath)
		os.Exit(1)
	}
	if fileInfo.Size() == 0 {
		slog.Error("CSV file is empty", "filePath", filePath)
		os.Exit(1)
	}

	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	// Get header row to map column names
	headers := records[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	// Process data rows
	for _, record := range records[1:] {
		recipient := types.Recipient{
			Extra: make(map[string]string),
		}

		// Extract Name and Email (required fields)
		if nameIdx, ok := headerMap["name"]; ok && nameIdx < len(record) {
			recipient.Name = strings.TrimSpace(record[nameIdx])
		}
		if emailIdx, ok := headerMap["email"]; ok && emailIdx < len(record) {
			recipient.Email = strings.TrimSpace(record[emailIdx])
		}

		// Store any other fields in Extra map
		for i, header := range headers {
			headerLower := strings.ToLower(strings.TrimSpace(header))
			if i < len(record) && headerLower != "name" && headerLower != "email" {
				recipient.Extra[headerLower] = strings.TrimSpace(record[i])
			}
		}

		emailChannel <- recipient
	}
	return nil
}